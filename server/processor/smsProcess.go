package processor

import (
	"encoding/json"
	"fmt"
	"go_code/chapter18/project3/common/message"
	"go_code/chapter18/project3/server/utils"
	"net"
	"time"
)

//SmsProcess 负责和用户相关的操作
type SmsProcess struct {
	Conn net.Conn
}

//ServerProcessSms 处理用户发送的消息
func (smsp *SmsProcess) ServerProcessSms(mes *message.Message) (err error) {
	var smsMes message.SmsMes
	err = json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		return fmt.Errorf("serverProcessSms failed : %v", err)
	}

	data, err := json.Marshal(mes)
	if err != nil {
		return fmt.Errorf("serverProcessSms failed : %v", err)
	}

	if smsMes.RcverID == nil {
		smsMes.RcverID = make([]int, 0)
		for userID := range UserMgr.users {
			smsMes.RcverID = append(smsMes.RcverID, userID)
		}
	}

	fmt.Println("recver : ", smsMes.RcverID)
	//构造转发给客户端的消息
	var smsResMes message.SmsResMes
	smsResMes.Content = smsMes.Content
	smsResMes.SendTime = time.Now().Unix()
	smsResMes.SenderID = smsMes.SenderID

	data, err = json.Marshal(smsResMes)
	if err != nil {
		return fmt.Errorf("serverProcessSms failed : %v", err)
	}

	//添加消息类型
	var resMes message.Message
	resMes.Type = message.SmsResMesType
	resMes.Data = string(data)

	data, err = json.Marshal(resMes)

	for _, userID := range smsMes.RcverID {
		temp := &SmsProcess{
			Conn: UserMgr.users[userID].Conn,
		}
		err = temp.transferSmsMes(data)
		if err != nil {
			fmt.Printf("serverProcessSms failed : %v\n", err)
		}
	}
	return nil

}

//
func (smsp *SmsProcess) transferSmsMes(data []byte) (err error) {
	tf := utils.Transfer{
		Conn: smsp.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		return fmt.Errorf("tansferSmsMes failed : %v", err)
	}
	return
}
