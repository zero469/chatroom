package processor

import (
	"encoding/json"
	"fmt"
	"go_code/chapter18/project3/common/message"
	"go_code/chapter18/project3/server/model"
	"go_code/chapter18/project3/server/utils"
	"net"
)

//UserProcess 负责和用户相关的操作
type UserProcess struct {
	Conn net.Conn
}

//ServerProcessRegister 处理注册消息
func (ups *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	var regiMes message.RegisterMes

	err = json.Unmarshal([]byte(mes.Data), &regiMes)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err)
		return
	}

	user := &regiMes.User
	err = model.MyUserDao.Register(user)

	var regiResMes message.RegisterResMes
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			regiResMes.Code = message.UserIdBeenUsedCode
			regiResMes.Error = err.Error()
		} else {
			regiResMes.Code = message.ServerErrorCode
			regiResMes.Error = "服务器内部错误"
		}
	} else {
		regiResMes.Code = message.RegisterSuccessCode
		fmt.Printf("用户 %v 注册成功", user.UserName)
	}

	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	data, err := json.Marshal(regiResMes)
	if err != nil {
		fmt.Println("json.Marshal(regiResMes) err=", err)
		return
	}
	resMes.Data = string(data)

	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal(resMes) err=", err)
		return
	}

	//创建Transfer实例

	tfer := &utils.Transfer{
		Conn: ups.Conn,
	}
	err = tfer.WritePkg(data)
	return
}

//ServerProcessLogin 处理登录消息
func (ups *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
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

	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = message.UnRegisterCode
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = message.WrongPasswordCode
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = message.ServerErrorCode
			loginResMes.Error = "服务器内部错误"
		}
	} else {
		loginResMes.Code = message.LoginSuccessCode
		//登录成功后将该用户加入到UserMgr中
		UserMgr.Add(loginMes.UserId, ups)

		//广播在线用户列表
		//TODO:重构上线广播代码，使用独立的消息类型，后续可扩展下线广播
		//1. 给登录成功的用户发送在线用户列表
		loginResMes.Users = make([]int, 0)
		for id := range UserMgr.users {
			fmt.Println("userid : ", id, "up : ", UserMgr.users[id])
			loginResMes.Users = append(loginResMes.Users, id)
		}
		//2. 更新其他在线用户的在线用户列表
		UpdateUserState(loginMes.UserId, message.UserOnlineState)

		fmt.Printf("用户 %v 登录成功", user.UserName)
	}

	data, err := json.Marshal(loginResMes)
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

	//创建Transfer实例

	tfer := &utils.Transfer{
		Conn: ups.Conn,
	}
	err = tfer.WritePkg(data)
	return
}

func (ups *UserProcess) updateUserState(mes *message.UpdataUserStateMes) (err error) {
	//构造消息结构体
	var updateMes message.Message
	updateMes.Type = message.UpdataUserStateMesType

	data, err := json.Marshal(*mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return err
	}

	updateMes.Data = string(data)

	data, err = json.Marshal(updateMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return err
	}

	tf := &utils.Transfer{
		Conn: ups.Conn,
	}
	err = tf.WritePkg(data)
	return err
}
