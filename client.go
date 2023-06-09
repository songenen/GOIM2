package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ClientIP   string
	ClientPort int
	Name       string
	Conn       net.Conn
	flag       int
}

func NewClient(serverIP string, serverPort int) *Client {
	client := &Client{
		ClientIP:   serverIP,
		ClientPort: serverPort,
		flag:       77,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		fmt.Println("net Dial error:", err)
		return nil
	}
	client.Conn = conn

	return client
}

var serverIP string
var serverPort int

func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "设置链接ip地址")
	flag.IntVar(&serverPort, "port", 8090, "设置链接port")
}
func (c *Client) UpdateName() bool {
	fmt.Println("请输入用户名...")
	fmt.Scanln(&c.Name)

	sendMsg := "rename|" + c.Name + "\n"
	_, err := c.Conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write error:", err)
		return false
	}
	return true
}
func (c *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := c.Conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write error:", err)
		return
	}
}
func (c *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	c.SelectUsers()
	fmt.Println("请输入用户名>>>,exit退出")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println("请输入消息内容>>>,exit退出")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := c.Conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write error:", err)
					return
				}
			}

			chatMsg = ""
			fmt.Println("请输入消息内容>>>,exit退出")
			fmt.Scanln(&chatMsg)
		}
		c.SelectUsers()
		fmt.Println("请输入用户名>>>,exit退出")
		fmt.Scanln(&remoteName)
	}
}
func (c *Client) PublicChat() {
	var chatMsg string

	fmt.Println(">>>>请输入聊天内容，exit退出。<<<<")

	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) > 0 {
			sendMsg := chatMsg + "\n"
			_, err := c.Conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write error:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println(">>>>请输入聊天内容，exit退出。<<<<")
		fmt.Scanln(&chatMsg)

	}
}
func (c *Client) DealResponse() {
	io.Copy(os.Stdout, c.Conn)
}
func (c *Client) Run() {
	for c.flag != 0 {
		for c.Menu() != true {
		}
		switch c.flag {
		case 1:
			c.PublicChat()
			break

		case 2:
			c.PrivateChat()
			break

		case 3:
			c.UpdateName()
			break

		}
	}
	fmt.Println("关闭链接")
}
func (c *Client) Menu() bool {
	var flag int
	fmt.Println("1.公聊模式...")
	fmt.Println("2.私聊模式...")
	fmt.Println("3.更名模式...")
	fmt.Println("0.退出...")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		c.flag = flag
		return true
	} else {
		fmt.Println(">>>>>>请输入合法数字<<<<<<")
		return false
	}
}
func main() {
	flag.Parse()
	client := NewClient(serverIP, serverPort)
	if client == nil {
		fmt.Println("=======>链接失败<=======")
		return
	}
	go client.DealResponse()
	fmt.Println("=======>链接成功<=======")
	client.Run()
}
