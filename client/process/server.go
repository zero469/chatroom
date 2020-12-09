package process

import (
	"fmt"
	"net"
	"os"

	"chatroom/client/utils"
	"chatroom/common/message"
)

//ShowMenu 展示登录成功后的界面
func ShowMenu() {
	fmt.Println("----------------------主菜单----------------------")
	fmt.Println("                      1.显示在线用户列表")
	fmt.Println("                      2.发送信息")
	fmt.Println("                      3.信息列表")
	fmt.Println("                      4.退出系统")
	fmt.Println("                      5.修改密码")
	fmt.Println("请选择(1-4):")
	var key int
	var mesContent string
	fmt.Scanln(&key)
	switch key {
	case 1:
		onlineUsers.showOnlineUsers()
	case 2:
		fmt.Println("请输入信息：")
		_, mesContent = utils.ReadLine()
		smsP := &SmsProcess{}
		smsP.SendGroupMes(mesContent)
	case 3:
		fmt.Println("信息列表: ")
		HistoryMes.ShowMesList()
	case 4:
		fmt.Println("退出系统")
		os.Exit(0)
	case 5:
		fmt.Println("修改密码") //分两步 1.验证原密码 2.输入新密码
		up := &UserProcess{}
		up.ChangePwd()
	default:
		fmt.Println("输入错误")
	}
}

func initAll() {
	initMesMgr()
}

func serverProcessMes(Conn net.Conn) {

	initAll()
	tf := &utils.Transfer{
		Conn: Conn,
	}
	for {
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg err=", err)
			return
		}
		switch mes.Type {
		//更新在线用户状态
		case message.UpdataUserStateMesType:
			onlineUsers.updateUserState(mes)
		//接受其他用户发送的消息
		case message.SmsResMesType:
			rcvSmsMes(mes)
		default:
			fmt.Println("消息类型错误")
		}
	}
}
