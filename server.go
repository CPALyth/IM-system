package main

import (
	"fmt"
	"io"
	"net"
	"sync"
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
	sendMsg := "[" + user.Addr + "] " + msg
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
	user := NewUser(conn)

	// 用户上线, 将用户加入到OnlineMap中
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	// 广播当前用户的上线消息
	this.BroadCast(user, "已上线")

	// 接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				this.BroadCast(user, "已下线")
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			// 提取用户的消息
			msg := string(buf[:n-1])
			// 将消息进行广播
			this.BroadCast(user, msg)
		}
	}()

	// 当前handler阻塞
	select {}
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
