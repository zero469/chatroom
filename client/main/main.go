package main

import (
	"fmt"
	"go_code/chapter18/project3/client/process"
)

var userID int
var userPWD string
var userName string

func main() {
	var key int
	var loop = true
	for loop {
		fmt.Println("-----------------欢迎登陆多人聊天系统-----------------")
		fmt.Println("\t\t\t1 登陆聊天室")
		fmt.Println("\t\t\t2 注册用户")
		fmt.Println("\t\t\t3 退出系统")
		fmt.Printf("\t\t\t请选择(1-3):")
		fmt.Scanln(&key)
		switch key {
		case 1:
			fmt.Println("登陆聊天室")
			fmt.Printf("请输入用户ID :")
			fmt.Scanln(&userID)

			fmt.Printf("请输入用户密码 :")
			fmt.Scanln(&userPWD)

			up := &process.UserProcess{}
			_ = up.Login(userID, userPWD)

		case 2:
			fmt.Println("注册用户")
			fmt.Printf("请输入用户ID :")
			fmt.Scanln(&userID)

			fmt.Printf("请输入用户密码 :")
			fmt.Scanln(&userPWD)

			fmt.Printf("请输入用户名 :")
			fmt.Scanln(&userName)
			
			up := &process.UserProcess{}
			_ = up.Register(userID, userPWD, userPWD)
		case 3:
			fmt.Println("退出系统")
			loop = false
		default:
			fmt.Println("输入有误，请重新输入")
		}
	}

}
