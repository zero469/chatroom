package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"go_code/chapter18/project3/common/message"
	"io"
	"net"
)

/*
	读包，并将包反序列化成结构体
*/
func readPkg(conn net.Conn) (mes message.Message, err error) {
	buf := make([]byte, 8096)

	//1 读取mesData的长度
	_, err = conn.Read(buf[:4])
	if err != nil {
		//err = errors.New("read pkg header error")
		return
	}
	pkgLen := int(binary.BigEndian.Uint32(buf[:4]))
	fmt.Println("服务器读取到mesData的长度 : ", pkgLen)

	//2 读取mesData
	//TODO: conn.Read()能保证读到这么多的消息吗？如果不指定要读的长度会发生什么？
	//这里表示期望读到pkgLen这么长的数据，但是实际可能读不到这么多（丢包？？）
	n, err := conn.Read(buf[:pkgLen])
	if n != pkgLen || err != nil {
		//err = errors.New("read pkg body error")
		return
	}

	//3. 反序列化message结构体
	err = json.Unmarshal(buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return
}

/*
	发包函数
*/
func writePkg(conn net.Conn, data []byte) (err error) {
	//1. 发送包长度
	pkgByte := make([]byte, 4)
	binary.BigEndian.PutUint32(pkgByte[0:4], uint32(len(data)))
	n, err := conn.Write(pkgByte)
	if n != 4 || err != nil {
		fmt.Println("conn.Write(pkgByte) fail ", err)
		return err
	}

	//2. 发送包本身
	n, err = conn.Write(data)
	if n != len(data) || err != nil {
		fmt.Println("conn.Write(data) fail ", err)
		return
	}
	return
}

func serverProcessLogin(conn net.Conn, mes *message.Message) (err error) {
	var loginMes message.LoginMes

	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err)
		return
	}

	//构造response消息
	var resMes message.Message
	resMes.Type = message.LoginMesResType

	var loginResMes message.LoginResMes
	if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
		//合法
		loginResMes.Code = message.LoginSuccessCode
	} else {
		//不合法
		loginResMes.Code = message.UnRegisterCode
		loginResMes.Error = "该用户不存在，请注册再使用"
	}

	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal(loginMes) err=", err)
		return
	}
	resMes.Data = string(data)

	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal(resMes) err=", err)
		return
	}
	err = writePkg(conn, data)
	return
}

//ServerProcess 根据客户端发送的消息种类调用不同的处理函数
func ServerProcess(conn net.Conn, mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:
		//处理登陆逻辑
		return serverProcessLogin(conn, mes)
	default:
		err = errors.New("消息类型不存在")
		return
	}
}

/*
	处理客户端连接函数
*/
func process(conn net.Conn) {
	defer conn.Close()

	for {
		//1. 读取客户端发送的包并反序列化为结构体
		mes, err := readPkg(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出")
				return
			} else {
				fmt.Println("readPkg err=", err)
				return
			}
		}

		fmt.Println(mes)

		//2. 调用统一处理接口处理不同的消息
		err = ServerProcess(conn, &mes)
		if err != nil {
			fmt.Println("ServerProcess err=", err)
			return
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
