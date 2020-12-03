package model

import (
	"net"
)

type ClientConn struct {
	UserName string
	Conn     net.Conn
}
