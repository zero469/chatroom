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

func (userProcess *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
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
		fmt.Printf("%v 登陆成功", user.UserName)
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
		Conn: userProcess.Conn,
	}
	err = tfer.WritePkg(data)
	return
}
