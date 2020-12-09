# Golang 聊天室项目
## 项目简介
本项目为使用Golang和redis实现的Socket聊天室，主要功能如下：
1. 用户注册、登录
2. 在线用户列表
3. 群聊
4. 历史消息列表
## 目录结构
```
.
├── README.md
├── client                              //客户端
│   ├── main                            
│   │   └── main.go                     //主函数
│   ├── model                           //model层
│   │   └── curuser.go 
│   ├── process                         //处理具体的事件
│   │   ├── mesMgr.go                   //历史消息管理
│   │   ├── server.go                   //与服务器交互
│   │   ├── smsProcess.go               //处理用户聊天发送的消息
│   │   ├── userMgr.go                  //用户信息管理
│   │   └── userProcess.go              //注册、登录
│   └── utils                           //工具函数
│       └── utils.go
├── common                              //服务器客户端共有代码
│   └── message
│       ├── message.go                  //消息类型和状态码
│       └── user.go
└── server                              //服务器
    ├── main                            //main包
    │   ├── main.go                     //主函数
    │   ├── processor.go                //
    │   └── redis.go                    //维护一个redis连接池
    ├── model                           //model层
    │   ├── error.go            
    │   ├── user.go             
    │   └── userDao.go
    ├── processor                       //处理具体的事件
    │   ├── smsProcess.go
    │   ├── userMgr.go                  //用户管理
    │   └── userProcess.go              //注册、登录、用户状态维护
    └── utils                           //工具函数
        └── utils.go
```

## 如何运行（windows）
### 前置要求
1. Go编译环境
2. redis
3. redigo包
### 下载代码
```
cd ${GOPATH}/src
git clone https://github.com/zero469/chatroom.git
```
### 编译运行
```
//服务器
go build -o server.exe chatroom/server/main
server.exe 
//客户端
go build -o client.exe chatroom/client/main
client.exe
```
---
# TODO:
- [x] 重构客户端userMgr，增加互斥锁
- [ ] 点对点聊天
- [x] 在线用户能显示用户Name
- [ ] 修改用户名和密码 //serverProcessMes已经在监听服务器发送的消息，如何将监听的权限交给ChangePwd