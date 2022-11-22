package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"sync"
)

type Message struct {
	gorm.Model
	//发送者id，即消息是从哪个人发过来的
	FormId int64
	//接收者id，消息是发给谁的
	TargetId int64
	//发送类型，2群聊，1私聊，3广播等等
	Type int
	//消息类型，1文字，4音频，2表情包，3图片等
	Media int
	//消息内容
	Content string
	//图片相关
	Pic string
	//url相关
	Url  string
	Desc string
	//其他数字统计
	Amount int
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

//映射关系
var clientMap = make(map[int64]*Node, 0)

//读写锁
var rwLocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {
	//1,获取参数并校验token等合法性

	//
	query := request.URL.Query()
	id := query.Get("userId")
	userId, _ := strconv.ParseInt(id, 10, 64)
	/*if err2 != nil {
		fmt.Println(err2)
		return
	}*/
	/*msgType := query.Get("type")
	//token := query.Get("token")
	targetId := query.Get("targetId")
	context := query.Get("context")*/
	isValida := true //checkToken()
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isValida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//2,获取连接conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	//3,用户关系

	//4,userId node 绑定并加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	//5,完成发送逻辑
	go sendProc(node)
	//6,完成接收逻辑
	go recvProc(node)
	sendMsg(userId, []byte("欢迎进入聊天室"))
}
func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()

		if err != nil {
			fmt.Println(err)
			return
		}
		broadMsg(data)
		fmt.Println("[ws]>>>>>>>>>>>>>>>", data)
	}
}

var udpSendChan = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpSendChan <- data
}

func init() {
	go udpSendProc()
	go udpRecvProc()
}

//完成udp数据发送协程
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 255),
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	for {
		select {
		case data := <-udpSendChan:
			_, err := con.Write(data)
			fmt.Println("udp客户端写数据发送给服务端")
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

//完成udp数据接收协程
func udpRecvProc() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	for {
		var buf [512]byte
		read, err := conn.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("udpRecvProc  data :", string(buf[0:read]))
		dispatch(buf[0:read])
	}
}

//后端逻辑调度处理
func dispatch(data []byte) {
	msg := &Message{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1:
		//私聊
		sendMsg(msg.TargetId, data)
		/*	case 2:
				//群聊
				sendGroupMsg()
			case 3:
				//广播
				sendAllMsg()
			case 4:
				//默认
		*/
	}
}

func sendMsg(userId int64, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}

}
