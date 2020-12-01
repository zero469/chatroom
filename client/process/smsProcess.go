package process

import (
	"chatroom/client/model"
	"chatroom/client/utils"
	"chatroom/common/message"
	"encoding/json"
	"fmt"
)

type SmsProcess struct {
}

func (smsp *SmsProcess) SendGroupMes(content string) (err error) {
	var mes message.Message
	mes.Type = message.SmsMesType

	var smsMes message.SmsMes
	smsMes.SenderID = model.CurUser.UserID
	smsMes.Content = content

	data, err := json.Marshal(smsMes)
	if err != nil {
		return fmt.Errorf("sendGroupMes failed : %v", err)
	}

	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		return fmt.Errorf("sendGroupMes failed : %v", err)
	}

	tf := utils.Transfer{
		Conn: model.CurUser.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		return fmt.Errorf("sendGroupMes failed : %v", err)
	}
	return nil
}

//将服务器发送的用户消息添加到历史消息列表中
func rcvSmsMes(mes message.Message) (err error) {
	var smsResMes message.SmsResMes
	err = json.Unmarshal([]byte(mes.Data), &smsResMes)
	if err != nil {
		return fmt.Errorf("rcvSmsMes failed : %v", err)
	}
	HistoryMes.AddMes(userMes{
		sender:   smsResMes.SenderID,
		sendTime: smsResMes.SendTime,
		content:  smsResMes.Content,
	})
	return nil
}
