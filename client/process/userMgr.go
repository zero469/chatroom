package process

import (
	"chatroom/common/message"
	"encoding/json"
	"fmt"
	"sync"
)

type UserMgr struct {
	userList map[int]*message.User
	lock     sync.RWMutex
}

var onlineUsers *UserMgr

func initUserMgr() {
	onlineUsers = &UserMgr{
		userList: make(map[int]*message.User, 16),
	}
}

func (userMgr *UserMgr) showOnlineUsers() {
	userMgr.lock.RLock()
	fmt.Printf("当前在线人数：%v\n", len(userMgr.userList))
	for id, user := range userMgr.userList {
		if user.UserState == message.UserOfflineState {
			continue
		}
		fmt.Println("* 用户id : ", id)
	}
	userMgr.lock.RUnlock()
}

func (userMgr *UserMgr) updateUserState(mes message.Message) {
	//1.反序列化解析消息结构体
	var updateMes message.UpdataUserStateMes
	err := json.Unmarshal([]byte(mes.Data), &updateMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	//2.检查该用户是否在onlineUsers中，如果不在则添加，如果在则修改其状态
	userMgr.lock.Lock()
	user, ok := userMgr.userList[updateMes.UserID]
	if !ok || user == nil {
		userMgr.userList[updateMes.UserID] = &message.User{
			UserID:    updateMes.UserID,
			UserState: updateMes.State,
		}
	} else {
		//如果该用户下线，则将其从onlineUsers中删除
		if updateMes.State == message.UserOfflineState {
			delete(onlineUsers.userList, updateMes.UserID)
		} else {
			onlineUsers.userList[updateMes.UserID].UserState = updateMes.State
		}
	}
	userMgr.lock.Unlock()
}
