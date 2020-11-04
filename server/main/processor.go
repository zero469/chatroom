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
			//此处表示用户退出客户端，即以及下线，需要广播该用户的下线状态，从UserMgr中拿到userID，然后将该用户从UserMgr中删除并广播该用户已下线
			userID, err := processor.UserMgr.GetIDbyConn(ps.Conn)
			if err == nil {
				processor.UserMgr.Del(userID)
				processor.UpdateUserState(userID, message.UserOfflineState)
			}
			if err == io.EOF {
				fmt.Println("客户端退出")
				return err
			}
			return fmt.Errorf("MainProcess tfer.ReadPkg failed : %v", err)

		}

		fmt.Println(mes)

		//2. 调用统一处理接口处理不同的消息
		err = ps.ServerProcess(&mes)
		if err != nil {
			fmt.Println("ServerProcess failed :", err)
			return err
		}
	}
}

//ServerProcess 根据客户端发送的消息种类调用不同的处理函数
func (ps *Process) ServerProcess(mes *message.Message) (err error) {
	fmt.Println(mes)
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
		err = errors.New("ServerProcess failed : 消息类型不存在")
		return
	}
}
