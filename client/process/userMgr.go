package process

import (
	"fmt"
	"go_code/chapter18/project3/common/message"
)

var onlineUsers map[int]*message.User = make(map[int]*message.User, 16)

func showOnlineUsers() {
	for id := range onlineUsers {
		if id == myId {
			continue
		}
		fmt.Println("用户id : ", id)
	}
}
