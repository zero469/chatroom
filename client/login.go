package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go_code/chapter18/project3/common/message"
	"net"
)

func login(userID int, userPWD string) (err error) {
	//1. 连接服务器
	conn, err := net.Dial("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net.Dial() err=", err)
		return err
	}
	defer conn.Close()

	//2. 通过conn发送消息
	var mes message.Message
	mes.Type = message.LoginMesType

	//3. 构建loginMes结构体并序列化
	var loginMes message.LoginMes
	loginMes.UserId = userID
	loginMes.UserPwd = userPWD
	loginData, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return err
	}

	//4. 组合Message结构体并序列化
	mes.Data = string(loginData)
	mesData, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
	}

	//5. 先发送mesData的长度
	pkgByte := make([]byte, 4)
	binary.BigEndian.PutUint32(pkgByte[0:4], uint32(len(mesData)))
	n, err := conn.Write(pkgByte)
	if n != 4 || err != nil {
		fmt.Println("", err)
		return err
	}
	fmt.Println("客户端发送消息长度成功")

	//6. 发送mesData
	n, err = conn.Write(mesData)
	if err != nil {
		fmt.Println("conn.Write err=", err)
		return err
	}

	fmt.Printf("客户端成功发送%v个字节的数据%v\n", n, string(mesData))
	return nil
}
