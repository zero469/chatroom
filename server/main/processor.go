package main

import (
	"errors"
	"fmt"
	"go_code/chapter18/project3/common/message"
	"go_code/chapter18/project3/server/processor"
	"go_code/chapter18/project3/server/utils"
	"io"
	"net"
)

//Process 总控结构体
type Process struct {
	Conn net.Conn
}

//MainProcess 处理客户端连接函数
func (ps *Process) MainProcess() (err error) {
	for {

		//0. 创建transfer实例
		tfer := utils.Transfer{
			Conn: ps.Conn,
		}

		//1. 读取客户端发送的包并反序列化为结构体
		mes, err := tfer.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出")
				return err
			} else {
				fmt.Println("readPkg err=", err)
				return err
			}
		}

		fmt.Println(mes)

		//2. 调用统一处理接口处理不同的消息
		err = ps.ServerProcess(&mes)
		if err != nil {
			fmt.Println("ServerProcess err=", err)
			return err
		}
	}
}

//ServerProcess 根据客户端发送的消息种类调用不同的处理函数
func (ps *Process) ServerProcess(mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:
		//处理登陆逻辑
		ups := &processor.UserProcess{
			Conn: ps.Conn,
		}
		return ups.ServerProcessLogin(mes)

	case message.RegisterMesType:
		ups := &processor.UserProcess{
			Conn: ps.Conn,
		}
		return ups.ServerProcessRegister(mes)
	default:
		err = errors.New("消息类型不存在")
		return
	}
}
