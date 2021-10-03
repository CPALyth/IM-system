package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string // 当前客户端的ip地址
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	// 启动go程
	go user.ListenMessage()
	return user
}

// 用户的上线业务
func (this *User) Online() {
	// 用户上线, 将用户加入到OnlineMap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// 广播当前用户的上线消息
	this.server.BroadCast(this, "已上线")
}

// 用户的下线业务
func (this *User) Offline() {
	// 用户下线, 将用户从OnlineMap中删除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// 广播当前用户的下线消息
	this.server.BroadCast(this, "已下线")
}

func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// 用户处理消息的业务
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		// 查询当前在线用户
		this.server.mapLock.Lock()
		for userName, _ := range this.server.OnlineMap {
			this.SendMsg(userName + "在线\n")
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 更改自己的用户名, 消息格式: rename|张三
		newName := strings.Split(msg, "|")[1]
		// 判断name是否存在
		_, err := this.server.OnlineMap[newName]
		if err {
			this.SendMsg("失败, 该用户名已被其它人使用\n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
			this.SendMsg("您已成功更新为用户名:" + this.Name + "\n")
		}
	} else {
		this.server.BroadCast(this, msg)
	}
}

// 监听当前User的channel,一旦有消息就直接发给客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
