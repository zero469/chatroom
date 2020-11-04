package model

import (
	"net"
)

//CurUser 当前在线用户
var CurUser curUser

type curUser struct {
	UserID int
	Conn   net.Conn
}

func InitCurUser(userID int, conn net.Conn) {
	CurUser.UserID = userID
	CurUser.Conn = conn
}
