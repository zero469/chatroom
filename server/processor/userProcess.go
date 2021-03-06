package processor

import (
	"chatroom/common/message"
	"chatroom/server/model"
	"chatroom/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

//UserProcess 负责和用户相关的操作
type UserProcess struct {
	Conn net.Conn
}

//Register 处理注册消息
func (ups *UserProcess) Register(mes *message.Message) (err error) {
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
		fmt.Printf("用户 %v 注册成功\n", user.UserName)
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

//Login 处理登录消息
func (ups *UserProcess) Login(mes *message.Message) (err error) {
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
		UserMgr.Add(loginMes.UserId, &model.ClientConn{
			Conn:     ups.Conn,
			UserName: user.UserName,
		})

		//广播在线用户列表
		//1. 给登录成功的用户发送在线用户列表
		loginResMes.Users = make([]message.User, 0)

		UserMgr.lock.RLock()
		for id, cc := range UserMgr.users {
			loginResMes.Users = append(loginResMes.Users,
				message.User{
					UserID:   id,
					UserName: cc.UserName,
				})
		}
		UserMgr.lock.RUnlock()

		//2. 更新其他在线用户的在线用户列表
		UpdateUserState(loginMes.UserId, message.UserOnlineState)

		fmt.Printf("用户 %v 登录成功\n", user.UserName)
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

func (ups *UserProcess) CheckPwd(mes *message.Message) (err error) {
	var checkPwdMes message.CheckOldPwdMes
	err = json.Unmarshal([]byte(mes.Data), &checkPwdMes)
	if err != nil {
		return fmt.Errorf("checkPwd failed : %v", err)
	}

	ok, err := model.MyUserDao.CheckPwd(checkPwdMes.ID, checkPwdMes.OldPwd)
	if err != nil {
		return err
	}

	var resMesData message.ChangePwdResMes
	if ok {
		resMesData.Code = message.CheckOldPwdSuccessCode
	} else {
		resMesData.Code = message.WrongPasswordCode
	}

	data, err := json.Marshal(resMesData)
	if err != nil {
		return
	}
	var resMes message.Message
	resMes.Type = message.ChangePwdResMesType
	resMes.Data = string(data)

	tfer := utils.Transfer{
		Conn: ups.Conn,
	}

	err = tfer.WriteMes(resMes)
	return err
}

func (ups *UserProcess) ChangePwd(mes *message.Message) (err error) {
	var changePwdMes message.ChangeNewPwdMes
	err = json.Unmarshal([]byte(mes.Data), &changePwdMes)
	if err != nil {
		return fmt.Errorf("checkPwd failed : %v", err)
	}

	err = model.MyUserDao.ChangePwd(changePwdMes.ID, changePwdMes.NewPwd)

	//服务器内部出错了是否还要给客户端发送response
	var resMesData message.ChangePwdResMes

	if err != nil {
		//服务器内部错误，将该信息返回给用户
		resMesData.Code = message.ServerErrorCode
	} else {
		resMesData.Code = message.ChangePwdSuccessCode
	}

	data, err := json.Marshal(resMesData)
	if err != nil {
		//序列化失败，消息无法发送，所以直接return
		return

	}
	var resMes message.Message
	resMes.Type = message.ChangePwdResMesType
	resMes.Data = string(data)

	tfer := utils.Transfer{
		Conn: ups.Conn,
	}

	err = tfer.WriteMes(resMes)
	return err
}
