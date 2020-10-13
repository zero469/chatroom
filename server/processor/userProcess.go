package processor


import (
	"go_code/chapter18/project3/server/utils"
	"encoding/json"
	"fmt"
	"go_code/chapter18/project3/common/message"
	"net"
)

//UserProcess 负责和用户相关的操作
type UserProcess struct{
	Conn net.Conn
}

func (userProcess *UserProcess)ServerProcessLogin(mes *message.Message) (err error) {
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
	if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
		//合法
		loginResMes.Code = message.LoginSuccessCode
	} else {
		//不合法
		loginResMes.Code = message.UnRegisterCode
		loginResMes.Error = "该用户不存在，请注册再使用"
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
		Conn : userProcess.Conn,
	}
	err = tfer.WritePkg(data)
	return
}