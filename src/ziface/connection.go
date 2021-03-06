/**
  * Author: JeffreyBool
  * Date: 2019/5/15
  * Time: 02:20
  * Software: GoLand
*/

package ziface

import (
	"net"
)

type IConnection interface {
	//启动链接 让当前的链接准备工作
	Start()

	//停止链接 结束当前链接的工作
	Stop()

	//获取当前链接绑定的 socket conn
	GetTCPConnection() *net.TCPConn

	//获取当前链接模块的链接 ID
	GetConnID() uint32

	//获取远程客户端的 TCP状态 IP Port
	RemoteAddr() net.Addr

	//发送数据， 将数据发送给客户端
	SendMsg(Id uint32, data []byte) error

	//设置连接属性
	SetProperty(key string, value interface{})

	//获取连接属性
	GetProperty(key string) (value interface{}, err error)

	//移除连接属性
	RemoveProperty(key string)
}

type ConnFunc func(conn IConnection)
