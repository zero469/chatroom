package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go_code/chapter18/project3/common/message"
	"io"
	"net"
)

//读客户端的包，并将包反序列化成结构体
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

//处理客户端连接
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
		// //2. 根据Type反序列化data 结构体
		// switch mes.Type {
		// case message.LoginMesType:
		// 	var loginMes message.LoginMes
		// 	err = json.Unmarshal([]byte(mes.Data), &loginMes)
		// 	fmt.Printf("用户id : %v, 用户密码 : %v\n", loginMes.UserId, loginMes.UserPwd)

		// default:
		// 	fmt.Println("未定义Type = ", mes.Type)
		// }
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
