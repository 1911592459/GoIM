package service

import (
	"GoIM/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"strconv"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	/*	files, err := template.ParseFiles("index.html")
		if err != nil {
			panic(err)
		}
		files.Execute(c.Writer, "index")
	*/

	ind, err := template.ParseFiles("index.html", "views/chat/head.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "index")

}
func ToRegister(c *gin.Context) {
	ind, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "register")
}

func ToChat(c *gin.Context) {
	ind, err := template.ParseFiles("views/chat/index.html",
		"views/chat/head.html",
		"views/chat/foot.html",
		"views/chat/tabmenu.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/createcom.html",
		"views/chat/userinfo.html",
		"views/chat/main.html")
	if err != nil {
		panic(err)
	}
	userId, err := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")
	if err != nil {
		fmt.Println(err)
		return
	}
	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token

	ind.Execute(c.Writer, user)
}
func Chat(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
