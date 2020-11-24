package process

import (
	"fmt"
	"time"
)

//用户消息类型，用于在客户端保存历史信息
type userMes struct {
	sender   int
	sendTime int64
	content  string
}

//消息列表，管理历史消息
type MesMgr struct {
	mesList []userMes
}

var HistoryMes MesMgr

func initMesMgr() {
	HistoryMes.mesList = make([]userMes, 128)
}

func (mesMgr MesMgr) AddMes(mes userMes) {
	mesMgr.mesList = append(mesMgr.mesList, mes)
}

//显示消息列表
func (mesMgr MesMgr) ShowMesList() {
	for _, v := range mesMgr.mesList {
		tm := time.Unix(v.sendTime, 0)
		fmt.Printf("[%v]%v : %v\n", tm.Format("01-02 15:04:05"), v.sender, v.content)
	}
}
