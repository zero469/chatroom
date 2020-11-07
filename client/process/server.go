package process

import (
	"fmt"
	"go_code/chapter18/project3/client/utils"
	"go_code/chapter18/project3/common/message"
	"net"
	"os"
)

var myID int

//ShowMenu 展示登录成功后的界面
func ShowMenu() {
	fmt.Println("----------------------恭喜登录成功----------------------")
	fmt.Println("----------------------1.显示在线用户列表-----------------")
	fmt.Println("----------------------2.发送信息------------------------")
	fmt.Println("----------------------3.信息列表------------------------")
	fmt.Println("----------------------4.退出系统------------------------")
	fmt.Println("请选择(1-4):")
	var key int
	var mesContent string
	fmt.Scanln(&key)
	switch key {
	case 1:
		fmt.Println("在线用户：")
		showOnlineUsers()
	case 2:
		fmt.Println("发送信息")
		fmt.Scanln(&mesContent)
		smsP := &SmsProcess{}
		smsP.SendGroupMes(mesContent)
	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Println("退出系统")
		os.Exit(0)
	default:
		fmt.Println("输入错误")
	}
}

func serverProcessMes(Conn net.Conn) {
	tf := &utils.Transfer{
		Conn: Conn,
	}
	for {
		fmt.Println("客户端等待服务器发送消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg err=", err)
			return
		}
		fmt.Println("mes = ", mes)
		switch mes.Type {
		case message.UpdataUserStateMesType:
			updateUserState(mes)
		default:
			fmt.Println("消息类型错误")
		}
	}
}
