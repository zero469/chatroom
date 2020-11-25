package process

import (
	"fmt"
	"sync"
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
	lock    sync.RWMutex
}

var HistoryMes *MesMgr

func initMesMgr() {
	HistoryMes = &MesMgr{
		mesList: make([]userMes, 0, 128),
	}
}

func (mesMgr *MesMgr) AddMes(mes userMes) {
	mesMgr.lock.Lock()
	mesMgr.mesList = append(mesMgr.mesList, mes)
	mesMgr.lock.Unlock()
}

//显示消息列表
func (mesMgr *MesMgr) ShowMesList() {
	mesMgr.lock.RLock()
	for _, v := range mesMgr.mesList {
		tm := time.Unix(v.sendTime, 0)
		fmt.Printf("+ [%v] %v: %v\n", tm.Format("15:04:05"), v.sender, v.content)
	}
	mesMgr.lock.RUnlock()
}
