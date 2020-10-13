package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go_code/chapter18/project3/common/message"
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
	//TODO: conn.Read()能保证读到这么多的消息吗？如果不指定要读的长度会发生什么？ 这里表示期望读到pkgLen这么长的数据，但是实际可能读不到这么多（丢包？？）
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
