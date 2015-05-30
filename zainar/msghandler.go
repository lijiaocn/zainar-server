//Copyright 2014. All rights reserved.
//Author: lja  
//Createtime: 2014/07/30 05:24:03  Lastchange: 2014/08/26 06:44:06
//changlog:  1. create by lja

package zainar

import (
	"net"
	"database/sql"
	"encoding/json"
	"bufio"
	"errors"
	"strings"
	"io"
	"time"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

//C->S的消息头
//消息格式: json表示的消息头,后面紧邻数据
type MsgHdr struct{
	Si  int64    //SessionID
	Nu  int32    //消息编号
	Cm  string    //指令
	Da  int       //数据大小
}

type SessionInfo struct{
	team     int64     //当前所在teamid
	itemid   int64
	itemtype int16
	cw       int32
	sw       int32
	cp       int8
	sp       int8
	conn     *net.TCPConn      //不要用这个conn进行读写操作，这个只被用来作为值传递
}

func MsgHandler(conn *net.TCPConn) {
	defer conn.Close()
	addbd := true    //添加留言板
	remote := conn.RemoteAddr()
	local  := conn.LocalAddr()
	rd := bufio.NewReader(conn)
	for {
		jsonline,err := rd.ReadString('\n')
		//---------------Test Start Time
		st := time.Now()
		//---------------
		if err != nil{
			m.Printf("%s->%s():%s",remote,local,err)
			response(conn, []byte(CmdReconnect))
			return
		}
		if len(jsonline) == 1{  //只有一个换行符
			continue
		}
		dec := json.NewDecoder(strings.NewReader(jsonline))
		var hdr MsgHdr
		err = dec.Decode(&hdr)
		if err != nil && err != io.EOF{
			m.Printf("%s->%s():%s",remote,local,err)
			response(conn, []byte(CmdReconnect))
			return
		}
		//---------------Test End Time
		ns := time.Now().Sub(st).Nanoseconds()
		fmt.Printf("%dns (%dms) begin find session \n",ns,ns/1000000)
		//---------------
		s := SessionInfo{conn:conn}
		err = sqlSesInfo.QueryRow(hdr.Si).
			Scan(&s.itemid,&s.itemtype,&s.team,&s.cw,&s.cp,&s.sw,&s.sp)
		//err = SidStmt.QueryRow(hdr.Si).Scan(&s.itemid,&s.itemtype,&s.team,
		//	&s.cw,&s.cp,&s.sw, &s.sp)
		switch {
		case err == sql.ErrNoRows:  //Session不存在, 有人冒充
			m.Printf("%s->%s(sid:%d,num:%d):%s",
				remote,local,hdr.Si,hdr.Nu,"Sid doesn't exist")
			response(conn, []byte(CmdReLogin))
			return
		case err != nil:   //内部错误
			m.Printf("%s->%s(sid:%d,num:%d):%s",
				remote,local,hdr.Si,hdr.Nu,err)
			response(conn, []byte(CmdReLogin))
			return
		}

		if addbd == true{
			identi := Identi{ItemType:s.itemtype, ItemID:s.itemid}
			err = AddBoard(identi,conn)
			if err == nil{
				addbd = false   //添加成功
			}else{
				m.Print(err)
			}
		}

		//---------------Test End Time
		ns = time.Now().Sub(st).Nanoseconds()
		fmt.Printf("%dns (%dms) begin check pkt num \n",ns,ns/1000000)
		//---------------

		//Msg Num错误，相当于使用了错误的Cookie, 当前连接已经不可信
		if  hdr.Nu != s.sw {
			m.Printf("%s->%s(sid:%d,num:%d):%s",
				remote,local,hdr.Si,hdr.Nu,"Msg Num Wrong")
			response(conn, []byte(CmdReLogin))
			return
		}

		//---------------Test End Time
		ns = time.Now().Sub(st).Nanoseconds()
		fmt.Printf("%dns (%dms) begin up pkt num \n",ns,ns/1000000)
		//---------------

		//Msg Num正确，消耗掉这个Num
		_,err = sqlSesNewSW.Exec(s.sw+int32(s.sp), hdr.Si)
		if err != nil{
			m.Printf("%s->%s(sid:%d,num:%d):%s",remote,local,hdr.Si,hdr.Nu,err)
			response(conn, []byte(CmdReLogin))
			return
		}

		//---------------Test End Time
		ns = time.Now().Sub(st).Nanoseconds()
		fmt.Printf("%dns (%dms) begin cmd find\n",ns,ns/1000000)
		//---------------

		//cmd处理
		cmdhandler, ok := CmdHandlers[hdr.Cm]
		if ok == false{
			m.Printf("%s->%s(sid:%d,num:%d):%s",
				remote,local,hdr.Si,hdr.Nu,"Don't have this cmd,%s",hdr.Cm)
			return
		}

		//---------------Test End Time
		ns = time.Now().Sub(st).Nanoseconds()
		fmt.Printf("%dns (%dms) begin cmd deal\n",ns,ns/1000000)
		//---------------

		buf, err := cmdhandler(rd, hdr, s)
		if err == ErrFatalErr{   //FatalErr关闭连接
			response(conn, []byte(buf))
			return
		}

		//---------------Test End Time
		ns = time.Now().Sub(st).Nanoseconds()
		fmt.Printf("%dns (%dms) begin resp\n",ns,ns/1000000)
		//---------------

		response(conn, []byte(buf+"\n"))
		fmt.Printf("%s\n",buf)

		//---------------Test End Time
		ns = time.Now().Sub(st).Nanoseconds()
		fmt.Printf("%dns (%dms) finish\n",ns,ns/1000000)
		//---------------
	}

//		//带有数据
//		if hdr.Da > 65536 {
//			m.Print(remote, "->", local, "data size err: ", hdr.Da)
//			response(conn, []byte(CmdReconnect))
//			return
//		}
//
//		buf, err := getdat(pkt, hdr.Da)
//		if err != nil {
//			m.Print(remote, "->", local, "data size err: ", err)
//			response(conn, []byte(CmdReconnect))
//			return
//		}
//
//		//TODO 处理数据
//		fmt.Print(buf, "\n")
//		response(conn, []byte(CmdOK))
}

func response(conn *net.TCPConn, res []byte) error{
	//TODO: 附带上Client端的PktWait
	if conn == nil{
		m.Print("conn is nil, must check it\n")
		return nil
	}
	_, err := conn.Write(res)
	if err != nil{
		m.Print("response: ", err)
		conn.Write([]byte(CmdReconnect))
		return err
	}
	return nil
}

func getdat(rd *bufio.Reader, total int) ([]byte, error){
	buf := make([]byte, total)
	n, err := rd.Read(buf)
	if err != nil {
		m.Print("Conn Read err")
		return nil,err
	}
	for n < total{
		r, err := rd.Read(buf[n:])
		if err != nil{
			m.Print("Conn Read err")
			return nil,err
		}
		n = n + r
		if n > total{
			return nil, errors.New("real data size is larger")
		}
	}
	return buf,nil
}
