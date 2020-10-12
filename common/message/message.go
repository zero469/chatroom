package message

const (
	//消息类型
	LoginMesType    = "LoginMes"
	LoginMesResType = "LoginMesRes"
	//返回状态码
	UnRegisterCode = 500
	LoginSuccess   = 200
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
}
