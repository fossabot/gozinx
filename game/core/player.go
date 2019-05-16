/**
  * Author: JeffreyBool
  * Date: 2019/5/17
  * Time: 01:39
  * Software: GoLand
*/

package core

import (
	"github.com/JeffreyBool/gozinx/src/ziface"
	"sync/atomic"
	"math/rand"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

var PidGen uint32 = 1 //用来生产玩家的 Id计数器

//玩家对象
type Player struct {
	Pid  uint32             //玩家 Id
	Conn ziface.IConnection //当前玩家的连接(用于和客户端的连接)
	X    float32            //平面的 X 坐标
	Y    float32            //高度
	Z    float32            //平面 Y 坐标
	V    float32            //旋转的 0-360 角度
}

//创建玩家方法
func NewPlayer(conn ziface.IConnection) *Player {
	//生成玩家的 Id
	atomic.AddUint32(&PidGen, 1)

	//创建玩家对象
	return &Player{
		Pid:  PidGen,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)), //随机在160坐标点 基于X轴偏移若干坐标
		Y:    0,                            //高度为0
		Z:    float32(134 + rand.Intn(17)), //随机在134坐标点 基于Y轴偏移若干坐标
		V:    0,                            //角度为0，尚未实现
	}
}

//告知客户端pid,同步已经生成的玩家ID给客户端

//广播玩家自己的出生地点

/*
	发送消息给客户端，
	主要是将pb的protobuf数据序列化之后发送
*/
func (p *Player) SendMsg(msgId uint32, data proto.Message) error {
	//将proto Message结构体序列化
	msg, err := proto.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "player to sendMsg marshal error")
	}

	//判断客户端是否已经离线
	if p.Conn == nil {
		return errors.New("connection in player is nil")
	}

	//调用Zinx框架的SendMsg发包
	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		return errors.New("player sendMsg error")
	}
	return nil
}
