package utils

import (
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

//Transfer 处理收发包
type Transfer struct {
	Conn net.Conn
	Buf  [8096]byte //传输缓冲
}

//读包，并将包反序列化成结构体
func (transfer *Transfer) ReadPkg() (mes message.Message, err error) {

	//1 读取mesData的长度
	_, err = transfer.Conn.Read(transfer.Buf[:4])
	if err != nil {
		return
	}
	pkgLen := int(binary.BigEndian.Uint32(transfer.Buf[:4]))

	//2 读取mesData
	n, err := transfer.Conn.Read(transfer.Buf[:pkgLen])
	if n != pkgLen || err != nil {
		return
	}

	//3. 反序列化message结构体
	err = json.Unmarshal(transfer.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return
}

// 发包函数
func (transfer *Transfer) WritePkg(data []byte) (err error) {
	//1. 发送包长度
	pkgByte := make([]byte, 4)
	binary.BigEndian.PutUint32(pkgByte[0:4], uint32(len(data)))
	n, err := transfer.Conn.Write(pkgByte)
	if n != 4 || err != nil {
		fmt.Println("conn.Write(pkgByte) fail ", err)
		return err
	}

	//2. 发送包本身
	n, err = transfer.Conn.Write(data)
	if n != len(data) || err != nil {
		fmt.Println("conn.Write(data) fail ", err)
		return
	}
	return
}

func (transfer *Transfer) WriteMes(mes message.Message) (err error) {
	data, err := json.Marshal(mes)
	if err != nil {
		return
	}
	err = transfer.WritePkg(data)
	return
}
