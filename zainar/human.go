//Copyright 2014. All rights reserved.
//Author: lja  
//Createtime: 2014/07/29 06:07:46  Lastchange: 2014/08/13 13:33:24
//changlog:  1. create by lja

package zainar

import (
	"net/http"
	"crypto/rand"
	"database/sql"
	"math/big"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
)

type HumanInfo struct{
	ID    int64
	SI    int64
	CP    int8
	CW    int32
	SP    int8
	SW    int32
	NI    string    //NickName
}

func HumanLogin(w http.ResponseWriter, r *http.Request, pwd, mail string,
	itemtype int16) error{

	var wait string
	err := sqlHumWait.QueryRow(mail).Scan(&wait)
	switch {
	case err == sql.ErrNoRows:    //Human不存在
		//TODO: Log 统计登陆尝试
		_,err = w.Write([]byte(LogNotExist))
		if err != nil{
			l.Print(err)
			return err
		}
		return err
	case err != nil:
		l.Print(err)
		return err
	}

	if wait[0] != '-' {           //但是已经锁定, 需要等待wait长的时间
		_, err = w.Write([]byte(LogLock+wait+"\n"))
		if err != nil{
			l.Print(err)
			return err
		}
		return nil
	}

	var humanid int64
	err = sqlHumMailUsed.QueryRow(mail).Scan(&humanid)
	switch {
	case err == sql.ErrNoRows:    //Human不存在
		//TODO: Log 统计登陆尝试
		_,err = w.Write([]byte(LogNotExist))
		if err != nil{
			l.Print(err)
			return err
		}
		return err
	case err != nil:
		l.Print(err)
		return err
	}

	//Human存在，检查是否已经有Session
	var sid int64
	var loginok int8
	var loginerr int8

	err = sqlSesTry.QueryRow(itemtype, humanid).Scan(&sid,&loginok,&loginerr)
	switch {
	case err == sql.ErrNoRows:    //Session不存在, 创建
		res, err := sqlSesNew.Exec(humanid,itemtype)
		xsid, err := res.LastInsertId()
		sid = int64(xsid)
		if err != nil{
			l.Print(err)
			return err
		}
		loginok = 0;
		loginerr = 0;
	case err != nil:
		l.Print(err)
		return err
	}
/*
	//在这里可以对重复登陆情况处理，这段代码注释掉，表示目前不做处理
	if loginok == 1 {    //重复登陆
		http.Redirect(w, r, Uri_repeat_login, 307)
		return nil
	}
*/
	if loginerr == LogTryNum {  //登陆错误达到五次
		_, err =w.Write([]byte(LogTryMax))
		if err != nil{
			l.Print(err)
			return err
		}
		return nil
	}

	//验证密码
	var nick string
	var realmail int8

	err = sqlHumChkPwd.QueryRow(humanid,AddSalt(pwd,mail)).Scan(&nick,&realmail)

	switch  {
	case err == sql.ErrNoRows:   //密码错误
		loginerr++

		if loginerr == LogTryNum{    //错误达到五次, 锁定用户，删除会话
			_,err = sqlHumLock.Exec(humanid)    //锁定用户
			if err != nil{
				l.Print(err)
				return err
			}
			_,err = sqlSesDel.Exec(sid)    //删除会话
			if err != nil{
				l.Print(err)
				return err
			}
			_, err =w.Write([]byte(LogTryMax))
			if err != nil{
				l.Print(err)
				return err
			}
			return nil
		}
		//不足五次
		_,err = sqlSesUpErr.Exec(loginerr,sid)
		if err != nil{
			l.Print(err)
			return err
		}

		_, err := w.Write([]byte(LogWrong))
		if err != nil{
			l.Print(err)
			return err
		}
		return nil
	case err != nil:
		l.Print(err)
		return err

	default:     //密码正确
		//登陆成功
		random, err := rand.Int(rand.Reader, big.NewInt(127))
		if err != nil{
			_, _ = sqlSesDel.Exec(sid)
			l.Print(err)
			return err
		}
		cp := int8(random.Int64())

		random, err = rand.Int(rand.Reader, big.NewInt(127))
		if err != nil {
			_, _ = sqlSesDel.Exec(sid)
			l.Print(err)
			return err
		}
		sp := int8(random.Int64())

		random, err = rand.Int(rand.Reader, big.NewInt(2147483647))
		if err != nil {
			_, _ = sqlSesDel.Exec(sid)
			l.Print(err)
			return err
		}
		cw := int32(random.Int64())

		random,err = rand.Int(rand.Reader, big.NewInt(2147483647))
		if err != nil {
			_, _ = sqlSesDel.Exec(sid)
			l.Print(err)
			return err
		}
		sw := int32(random.Int64())
		
		_,err = sqlSesUpOk.Exec(cw,cp,sw,sp,sid)

		if err != nil{
			_, _ = sqlSesDel.Exec(sid)
			l.Print(err)
			return err
		}

		humaninfo := HumanInfo{NI: nick, ID: humanid, SI:sid, CP: cp, CW: cw, SP: sp, SW: sw}

		jsonline, err := json.Marshal(humaninfo)
		if err != nil{
			l.Print(err)
			return err
		}
		_, err = w.Write(jsonline)
		if err != nil{
			l.Print(err)
			return err
		}
	}
	return nil
}

func HumanRegist(w http.ResponseWriter, r *http.Request, pwd,mail,nick string,
	itemtype int16) error{
	var humanid int64;

	err := sqlHumMailUsed.QueryRow(mail).Scan(&humanid)
	if err == nil{
		_, err = w.Write([]byte(RegMailUsed))
		return nil
	}else if err != sql.ErrNoRows{
		l.Print(err)
		return err
	}

	err = sqlHumNickUsed.QueryRow(nick).Scan(&humanid)
	if err == nil{
		_, err = w.Write([]byte(RegNickUsed))
		return nil
	}else if err != sql.ErrNoRows{
		l.Print(err)
		return err
	}

	res,err := sqlHumNew.Exec(mail,AddSalt(pwd,mail),nick)
	if err != nil {
		l.Print(err)
		return err
	}

	tmpid, err := res.LastInsertId();
	humanid = int64(tmpid)
	if err != nil {
		l.Print(err)
		return err
	}

	_, err = sqlHumStatNew.Exec(humanid)
	if err != nil {
		l.Print(err)
		return err
	}

	var humaninfo = HumanInfo{NI: nick, ID: int64(humanid)}

	jsonline, err := json.Marshal(humaninfo)
	if err != nil{
		l.Print(err)
		return err
	}
	_, err = w.Write(jsonline)
	if err != nil{
		l.Print(err)
		return err
	}

	return nil
}
