//Copyright 2014. All rights reserved.
//Author: lja  
//Createtime: 2014/08/01 12:06:55  Lastchange: 2014/08/06 09:42:53
//changlog:  1. create by lja

package zainar

import (
	"sync"
	"database/sql"
	"net"
	"errors"
	"encoding/json"
	"time"
)

//消息队列 
//队列的添加必须要互斥,防止一个刚添加的队列被后续的添加替换了
var (
	//MsgRT订阅的目标是Team消息, 一个Team有多个听众
	MsgRealTimeQu map[Identi] *Audience
)


//观众
//听众的添加以及删除不需要互斥
type Audience struct {
	mutex sync.Mutex   //向Audience广播的的过程必须互斥,lastmsg的更新必须互斥
	lastmsg  int64
	conns  map[Identi] *net.TCPConn
}

//新建的Msg表中的实时消息
type MsgRTNew struct{
	Si    int64         //发送者ID
	Ri    int64         //接受者ID
	St    int16         //发送者类型
	Rt    int16         //接收者类型
	Mt    int8          //消息类型
	Bo    string         //消息体
}

func init(){
	MsgRealTimeQu = make(map[Identi] *Audience)
}

//life 存活时间，以分钟为单位, 0表示默认时间10080分钟(7天)
func SendMsg(msg MsgRTNew, life int16) error {
	if life == 0{
		life = 10080   //7天，10080分钟
	}

	res,err := sqlMsgNew.Exec(msg.Mt,msg.St,msg.Si,msg.Rt,msg.Ri,msg.Bo,life)
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
		//TODO
		m.Print(errUnFi)
		return errUnFi
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
	case ItemTeam:
		var q Identi
		q.ItemType = msg.Rt
		q.ItemID = msg.Ri

		au,ok := MsgRealTimeQu[q]
		if ok == false{
			m.Print("Not find Audience: ",q)
			return errNotFd
		}

		//每产生一条消息，就触发一次广播, 广播可能与定时触发的广播同时进行
		//MsgRTCast通过lastmsgid,避免消息重复
		MsgRTCast(q, au, int64(msgid))
		return nil
	case ItemPublic:
		//TODO
		m.Print(errUnFi)
		return errUnFi
	default:
		m.Print(errUnFi)
		return errUnFi
	}
	return nil
}

//q：队列标识
func AddMsgRT(q Identi) error{
	_, ok := MsgRealTimeQu[q]
	if ok == true{
		return nil //队列已经存在
	}

	//新建队列
	switch q.ItemType{
	case ItemHuman:
		return errors.New("Unfinish")
	case ItemPhyTrackTypeA:
		return errors.New("Unfinish")
	case ItemVirtualTrackerA1:
		return errors.New("Unfinish")
	case ItemVirtualTrackerA2:
		return errors.New("Unfinish")
	case ItemTeam:
		var lastmsg int64
		//err := sqlMsgLastRec.QueryRow(q.ItemType,q.ItemID).Scan(&lastmsg)
		err := sqlMsgLast.QueryRow().Scan(&lastmsg)
		if err == sql.ErrNoRows{
			err := sqlMsgLast.QueryRow().Scan(&lastmsg)
			if err == sql.ErrNoRows{
				lastmsg = 0
			}else if err != nil{
				m.Print(err)
				return err
			}
		}else if err != nil{
				m.Print(err)
				return err
		}
		MsgRealTimeQu[q] = &Audience{lastmsg:lastmsg, conns: make(map[Identi]*net.TCPConn)}
		return nil
	case ItemPublic:
		return errors.New("Unfinish")
	default:
		return errors.New("doesn't exist")
	}
}

func DelMsgRT(identi Identi) error{
	MsgRealTimeQu[identi] = nil
	return nil
}

//q: 队列标识, 如果队列不存在，新建
//a: 听众标识 
//conn: 听众的接收地址
func AddAudience(q Identi, a Identi, conn *net.TCPConn) error{
	au, ok := MsgRealTimeQu[q]
	if ok == false{
		err := AddMsgRT(q)
		if err != nil{
			return err
		}
		au, ok = MsgRealTimeQu[q]
		if ok == false{
			return errors.New("go's map error!")
		}
	}
	au.conns[a] = conn
	return nil
}

//q: 队列标识
//a: 听众标识 
func DelAudience(q Identi, a Identi) error{
	au, ok := MsgRealTimeQu[q]
	if ok == false{
		return errors.New("queue doesn't exist")
	}
	au.conns[a] = nil
	return nil
}

//消息广播,sec 定时调度时间
//程序启动时,将其运行在一个goroutine中
func TimerMsgRTCast(d time.Duration){
	c := time.Tick(d * time.Second)
	for  _ = range c{
		for  i,au := range MsgRealTimeQu{
			if au == nil{
				continue
			}
			go MsgRTCast(i,au,0)
		}
	}
}

//将目标是rid,rtype的消息广播给Audience
//stopid: 查询时的截止id, 如果为0，需要查询获得
//au: Audience, 如果为nil, 查找q对应的Audience
func MsgRTCast(q Identi, au *Audience, stopid int64){
	defer au.mutex.Unlock()
	au.mutex.Lock()
	if stopid == 0{
		err := sqlMsgLast.QueryRow().Scan(&stopid)
		if err != nil{
			m.Print(err)
			return
		}
	}

	var ok bool
	if au == nil{
		au, ok = MsgRealTimeQu[q]
		if ok != true{
			m.Print("No Audience")
			return
		}
	}

	rows,err := sqlMsg.Query(au.lastmsg,stopid,q.ItemType,q.ItemID)
	defer rows.Close()
	if err != nil{
		m.Print(err)
		return
	}

	buf := []byte(PushMsgRT)
	for rows.Next(){
		var msg MsgRT
		err := rows.Scan(&msg.Mt,&msg.St,&msg.Si,&msg.Rt,&msg.Ri,&msg.Bo)
		if err != nil{
			m.Print(err)
			return
		}
		msg.Rt = q.ItemType
		msg.Ri = q.ItemID

		j, err := json.Marshal(msg)
		if err != nil{
			m.Print(err)
			return
		}
		buf = AppendSlice(buf,j)
	}

	if len(buf) == len(PushMsgRT){
		return
	}

	err = rows.Err()
	if err != nil{
		m.Print(err)
		return
	}

	buf = append(buf,'\n')

	for identi,conn := range au.conns{
		if conn == nil{
			continue
		}
		err = response(conn, buf)
		if err != nil{
			m.Print(err, identi)
		}
	}
	au.lastmsg = stopid
}
