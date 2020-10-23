package model

import (
	"encoding/json"
	"fmt"

	"github.com/redigo/redis"
)

//UserDao User Database Access Object
type UserDao struct {
	pool *redis.Pool
}

var (
	MyUserDao *UserDao
)

//NewUserDao 工厂模式创建结构体
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

//根据用户ID返回User实例
func (dao *UserDao) getUserByID(conn redis.Conn, id int) (user *User, err error) {
	res, err := redis.String(conn.Do("HGET", "users", id))
	if err != nil {
		if err == redis.ErrNil {
			//没找到
			err = ERROR_USER_NOTEXISTS
		}
		return
	}

	user = &User{}
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return
}

//Login 完成对用户信息的校验
func (dao *UserDao) Login(userID int, userPwd string) (user *User, err error) {
	conn := dao.pool.Get()
	defer conn.Close()
	user, err = dao.getUserByID(conn, userID)
	if err != nil {
		return
	}
	if userPwd != user.UserPwd {
		return nil, ERROR_USER_PWD
	}
	return user, nil
}
