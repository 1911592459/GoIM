package service

import (
	"GoIM/models"
	"GoIM/utils"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetIndex
// @Summary 所有用户
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	userList := models.GetUserList()
	c.JSON(http.StatusOK, gin.H{
		//0 成功  -1失败
		"code":    0,
		"Message": userList,
	})
}

// FindUserByNameAndPwd
// @Summary 查询指定用户
// @Tags 用户模块
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	name := c.PostForm("name")
	password1 := c.PostForm("password")
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(200, gin.H{
			//0 成功  -1失败
			"code":    -1,
			"message": "该用户不存在！"})
		return
	}
	flag := utils.ValidPassword(password1, user.Salt, user.PassWord)
	if !flag {
		c.JSON(200, gin.H{
			//0 成功  -1失败
			"code":    -1,
			"message": "密码不正确！"})
		return
	}
	password := utils.MakePassword(password1, user.Salt)
	data := models.FindUserByNameAndPwd(name, password)
	c.JSON(http.StatusOK, gin.H{
		//0 成功  -1失败
		"code":    0,
		"message": "消息发送成功",
		"data":    data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @param repassword formData string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [post]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.PostForm("name")
	password1 := c.PostForm("password")
	password2 := c.PostForm("repassword")
	if user.Name == "" || password1 == "" {
		c.JSON(200, gin.H{
			//0 成功  -1失败
			"code":    -1,
			"message": "用户名和密码不能为空！"})
		return
	}
	data := models.FindUserByName(user.Name)

	if data.Name != "" {
		c.JSON(200, gin.H{
			//0 成功  -1失败
			"code":    -1,
			"message": "用户名已注册！"})
		return
	}
	if password1 != password2 {
		c.JSON(200, gin.H{
			//0 成功  -1失败
			"code":    -1,
			"message": "俩次密码不一致！"})
		return
	}
	salt := fmt.Sprintf("%06d", rand.Int31())
	//加密
	user.PassWord = utils.MakePassword(password1, salt)
	user.Salt = salt
	//user.PassWord = password1
	models.CreateUser(&user)
	c.JSON(http.StatusOK, gin.H{
		//0 成功  -1失败
		"code":    0,
		"message": "新增用户成功",
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id formData string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [post]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	models.DeleteUser(&user)
	c.JSON(http.StatusOK, gin.H{
		//0 成功  -1失败
		"code":    0,
		"message": "删除用户成功",
	})
}

// UpdateUser
// @Summary 更新用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	user.ID = uint(id)
	//使用第三方工具类govalidator去校验邮箱格式
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			//0 成功  -1失败
			"code":    -1,
			"message": "修改参数格式不正确",
		})
		return
	}
	models.UpdateUser(&user)
	c.JSON(http.StatusOK, gin.H{
		//0 成功  -1失败
		"code":    0,
		"message": "更新用户成功",
	})
}

//防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	conn, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err2 := ws.Close()
		if err2 != nil {
			fmt.Println(err2)
		}
	}(conn)
	MsgHandler(conn, c)
}

func MsgHandler(conn *websocket.Conn, c *gin.Context) {
	msg, err := utils.Subscribe(c, utils.PublishKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	format := time.Now().Format("2006-01-02 15:04:05")
	m := fmt.Sprintf("[ws][%s]:%s", format, msg)
	err = conn.WriteMessage(1, []byte(m))
	if err != nil {
		fmt.Println(err)
	}
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

func SearchFriends(c *gin.Context) {
	id, _ := strconv.Atoi(c.Request.FormValue("userId"))
	users := models.SearchFriend(uint(id))
	// c.JSON(200, gin.H{
	// 	"code":    0, //  0成功   -1失败
	// 	"message": "查询好友列表成功！",
	// 	"data":    users,
	// })
	utils.RespOKList(c.Writer, users, len(users))
}

func AddFriend(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	targetName := c.Request.FormValue("targetName")
	//targetId, _ := strconv.Atoi(c.Request.FormValue("targetId"))
	code, msg := models.AddFriend(uint(userId), targetName)
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

//新建群
func CreateCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	name := c.Request.FormValue("name")
	icon := c.Request.FormValue("icon")
	desc := c.Request.FormValue("desc")
	community := models.Community{}
	community.OwnerId = uint(ownerId)
	community.Name = name
	community.Img = icon
	community.Desc = desc
	code, msg := models.CreateCommunity(community)
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

//加载群列表
func LoadCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	//	name := c.Request.FormValue("name")
	data, msg := models.LoadCommunity(uint(ownerId))
	if len(data) != 0 {
		utils.RespList(c.Writer, 0, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

//加入群 userId uint, comId uint
func JoinGroups(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	comId := c.Request.FormValue("comId")

	//	name := c.Request.FormValue("name")
	data, msg := models.JoinGroup(uint(userId), comId)
	if data == 0 {
		utils.RespOK(c.Writer, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

func FindByID(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))

	//	name := c.Request.FormValue("name")
	data := models.FindByID(uint(userId))
	utils.RespOK(c.Writer, data, "ok")
}

func RedisMsg(c *gin.Context) {
	userIdA, _ := strconv.Atoi(c.PostForm("userIdA"))
	userIdB, _ := strconv.Atoi(c.PostForm("userIdB"))
	start, _ := strconv.Atoi(c.PostForm("start"))
	end, _ := strconv.Atoi(c.PostForm("end"))
	isRev, _ := strconv.ParseBool(c.PostForm("isRev"))
	res := models.RedisMsg(int64(userIdA), int64(userIdB), int64(start), int64(end), isRev)
	utils.RespOKList(c.Writer, "ok", res)
}
