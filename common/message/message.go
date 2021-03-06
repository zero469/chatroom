package message

const (
	//消息类型
	LoginMesType    = "LoginMes"
	LoginMesResType = "LoginMesRes"

	RegisterMesType    = "RegisterMesType"
	RegisterResMesType = "RegisterResMesType"

	UpdataUserStateMesType = "UpdataUserStateMesType"

	SmsMesType    = "SmsMesType"
	SmsResMesType = "SmsResMesType"

	CheckOldPwdMesType  = "CheckOldPwdMesType"
	ChangeNewPwdMesType = "ChangeNewPwdMesType"
	ChangePwdResMesType = "ChangePwdResMesType"

	//返回状态码
	UnRegisterCode     = 500
	WrongPasswordCode  = 403
	ServerErrorCode    = 505
	UserIdBeenUsedCode = 400

	LoginSuccessCode       = 200
	RegisterSuccessCode    = 201
	CheckOldPwdSuccessCode = 202
	ChangePwdSuccessCode   = 203

	//用户状态
	UserOnlineState  = "UserOnlineState"
	UserOfflineState = "UserOfflineState"
)

//Message 为客户端和服务器公用的消息结构
type Message struct {
	Type string `json:"type"` //消息类型
	Data string `json:"data"` //消息内容
}

//LoginMes 的Type为 LoginMesType
type LoginMes struct {
	UserId   int    `json:"userid"`   //用户id
	UserPwd  string `json:"userpwd"`  //用户密码
	UserName string `json:"username"` //用户名
}

type LoginResMes struct {
	Code  int    `json:"code"`  //返回状态码
	Error string `json:"error"` //错误信息
	Users []User `json:"users"`
}

type RegisterMes struct {
	User User `json:"user"`
}

type RegisterResMes struct {
	Code  int    `json:"code"`  //返回状态码
	Error string `json:"error"` //错误信息
}

type UpdataUserStateMes struct {
	UserID int    `json:"userid"`
	State  string `json:"state"`
}

//SmsMes 消息类型
type SmsMes struct {
	SenderID int    `json:"senderid"`
	Content  string `json:"content"`
	RcverID  []int  `json:"rcvid"` //由发送发指定，如果为空表示群发
}

//SmsResMes 服务器转发的消息，需要附带发送的时间戳
type SmsResMes struct {
	SenderID int    `json:"senderid"`
	Content  string `json:"content"`
	SendTime int64  `json:"sendtime"` //unix时间戳，由服务器获取
}

type CheckOldPwdMes struct {
	ID     int    `json:"id"`
	OldPwd string `json:"oldpwd"`
}

type ChangeNewPwdMes struct {
	ID     int    `json:"id"`
	NewPwd string `json:"newPwd"`
}

type ChangePwdResMes struct {
	Code int `json:"code"`
}
