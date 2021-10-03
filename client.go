package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}
	// 连接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}
	client.conn = conn
	return client
}

var serverIp string
var serverPort int

// 此函数会在main函数之前被自动运行, client -i 127.0.0.1 -p 8888
func init() {
	flag.StringVar(&serverIp, "i", "127.0.0.1", "设置服务器的IP地址(默认是127.0.0.1)")
	flag.IntVar(&serverPort, "p", 8888, "设置服务器的端口号(默认是8888)")
	flag.Parse() // 命令行解析
}

func main() {
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("连接服务器失败")
		return
	}
	fmt.Println("连接服务器成功")
	// 未完待续
	select {}
}
