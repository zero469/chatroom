package processor

import (
	"go_code/chapter18/project3/common/message"
)

//UserMgr 服务器端全局在线用户管理对象
var UserMgr *onlineUsers

type onlineUsers struct {
	//key : userId | value : UserProcess
	users map[int]*UserProcess
}

func init() {
	UserMgr = &onlineUsers{
		users: make(map[int]*UserProcess),
	}
}

//增删查改接口

func (userMgr *onlineUsers) Add(userID int, up *UserProcess) {
	userMgr.users[userID] = up
}

func (userMgr *onlineUsers) Del(userID int) {
	delete(userMgr.users, userID)
}

func (userMgr *onlineUsers) Update(userID int, up *UserProcess) {
	userMgr.Add(userID, up)
}

func (userMgr *onlineUsers) Get(userID int) (up *UserProcess) {
	return userMgr.users[userID]
}

func (userMgr *onlineUsers) GetAll() (ups []*UserProcess) {
	ups = make([]*UserProcess, 0, len(userMgr.users))
	for _, up := range userMgr.users {
		ups = append(ups, up)
	}
	return
}

//UpdateUserState 通知其他用户有用户状态发生该表，通知每个user使用的是up.updateUserState函数
func UpdateUserState(userID int, userStauts string) {
	//构造消息结构体
	mes := &message.UpdataUserStateMes{
		UserID: userID,
		State:  userStauts,
	}

	for id, up := range UserMgr.users {
		if id == userID {
			continue
		}

		up.updateUserState(mes)
	}
}
