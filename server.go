package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// 创建
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

// 处理业务
func (this *Server) Handler(conn net.Conn) {
	fmt.Println("成功建立连接")
}

// 启动服务器的接口函数
func (this *Server) Start() {
	// 创建一个监听套接字
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {

	}
	defer listener.Close()
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
