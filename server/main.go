package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go_code/chapter18/project3/common/message"
	"net"
)

//读取客户端发送的信息
func process(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 8096)

	for {
		//0.1 读取mesData的长度
		_, err := conn.Read(buf[:4])
		if err != nil {
			fmt.Println("conn.Read err = ", err)
			return
		}
		pkgLen := int(binary.BigEndian.Uint32(buf[:4]))
		fmt.Println("服务器读取到mesData的长度 : ", pkgLen)

		//0.2 读取mesData
		n, err := conn.Read(buf)
		if n != pkgLen || err != nil {
			fmt.Println("conn.Read err=", err)
			return
		}

		//1. 反序列化message结构体
		var mes message.Message
		err = json.Unmarshal(buf[:n], &mes)
		if err != nil {
			fmt.Println("json.Unmarshal err=", err)
			return
		}

		//2. 根据Type反序列化data 结构体
		switch mes.Type {
		case message.LoginMesType:
			var loginMes message.LoginMes
			err = json.Unmarshal([]byte(mes.Data), &loginMes)
			fmt.Printf("用户id : %v, 用户密码 : %v\n", loginMes.UserId, loginMes.UserPwd)

		default:
			fmt.Println("未定义Type = ", mes.Type)
		}
	}
}

func main() {
	fmt.Println("开始在8889端口监听")
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
