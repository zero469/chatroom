package process

import (
	"chatroom/client/model"
	"chatroom/client/utils"
	"chatroom/common/message"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

type UserProcess struct {
}

func (up *UserProcess) Register(userID int, userPWD string, userName string) (err error) {
	conn, err := net.Dial("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net.Dial() err=", err)
		return err
	}
	defer conn.Close()

	var mes message.Message
	mes.Type = message.RegisterMesType

	var RegiMes message.RegisterMes
	RegiMes.User.UserID = userID
	RegiMes.User.UserPwd = userPWD
	RegiMes.User.UserName = userName
	RegiData, err := json.Marshal(RegiMes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return err
	}

	mes.Data = string(RegiData)
	mesData, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
	}

	tf := &utils.Transfer{
		Conn: conn,
	}

	err = tf.WritePkg(mesData)
	if err != nil {
		fmt.Println("writePkg fail", err)
		return
	}

	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg fail :", err)
		return
	}

	var RegiResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &RegiResMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	if RegiResMes.Code == message.RegisterSuccessCode {
		fmt.Println("注册成功，请重新登陆")
	} else {
		fmt.Println(RegiResMes.Error)
	}
	return nil
}

func (up *UserProcess) Login(userID int, userPWD string) (err error) {
	//1. 连接服务器
	conn, err := net.Dial("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net.Dial() err=", err)
		return err
	}
	defer conn.Close()

	//2. 通过conn发送消息
	var mes message.Message
	mes.Type = message.LoginMesType

	//3. 构建loginMes结构体并序列化
	var loginMes message.LoginMes
	loginMes.UserId = userID
	loginMes.UserPwd = userPWD
	loginData, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return err
	}

	//4. 组合Message结构体并序列化
	mes.Data = string(loginData)
	mesData, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
	}

	tf := &utils.Transfer{
		Conn: conn,
	}
	//5. 发送mes包
	err = tf.WritePkg(mesData)
	if err != nil {
		fmt.Println("writePkg fail", err)
		return
	}

	//6. 读服务器的response包
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg fail :", err)
		return
	}

	//7. 解析response包
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	if loginResMes.Code == message.LoginSuccessCode {
		fmt.Println("登录成功")

		model.InitCurUser(userID, conn)

		//1. 初始化UserMgr
		initUserMgr()
		for _, user := range loginResMes.Users {
			onlineUsers.userList[user.UserID] = &message.User{
				UserID:    user.UserID,
				UserState: message.UserOnlineState,
				UserName:  user.UserName,
			}
		}
		//2. 打印在线用户
		onlineUsers.showOnlineUsers()

		go serverProcessMes(conn)
		//3. 显示登陆成功的界面
		for {
			ShowMenu()
		}
	} else {
		fmt.Println(loginResMes.Error)
	}

	return nil
}

func (up *UserProcess) checkOldPwd(oldPwd string) (ok bool, err error) {
	ok = true
	err = nil

	var mes message.Message
	mes.Type = message.CheckOldPwdMesType

	var dataMes message.CheckOldPwdMes
	dataMes.OldPwd = oldPwd
	dataMes.ID = model.CurUser.UserID

	data, err := json.Marshal(dataMes)
	if err != nil {
		return false, fmt.Errorf("checkOldPwd failed : %v", err)
	}
	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		return false, fmt.Errorf("checkOldPwd failed : %v", err)
	}

	tfer := utils.Transfer{
		Conn: model.CurUser.Conn,
	}
	err = tfer.WritePkg(data)
	if err != nil {
		return false, fmt.Errorf("checkOldPwd failed : %v", err)
	}

	res, err := tfer.ReadPkg()
	if err != nil {
		return false, fmt.Errorf("checkOldPwd failed : %v", err)
	}

	var resMes message.ChangePwdResMes
	err = json.Unmarshal([]byte(res.Data), &resMes)
	if err != nil {
		return false, fmt.Errorf("json.Unmarshal err= %v", err)
	}

	switch resMes.Code {
	case message.CheckOldPwdSuccessCode:
		return true, nil
	case message.WrongPasswordCode:
		return false, errors.New("Wrong password")
	default:
		return false, errors.New("Server internal error")
	}
}

func (up *UserProcess) changeNewPwd(newPwd string) (err error) {
	var mes message.Message
	mes.Type = message.ChangeNewPwdMesType

	var dataMes message.ChangeNewPwdMes
	dataMes.NewPwd = newPwd

	data, err := json.Marshal(dataMes)
	if err != nil {
		return fmt.Errorf("changeNewPwd failed : %v", err)
	}
	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		return fmt.Errorf("changeNewPwd failed : %v", err)
	}

	tfer := utils.Transfer{
		Conn: model.CurUser.Conn,
	}
	err = tfer.WritePkg(data)
	if err != nil {
		return fmt.Errorf("changeNewPwd failed : %v", err)
	}

	res, err := tfer.ReadPkg()
	if err != nil {
		return fmt.Errorf("changeNewPwd failed : %v", err)
	}

	var resMes message.ChangePwdResMes
	err = json.Unmarshal([]byte(res.Data), &resMes)
	if err != nil {
		return fmt.Errorf("json.Unmarshal err= %v", err)
	}

	switch resMes.Code {
	case message.ChangePwdSuccessCode:
		return nil
	default:
		return errors.New("Server internal error")
	}
}

/*ChangePwd is a public function which can change user password
1. get user's old password and check it
2. transfer new password to server
*/
func (up *UserProcess) ChangePwd() {
	//1.验证旧密码
	var oldPwd string
	fmt.Printf("请输入当前密码：")
	fmt.Scanln(&oldPwd)
	ok, err := up.checkOldPwd(oldPwd)
	if err != nil {
		fmt.Printf("failed : %v\n", err)
		return
	}
	if !ok {
		fmt.Println("密码错误")
		return
	}

	//2.修改新密码
	var newPwd string
	fmt.Printf("请输入新密码：")
	fmt.Scanln(&newPwd)
	err = up.changeNewPwd(newPwd)
	if err != nil {
		fmt.Printf("密码修改失败：%v\n", err.Error())
		return
	}
	return
}
