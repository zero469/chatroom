package processor

import (
	"encoding/json"
	"fmt"
	"go_code/chapter18/project3/common/message"
	"go_code/chapter18/project3/server/utils"
	"net"
)

//SmsProcess 负责和用户相关的操作
type SmsProcess struct {
	Conn net.Conn
}

//ServerProcessSms 处理用户发送的消息
func (smsp *SmsProcess) ServerProcessSms(mes *message.Message) (err error) {
	var smsMes message.SmsMes
	err = json.Unmarshal([]byte(mes.Data))
	if err != nil {
		return fmt.Errorf("serverProcessSms failed : %v", err)
	}

	data, err := json.Marshal(mes)
	if err != nil {
		return fmt.Errorf("serverProcessSms failed : %v", err)
	}

	if smsMes.RcverID == nil {
		smsMes.RcverID = make([]int, len(UserMgr))
		for userID := range UserMgr {
			smsMes.RcverID = append(smsMes.RcverID, userID)
		}
	}

	for userID := range smsMes.RcverID {
		temp := &SmsProcess{
			Conn: UserMgr[userID].Conn,
		}
		err = temp.transferSmsMes(data)
		if err != nil {
			fmt.Println("serverProcessSms failed : %v", err)
		}
	}
	return nil

}

//
func (smsp *SmsProcess) transferSmsMes(data []byte) (err error) {
	tf := utils.Transfer{
		Conn: smsp.Conn,
	}
	err = tf.WritePkg(mes)
	if err != nil {
		return fmt.Errorf("tansferSmsMes failed : %v", err)
	}
	return
}
