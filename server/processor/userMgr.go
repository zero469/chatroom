package processor

import (
	"chatroom/common/message"
	"chatroom/server/model"
	"fmt"
	"net"
)

//UserMgr 在线用户管理对象
var UserMgr *userMgr

type userMgr struct {
	//userId  ClientConn
	users map[int]*model.ClientConn
}

func init() {
	UserMgr = &userMgr{
		users: make(map[int]*model.ClientConn),
	}
}

func (um *userMgr) Add(userID int, cc *model.ClientConn) {
	um.users[userID] = cc
}

func (um *userMgr) Del(userID int) {
	delete(um.users, userID)
}

func (um *userMgr) Update(userID int, cc *model.ClientConn) {
	um.Add(userID, cc)
}

func (um *userMgr) Get(userID int) (cc *model.ClientConn) {
	return um.users[userID]
}

func (um *userMgr) GetAll() (ccs []*model.ClientConn) {
	ccs = make([]*model.ClientConn, 0, len(um.users))
	for _, cc := range um.users {
		ccs = append(ccs, cc)
	}
	return
}

func (um *userMgr) GetIDbyConn(Conn net.Conn) (id int, err error) {
	for id, cc := range um.users {
		if Conn == cc.Conn {
			return id, nil
		}
	}
	return 0, fmt.Errorf("GetIDbyConn failed : conn 不存在")
}

//UpdateUserState 通知其他用户有用户状态发生该表，通知每个user使用的是up.updateUserState函数
func UpdateUserState(userID int, userStauts string) {
	//构造消息结构体
	mes := &message.UpdataUserStateMes{
		UserID: userID,
		State:  userStauts,
	}

	for id, cc := range UserMgr.users {
		if id == userID {
			continue
		}

		//构造临时userprocess对象
		(&UserProcess{
			Conn: cc.Conn,
		}).updateUserState(mes)
	}
}
