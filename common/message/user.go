package message

//TODO:单独打一个package

type User struct {
	UserID    int    `json:"userid"`
	UserPwd   string `json:"userpwd"`
	UserName  string `json:"username"`
	UserState string `json:"userstate"`
}
