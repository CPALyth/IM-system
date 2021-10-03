package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播的channel
	Message chan string
}

// 创建
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 广播消息的方法
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Name + "] " + msg
	this.Message <- sendMsg
}

// 监听message消息的GO程, 一旦有消息就发给全部在线的User
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		// 将msg发送给全部在线的User
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 处理业务
func (this *Server) Handler(conn net.Conn) {
	user := NewUser(conn, this)
	user.Online()

	// 监听用户是否活跃的channel
	isLive := make(chan bool)

	// 接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			// 提取用户的消息
			msg := string(buf[:n-1])
			// 用户针对消息进行处理
			user.DoMessage(msg)
			// 把用户设为活跃
			isLive <- true
		}
	}()

	// 当前handler阻塞
	for {
		select {
		case <-isLive:
			// 不做任何事情, 为了激活select, 更新下面的定时器
		case <-time.After(time.Second * 300):
			// 已经超时, 将当前用户踢出群聊
			user.SendMsg("你被踢了\n")
			// 销毁管道
			close(user.C)
			// 关闭连接
			conn.Close()
			// 退出当前的handler
			return // runtime.Goexit()
		}
	}
}

// 启动服务器的接口函数
func (this *Server) Start() {
	// 创建一个监听套接字
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	defer listener.Close()

	// 启动监听消息的go程
	go this.ListenMessage()

	// 创建新套接字, 处理客户端的请求
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		go this.Handler(conn)
	}
}
