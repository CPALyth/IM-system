package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
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

// 处理server回应的消息
func (this *Client) DealResponse() {
	// 读取接收缓存区的内容 到 标准输出中, 永久阻塞监听
	io.Copy(os.Stdout, this.conn)
}

func (this *Client) menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag) // 读取用户输入
	if flag >= 0 && flag <= 3 {
		this.flag = flag
		return true
	} else {
		fmt.Println("请输入0,1,2,3")
		return false
	}
}

func (this *Client) UpdateName() bool {
	fmt.Println(">>> 请输入用户名:")
	fmt.Scanln(&this.Name)
	msg := "rename|" + this.Name + "\n"
	_, err := this.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("conn write err:", err)
		return false
	}
	return true
}

func (this *Client) PublicChat() {
	// 提示用户输入消息
	var chatMsg string
	fmt.Println(">>> 请输入聊天内容, exit退出")
	fmt.Println(&chatMsg)
	for chatMsg != "exit" {
		// 消息不为空则发送给服务器
		if len(chatMsg) != 0 {
			msg := chatMsg + "\n"
			_, err := this.conn.Write([]byte(msg))
			if err != nil {
				fmt.Println("conn Write err:", err)
				break
			}
		}
		time.Sleep(time.Second)
		chatMsg = ""
		fmt.Println(">>> 请输入聊天内容, exit退出")
		fmt.Scanln(&chatMsg)
	}
}

func (this *Client) Run() {
	for this.flag != 0 {
		for !this.menu() {
		}

		switch this.flag {
		case 1:
			this.PublicChat()
			break
		case 2:
			fmt.Println("已选择私聊模式")
			break
		case 3:
			fmt.Println("更新用户名")
			this.UpdateName()
			break
		}
	}
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
	go client.DealResponse()
	fmt.Println("连接服务器成功")
	client.Run()
}
