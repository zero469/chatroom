package main

import (
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

	//5. 发送mes包
	err = writePkg(conn, mesData)
	if err != nil {
		fmt.Println("writePkg fail", err)
		return
	}

	//6. 读服务器的response包
	mes, err = readPkg(conn)
	if err != nil {
		fmt.Println("readPkg fail :", err)
		return
	}

	//7. 解析response包
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	if loginResMes.Code == message.LoginSuccessCode {
		fmt.Println("登录成功")
	} else if loginResMes.Code == message.UnRegisterCode {
		fmt.Println(loginResMes.Error)
	}

	return nil
}
