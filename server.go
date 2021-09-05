package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip string
	Port int
}

// 创建
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip: ip,
		Port: port,
	}
	return server
}

// 启动服务器的接口函数
func (this *Server) Start() {
	// 创建一个监听套接字
	net.Listen()
}