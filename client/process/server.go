package process

import (
	"fmt"
	"os"
)

var myId int

//Show 展示登录成功后的界面
func ShowMenu() {
	fmt.Println("----------------------恭喜登录成功----------------------")
	fmt.Println("----------------------1.显示在线用户列表-----------------")
	fmt.Println("----------------------2.发送信息------------------------")
	fmt.Println("----------------------3.信息列表------------------------")
	fmt.Println("----------------------4.退出系统------------------------")
	fmt.Println("请选择(1-4):")
	var key int
	fmt.Scanln(&key)
	switch key {
	case 1:
		fmt.Println("显示在线用户列表")
		showOnlineUsers()
	case 2:
		fmt.Println("发送信息")
	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Println("退出系统")
		os.Exit(0)
	default:
		fmt.Println("输入错误")
	}
}
