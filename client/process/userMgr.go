package process

import (
	"encoding/json"
	"fmt"
	"go_code/chapter18/project3/common/message"
)

var onlineUsers map[int]*message.User = make(map[int]*message.User, 16)

func showOnlineUsers() {
	for id := range onlineUsers {
		if id == myID {
			continue
		}
		fmt.Println("用户id : ", id)
	}
}

func updateUserState(mes message.Message) {
	//1.反序列化解析消息结构体
	var updateMes message.UpdataUserStateMes
	err := json.Unmarshal([]byte(mes.Data), &updateMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	//2.检查该用户是否在onlineUsers中，如果不在则添加，如果在则修改其状态
	if onlineUsers[updateMes.UserID] == nil {
		onlineUsers[updateMes.UserID] = &message.User{
			UserID:    updateMes.UserID,
			UserState: updateMes.State,
		}
	} else {
		onlineUsers[updateMes.UserID].UserState = updateMes.State
	}

}
