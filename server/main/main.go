package main

import (
	"fmt"
	"net"
)



//处理客户端连接函数
func process(conn net.Conn) {
	defer conn.Close()

	psor := &Process{
		Conn : conn,
	}
	err := psor.MainProcess()
	if err != nil{
		fmt.Println("客户端和服务器通讯协程错误 err=", err)
		return
	}
}

func main() {
	fmt.Println("新结构 开始在8889端口监听")
	listen, err := net.Listen("tcp", "0.0.0.0:8889") //监听端口8889
	defer listen.Close()

	if err != nil {
		fmt.Println("监听失败...")
		return
	}

	for {
		fmt.Println("等待连接")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept() err=", err)
			continue
		}
		//连接成功，起一个协程保持通讯
		go process(conn)
	}
}
