package main

import "net"

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

// 用户处理消息的业务
func (this *User) DoMessage(msg string) {
	// 将消息进行广播
	this.server.BroadCast(this, msg)
}

// 监听当前User的channel,一旦有消息就直接发给客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
