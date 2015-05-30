//Copyright 2014. All rights reserved.
//Author: lja  
//Createtime: 2014/08/05 15:43:19  Lastchange: 2014/08/25 16:41:18
//changlog:  1. create by lja

package zainar

//留言板
import (
	"net"
	"time"
	"strconv"
)

var(
	//MsgBoard订阅的目标是自己，即发送给自己的消息, 只有自己一个听众
	MsgBoard map[Identi] *BdOwner
)

//留言板主人的接收连接
type BdOwner struct{
	lastmsg  int64
	conn    *net.TCPConn
}

//新建(MsgBoard)上的消息
type MsgBdNew struct{
	Mt   int8      //消息类型
	St   int16     //发送者类型
	Si   int64     //发送者ID
	Rt   int16     //接收者者类型
	Ri   int64     //接收者ID
	Bo   string     //消息体
	Ti   string     //留言时间
}

func init(){
	MsgBoard = make(map[Identi] *BdOwner)
}

//添加留言板
func AddBoard(identi Identi, conn *net.TCPConn) error{

	var bd BdOwner
	err := sqlMsgBdLast.QueryRow(MsgBdUnRd).Scan(&bd.lastmsg)
	if err != nil{
		m.Print(err)
		return err
	}

	bd.conn = conn
	MsgBoard[identi] = &bd

	return nil
}

//删除留言板
func DelBoard(identi Identi){
	MsgBoard[identi] = nil
}
//留言
func LeaveMsg(msg MsgBdNew)error {
	res,err := sqlMsgBdNew.Exec(MsgBdUnRd,msg.Mt,msg.St,msg.Si,msg.Rt,msg.Ri,msg.Bo)
	if err != nil{
		m.Print(err)
		return err
	}

	msgid,err := res.LastInsertId()
	if err != nil{
		m.Print(err)
		return err
	}

	switch msg.Rt{
	case ItemHuman:
		var q Identi
		q.ItemType = msg.Rt
		q.ItemID = msg.Ri

		MsgBdNuCast(q, int64(msgid))

		return nil
	case ItemPhyTrackTypeA:
		//TODO
		m.Print(errUnFi)
		return errUnFi
	case ItemVirtualTrackerA1:
		//TODO
		m.Print(errUnFi)
		return errUnFi
	case ItemVirtualTrackerA2:
		//TODO
		m.Print(errUnFi)
		return errUnFi
	case ItemTeam:   //如果是给Team的留言，认为是留给Team中的人类成员
		rows, err := sqlTeamMemID.Query(msg.Ri,ItemHuman)
		defer rows.Close()
		if err != nil{
			m.Print(err)
			return err
		}
		msg.Rt = ItemHuman   //注意这个必须修改，否则无限嵌套调用LeaveMsg
		for rows.Next(){
			err = rows.Scan(&msg.Ri)
			if err != nil{
				m.Print(err)
				return err
			}
			LeaveMsg(msg)
		}
		return nil
	case ItemPublic:
		//TODO
		m.Print(errUnFi)
		return errUnFi
	default:
		m.Printf("ItemType(%d) doesn't exist",msg.Rt)
	}
	return nil
}

//广播留言板中的未读消息数量
func MsgBdNuCast(q Identi, stopid int64){

	own, ok := MsgBoard[q]
	if ok != true{
		m.Print("No own")
		return
	}

	if stopid == 0{
		err := sqlMsgBdLast.QueryRow(MsgBdUnRd).Scan(&stopid)
		if err != nil{
			m.Print(err)
			return
		}
	}

	var num int64
	err := sqlMsgBdNewNu.QueryRow(MsgBdUnRd,own.lastmsg,stopid,q.ItemType,q.ItemID).Scan(&num)
	if err != nil{
		m.Print(err)
		return
	}
	if num == 0{
		return
	}

	//发送留言数量
	buf := PushMsgBdNum+strconv.FormatInt(num,10)+"\n"
	err = response(own.conn, []byte(buf))
	if err != nil{
		m.Print(err, q)
	}

	own.lastmsg = stopid
}

//只广播留言的数量，通过Cmd获取留言
func TimerMsgBdCast(d time.Duration){
	c := time.Tick(d * time.Second)
	for  _ = range c{
		for  i,owner := range MsgBoard{
			if owner == nil{
				continue
			}
			go MsgBdNuCast(i,0)
		}
	}
}
