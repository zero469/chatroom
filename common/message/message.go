package message

const (
	//消息类型
	LoginMesType    = "LoginMes"
	LoginMesResType = "LoginMesRes"

	RegisterMesType    = "RegisterMesType"
	RegisterResMesType = "RegisterResMesType"

	UpdataUserStateMesType = "UpdataUserStateMesType"

	SmsMesType = "SmsMesType"
	//返回状态码
	UnRegisterCode    = 500
	LoginSuccessCode  = 200
	WrongPasswordCode = 403
	ServerErrorCode   = 505

	UserIdBeenUsedCode  = 400
	RegisterSuccessCode = 200

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
	Users []int  `json:"users"`
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
