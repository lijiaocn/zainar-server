//Copyright 2014. All rights reserved.
//Author: lja  
//Createtime: 2014/07/31 05:11:16  Lastchange: 2014/08/27 07:37:23
//changlog:  1. create by lja

package zainar

import (
	"strconv"
	"bufio"
	"errors"
	"encoding/json"
	"strings"
	"database/sql"
	"log"
)

var (
	CmdHandlers = map[string]CmdHandler{
		"TL":      cmdTeamList,     //获取我的团队列表, Team List
		"NT":      cmdNewTeam,      //创建团队, New Team
		"FT":      cmdFindTeam,     //查找团队, Find Team
		"JT":      cmdJoinTeam,     //加入团队, Join Team
		"RJ":      cmdRejectJoin,   //拒绝加入
		"IT":      cmdEnterTeam,    //进入团队, In Team
		"UP":      cmdUpPos,        //更新位置, Update 
		"LT":      cmdLeaveTeam,    //离开团队, Leave Team
		"ET":      cmdExitTeam,     //退出团队, Out Team
		"GM":      cmdGetMsgBd,     //获取留言板上的消息  
		"LM":      cmdLookMsgBd,    //查看留言板留言数量
		"LO":      cmdLogout,       //退出系统, Logout 
		"SM":      cmdSetMsgBd,     //设置留言状态
		"AM":      cmdAddMem,       //添加新成员
		"DS":      cmdDisTeam,      //解散Team
		"TS":      cmdTeamSensi,    //Team的敏感信息
		"TU":      cmdUpdateTeam,   //更新Team
	}

	//返回ErrFatal, 表示严重错误，连接需要被关闭
	ErrFatalErr = errors.New("Fatal Error, close connection")
)

type CmdHandler func(rd *bufio.Reader, hdr MsgHdr, s SessionInfo)(string, error)

type RejectJoin struct{
	Tp int16    //Item Type
	ID int64    //拒绝的Human的ID
	TID int64   //TeamID
	De string   //拒绝原因
}

type DisTeam struct{
	ID int64   //要解散的Team的ID
}

type MemAdd struct{
	Rl int16   //角色
	Tp int16   //Mem类型
	ID int64   //MemID
	TI int64   //目标TeamID
	De string  //Mem描述
	Tg string  //Mem标签
}

type SetMsgBd struct{
	ID  int64
	Rd  int8      //MsgBdUnRd,MsgBdRd, 已读，未读, 以后可以增加收藏等。。
}

type GetMsgBd struct{
	ID   int64     //消息的ID要大于这个ID
	To   int8      //返回记录总数
}

type LookMsgBd struct{
	ID   int64     //消息的ID要大于这个ID
}

type Pos struct{
	Lo  float32  //经度
	La  float32  //纬度
	He  float32  //海拔
}

type NewTeam struct{
	Na      string    //name
	De      string    //Descripe
	Tg      string    //tag
	Pu      int8     //公开性
	Aw      int8     //Allow加入许可
	Pr      string     //加入码提示语, Prompt
	Co      string     //加入码,code
}

type UpdateTeam struct{
	ID      int64     //TeamID
	Na      string    //name
	De      string    //Descripe
	Tg      string    //tag
	Pu      int8     //公开性
	Aw      int8     //Allow加入许可
	Pr      string     //加入码提示语, Prompt
	Co      string     //加入码,code
}

type FindTeam struct{
	Wd    string    //Team名称关键字
	Tg   string     //Team标签
	Nu   int16     //读取查询结果的第Num到Num+Tot条记录
	To   int8     //返回记录总数
}

type JoinTeam struct{
	ID   int64     //要加入的Team的ID
	Co   string     //JoinCoe
	De   string     //Descripe加入说明
}

type ExitTeam struct{
	ID  int64     //要加入的Team的ID
}

type TeamIdenti struct{
	ID int64
}

//进入Team
type EnterTeam struct{
	ID int64    //要进入的Team的ID
}


//更新Team状态
func cmdUpdateTeam(rd *bufio.Reader, hdr MsgHdr, s SessionInfo)(string, error){
	jsonline,err := rd.ReadString('\n')
	if err != nil{
		m.Print(err)
		return CmdNetFault, ErrFatalErr
	}

	var t UpdateTeam
	dec := json.NewDecoder(strings.NewReader(jsonline))
	err = dec.Decode(&t)
	if err != nil {
		m.Print(err)    //MsgHdr解析出错
		return CmdFormatErr, nil
	}

	var rol int16
	err = sqlTeamRol.QueryRow(t.ID,s.itemtype,s.itemid).Scan(&rol)
	if err != nil{
		m.Print(err)
		return CmdDbErr,err
	}

	//权限检查,创建者才允许
	if rol != TeamOwer {
		return CmdNotAllow,nil
	}

	_,err = sqlTeamUpdate.Exec(t.Na,t.De,t.Tg,t.Pu,t.Aw,t.Pr,t.Co,t.ID)
	if err != nil{
		m.Print(err)
		return CmdDbErr,err
	}

	fname := strconv.FormatInt(int64(ItemTeam),10)+"-"+strconv.FormatInt(t.ID,10);
	err = WriteStaticFile(fileTeamPubInfo, fname, []byte(jsonline))
	if err != nil{
		m.Print(err)
		return InternalErr,nil
	}

	return CmdOK,nil
}

//获取Team的敏感信息,例如JoinCode
func cmdTeamSensi(rd *bufio.Reader, hdr MsgHdr, s SessionInfo)(string, error){
	jsonline,err := rd.ReadString('\n')
	if err != nil{
		m.Print(err)
		return CmdNetFault, ErrFatalErr
	}

	var t TeamIdenti
	dec := json.NewDecoder(strings.NewReader(jsonline))
	err = dec.Decode(&t)
	if err != nil {
		m.Print(err)    //MsgHdr解析出错
		return CmdFormatErr, nil
	}

	var rol int16
	err = sqlTeamRol.QueryRow(t.ID,s.itemtype,s.itemid).Scan(&rol)
	if err != nil{
		m.Print(err)
		return CmdDbErr,err
	}

	//权限检查,创建者才允许
	if rol != TeamOwer {
		return CmdNotAllow,nil
	}

	var ts TeamSensi
	ts.ID = t.ID
	err = sqlTeamSensi.QueryRow(t.ID).Scan(&ts.Pr,&ts.Co)
	if err != nil{
		m.Print(err)
		return CmdDbErr,err
	}

	j, err := json.Marshal(ts)
	if err != nil{
		m.Print(err)
		return InternalErr, nil
	}

	return CmdTeamSensiOK+string(j)+"\n",nil
}

//拒绝加入的申请
func cmdRejectJoin(rd *bufio.Reader, hdr MsgHdr, s SessionInfo)(string, error){

	jsonline,err := rd.ReadString('\n')
	if err != nil{
		m.Print(err)
		return CmdNetFault, ErrFatalErr
	}
	log.Print(jsonline)

	var t RejectJoin
	dec := json.NewDecoder(strings.NewReader(jsonline))
	err = dec.Decode(&t)
	if err != nil {
		m.Print(err)    //MsgHdr解析出错
		return CmdFormatErr, err
	}

	var rol int16
	err = sqlTeamRol.QueryRow(t.TID,s.itemtype,s.itemid).Scan(&rol)
	if err != nil{
		m.Print(err)
		return CmdDbErr,err
	}

	//权限检查,管理员和创建者和队长才允许
	if rol != TeamOwer && rol != TeamLeader && rol != TeamManage{
		return CmdNotAllow,nil
	}

	msgbd := MsgBdNew{Mt:MsgSysBdDisTeam, St:s.itemtype, Si:s.itemid, Rt:ItemHuman, Ri:t.ID}
	msgbd.Bo = strconv.FormatInt((int64)(s.itemtype), 10)+
		"-"+strconv.FormatInt(s.itemid, 10)+":"+strconv.FormatInt(t.TID,10)+":"+t.De;
	LeaveMsg(msgbd)

	return CmdRejectJoinOK,nil
}

//查看未读留言的数量
func cmdLookMsgBd(rd *bufio.Reader, hdr MsgHdr, s SessionInfo)(string, error){
	var num int64

	log.Print("read data")
	jsonline,err := rd.ReadString('\n')
	if err != nil{
		m.Print(err)
		return CmdNetFault, ErrFatalErr
	}
	log.Print(jsonline)

	var t LookMsgBd
	dec := json.NewDecoder(strings.NewReader(jsonline))
	err = dec.Decode(&t)
	if err != nil {
		m.Print(err)    //MsgHdr解析出错
		return CmdFormatErr, err
	}

	err = sqlMsgBdNu.QueryRow(t.ID, MsgBdUnRd,s.itemtype,s.itemid).Scan(&num)
	if(err != nil){
		m.Print(err)
		return InternalErr, err
	}
	return CmdLookMsgBd+strconv.FormatInt(num,10)+"\n", nil
}

//解散Team
func cmdDisTeam(rd *bufio.Reader, hdr MsgHdr, s SessionInfo)(string, error){
	jsonline,err := rd.ReadString('\n')
	if err != nil{
		m.Print(err)
		return CmdNetFault, ErrFatalErr
	}

	var t DisTeam
	dec := json.NewDecoder(strings.NewReader(jsonline))
	err = dec.Decode(&t)
	if err != nil {
		m.Print(err)    //MsgHdr解析出错
		return CmdFormatErr, err
	}

	var rol int16
	err = sqlTeamRol.QueryRow(t.ID,s.itemtype,s.itemid).Scan(&rol)
	if err != nil{
		m.Print(err)
		return CmdDbErr,err
	}

	strid := strconv.FormatInt(t.ID,10)

	//权限检查,创建者才允许
	if rol != TeamOwer {
		return CmdNotAllow,nil
	}

	_, err = sqlTeamDel.Exec(t.ID,s.itemtype,s.itemid)
	if err != nil{
		m.Print(err)
		return CmdDbErr, err
	}

	_, err = sqlTeamDelMemAll.Exec(t.ID)
	if err != nil{
		m.Print(err)
		return CmdDbErr, err
	}

	_, err = sqlHumDelTeam.Exec(","+strid+",",s.itemid)

	msgbd := MsgBdNew{Mt:MsgSysBdRejectJoin, St:s.itemtype, Si:s.itemid, Rt:ItemTeam, Ri:t.ID}
	LeaveMsg(msgbd)

	msgrt := MsgRTNew{Mt:MsgSysRTDisTeam, St:s.itemtype, Si:s.itemid, Rt:ItemTeam, Ri:t.ID}
	SendMsg(msgrt,5)

	return CmdOK,nil
}

//添加新成员
func cmdAddMem(rd *bufio.Reader, hdr MsgHdr, s SessionInfo)(string, error){
	jsonline,err := rd.ReadString('\n')
	if err != nil{
		m.Print(err)
		return CmdNetFault, ErrFatalErr
	}

	var t MemAdd
	dec := json.NewDecoder(strings.NewReader(jsonline))
	err = dec.Decode(&t)
	if err != nil {
		m.Print(err)    //MsgHdr解析出错
		return CmdFormatErr, err
	}

	var rol int16
	err = sqlTeamRol.QueryRow(t.TI,s.itemtype,s.itemid).Scan(&rol)
	if err != nil{
		m.Print(err)
		return CmdDbErr,err
	}

	//权限检查,管理员和创建者和队长才允许
	if rol != TeamOwer && rol != TeamLeader && rol != TeamManage{
		return CmdNotAllow,nil
	}

	if(t.Tp == ItemHuman){
		var num int64
		strid := strconv.FormatInt(t.TI,10)
		err = sqlHumChkTeam.QueryRow(t.ID,"%,"+strid+",%").Scan(&num)
		if err == nil{
			return  CmdRepeatJoin, nil
		}else if err != sql.ErrNoRows{
			m.Print(err)
			return CmdDbErr,err
		}

		_,err = sqlHumAddTeam.Exec(t.TI,t.ID)
		if err != nil{
			m.Print(err)
			return CmdTeamMax, nil
		}
	}

	_,err = sqlTeamAddMem.Exec(t.TI,t.Tp,t.ID,t.Rl,t.Tg,t.De)
	if err != nil{
		m.Print(err)
		return CmdDbErr,nil
	}
	
	//给被加入者发送留言
	msgbd := MsgBdNew{Mt:MsgSysBdAdded, St:s.itemtype, Si:s.itemid, Rt:t.Tp, Ri:t.ID}
	msgbd.Bo = strconv.FormatInt(t.TI,10)
	LeaveMsg(msgbd)

	//给Team发送实时消息
	msgrt := MsgRTNew{Mt:MsgSysRTNewMem, St:s.itemtype, Rt:ItemTeam, Si:s.itemid,Ri:t.ID}
	msgrt.Bo = strconv.FormatInt(int64(s.itemtype),10) + "," + strconv.FormatInt(s.itemid,10)
	SendMsg(msgrt,5)

	return CmdOK,nil
}

//设置留言的状态
func cmdSetMsgBd(rd *bufio.Reader, hdr MsgHdr, s SessionInfo)(string, error){
	jsonline,err := rd.ReadString('\n')
	if err != nil{
		m.Print(err)
		return CmdNetFault, ErrFatalErr
	}

	var t SetMsgBd
	dec := json.NewDecoder(strings.NewReader(jsonline))
	err = dec.Decode(&t)
	if err != nil {
		m.Print(err)    //MsgHdr解析出错
		return CmdFormatErr, nil
	}

	_,err = sqlMsgBdUpRd.Exec(t.Rd,t.ID,s.itemtype,s.itemid)
	if err != nil{
		m.Print(err)
		return CmdDbErr,nil
	}

	return CmdOK,nil
}

//获取留言板的留言
func cmdGetMsgBd(rd *bufio.Reader, hdr MsgHdr, s SessionInfo)(string, error){
	jsonline,err := rd.ReadString('\n')
	if err != nil{
		m.Print(err)
		return CmdNetFault, ErrFatalErr
	}

	var t GetMsgBd
	dec := json.NewDecoder(strings.NewReader(jsonline))
	err = dec.Decode(&t)
	if err != nil {
		m.Print(err)    //MsgHdr解析出错
		return CmdFormatErr, nil
	}

	rows,err := sqlMsgBd.Query(s.itemtype,s.itemid,MsgBdUnRd,t.ID,t.To)
	if err != nil{
		m.Print(err)
		return CmdDbErr,nil
	}

	var buf []byte
	defer rows.Close()
	for rows.Next(){
		var msgbd MsgBd
		err = rows.Scan(&msgbd.ID,&msgbd.Mt,&msgbd.St,&msgbd.Si,&msgbd.Bo,&msgbd.Ti)
		if err != nil{
			m.Print(err)
			return CmdDbErr, nil
		}
		j, err := json.Marshal(msgbd)
		if err != nil{
			m.Print(err)
			return InternalErr, nil
		}
		buf = AppendSlice(buf,j)
	}
	err = rows.Err()
	if err != nil{
		m.Print(err)  //查询结果时可能出错,只做记录，不影响返回给用户的数据
	}

	if len(buf) == 0{
		return CmdNoMsgBd, nil
	}
	return CmdGetMsgBdOK+string(buf[:])+"\n", nil
}

//更新位置
func cmdUpPos(rd *bufio.Reader, hdr MsgHdr, s SessionInfo)(string, error){
	jsonline,err := rd.ReadString('\n')
	if err != nil{
		m.Print(err)
		return CmdNetFault, ErrFatalErr
	}

	if s.team == 0{
		return  CmdNotEnter, nil
	}

	var t Pos
	dec := json.NewDecoder(strings.NewReader(jsonline))
	err = dec.Decode(&t)
	if err != nil {
		m.Print(err)    //MsgHdr解析出错
		return CmdFormatErr, nil
	}

	_,err = sqlTeamUpMemPos.Exec(t.Lo,t.La,t.He,s.team,s.itemtype,s.itemid)
	if err != nil{
		m.Print(err)
		return CmdDbErr,nil
	}
	
	msg := MsgRTNew{Mt:MsgSysRTPosUp, St:s.itemtype, Rt:ItemTeam, Si:s.itemid, Ri:s.team}
	msg.Bo = strconv.FormatFloat(float64(t.Lo),'f',6,32)+","+
				strconv.FormatFloat(float64(t.La),'f',6,32)+","+
				strconv.FormatFloat(float64(t.He),'f',2,32)
	SendMsg(msg,5)   //消息有效时间是5分钟

	return CmdOK,nil
}

//退出Team
func cmdExitTeam(rd *bufio.Reader, hdr MsgHdr, s SessionInfo)(string, error){
	if s.itemtype != 1 {   //only human can search team
		return CmdNotAllow, nil
	}

	jsonline,err := rd.ReadString('\n')
	if err != nil{
		m.Print(err)
		return CmdNetFault, ErrFatalErr
	}

	var t ExitTeam
	dec := json.NewDecoder(strings.NewReader(jsonline))
	err = dec.Decode(&t)
	if err != nil {
		m.Print(err)    //MsgHdr解析出错
		return CmdFormatErr, nil
	}

	strid := strconv.FormatInt(t.ID,10)
	var otype int16
	var oid int64
	err = sqlTeamOwner.QueryRow(t.ID).Scan(&otype,&oid)
	if err == sql.ErrNoRows{  //Team不存在
		_, err = sqlHumDelTeam.Exec(","+strid+",",s.itemid)
		if err != nil {
			m.Print(err)    //MsgHdr解析出错
			return CmdDbErr, nil
		}
		return CmdOK, nil
	}else if err != nil{
		m.Print(err)
		return CmdDbErr,err
	}

	if oid == s.itemid && otype == s.itemtype{
		return CmdOwnerDeny, nil   //创始人不能退出Team
	}


	var num int64
	err = sqlHumChkTeam.QueryRow(s.itemid,"%,"+strid+",%").Scan(&num)
	if err == sql.ErrNoRows{
		return  CmdNotMem, nil
	}else if err != nil {
		m.Print(err)
		return CmdDbErr,err
	}

	_, err = sqlTeamDelMem.Exec(t.ID,s.itemtype,s.itemid)
	if err != nil {
		m.Print(err)    //MsgHdr解析出错
		return CmdDbErr, nil
	}

	_, err = sqlHumDelTeam.Exec(","+strid+",",s.itemid)

	msg := MsgRTNew{Mt:MsgSysBdExitTeam, St:s.itemtype, Rt:ItemTeam, Si:s.itemid,Ri:t.ID}
	SendMsg(msg,5)

	return CmdOK,nil
}

//离开当前的Team
func cmdLeaveTeam(rd *bufio.Reader, hdr MsgHdr, s SessionInfo)(string, error){

	if s.team == 0{
		return  CmdNotEnter, nil
	}

	_,err := sqlSesUpTeam.Exec(NoTeam, hdr.Si)
	if err != nil{
		m.Print(err)
		return CmdDbErr, nil
	}

	_,err = sqlTeamUpMemSt.Exec(MemOut, s.team, s.itemtype, s.itemid)
	if err != nil{
		m.Print(err)
		return CmdDbErr, nil
	}

	q := Identi{ItemType: ItemTeam, ItemID: s.team}
	a := Identi{ItemType: s.itemtype, ItemID: s.itemid}
	DelAudience(q, a)

	msg := MsgRTNew{Mt:MsgSysRTOutTeam, St:s.itemtype, Rt:ItemTeam, Si:s.itemid, Ri:s.team}
	SendMsg(msg,5)

	return CmdOK,nil
}

//进入Team
func cmdEnterTeam(rd *bufio.Reader, hdr MsgHdr, s SessionInfo)(string, error){
	
	jsonline,err := rd.ReadString('\n')
	if err != nil{
		m.Print(err)
		return CmdNetFault, ErrFatalErr
	}

	var t EnterTeam
	dec := json.NewDecoder(strings.NewReader(jsonline))
	err = dec.Decode(&t)
	if err != nil {
		m.Print(err)    //MsgHdr解析出错
		return CmdFormatErr, nil
	}

	//查看Team是否还存在
	strid := strconv.FormatInt(t.ID,10)
	var otype int16
	var oid int64
	err = sqlTeamOwner.QueryRow(t.ID).Scan(&otype,&oid)
	if err == sql.ErrNoRows{    //Team已经解散
		_, err = sqlHumDelTeam.Exec(","+strid+",",s.itemid)
		if err != nil {
			m.Print(err)
			return CmdDbErr, nil
		}
		return CmdTeamDis,nil
	}else if err != nil{
		m.Print(err)
		return CmdDbErr, nil
	}

	//查看用户的Teams中是否有该Team
	var num int64
	err = sqlHumChkTeam.QueryRow(s.itemid,"%,"+strid+",%").Scan(&num)
	if err == sql.ErrNoRows{   //用户的Teams不包含该Team
		return  CmdNotMem, nil
	}else if err != nil{
		m.Print(err)
		return CmdDbErr,err
	}

	//查看用户是否在TeamStat中
	err = sqlTeamFiMem.QueryRow(t.ID,s.itemtype,s.itemid).Scan(&num)
	if err == sql.ErrNoRows { //用户的Teams中有该Team,但是用户不在TeamStat中
		_, err = sqlHumDelTeam.Exec(","+strid+",",s.itemid)
		if err != nil {
			m.Print(err)
			return CmdDbErr, nil
		}
		return CmdNotMem, errors.New("Not in The Team")
	}else if err != nil{
		m.Print(err)
		return CmdDbErr, nil
	}


	_,err = sqlSesUpTeam.Exec(t.ID,hdr.Si)
	if err != nil{
		m.Print(err)
		return CmdDbErr, nil
	}

	_,err = sqlTeamUpMemSt.Exec(MemIn,t.ID,s.itemtype,s.itemid)
	if err != nil{
		m.Print(err)
		return CmdDbErr, nil
	}

	//获取Team的最新状态
	rows,err := sqlTeamStat.Query(t.ID)
	if err != nil{
		m.Print(err)
		return CmdDbErr, nil
	}

	var buf []byte
	defer rows.Close()
	for rows.Next(){
		var st ItemStat
		err = rows.Scan(&st.Mt,&st.Md,&st.Lt,&st.Ld,&st.Rl,&st.St,&st.Ar,&st.De,&st.Tg,&st.Lo,&st.La,&st.He,&st.Up)
		if err != nil{
			m.Print(err)
			return CmdDbErr, nil
		}
		j, err := json.Marshal(st)
		if err != nil{
			m.Print(err)
			return InternalErr, nil
		}
		buf = AppendSlice(buf,j)
	}
	err = rows.Err()
	if err != nil{
		m.Print(err)  //查询结果时可能出错,只做记录，不影响返回给用户的数据
	}

	//订阅Team消息
	q := Identi{ItemType: ItemTeam, ItemID: t.ID}
	a := Identi{ItemType: s.itemtype, ItemID: s.itemid}
	err = AddAudience(q,a,s.conn)
	if err != nil{
		return InternalErr, nil
	}

	msg := MsgRTNew{Mt:MsgSysRTInTeam, St:s.itemtype, Rt:ItemTeam, Si:s.itemid, Ri:t.ID}
	SendMsg(msg,5)
	return PushTeamStat+string(buf[:])+"\n", nil
}


//加入Team
func cmdJoinTeam(rd *bufio.Reader, hdr MsgHdr, s SessionInfo)(string, error){
	if s.itemtype != 1 {   //only human can search team
		return CmdNotAllow, nil
	}

	jsonline,err := rd.ReadString('\n')
	if err != nil{
		m.Print(err)
		return CmdNetFault, ErrFatalErr
	}

	var t JoinTeam
	dec := json.NewDecoder(strings.NewReader(jsonline))
	err = dec.Decode(&t)
	if err != nil {
		m.Print(err)    //MsgHdr解析出错
		return CmdFormatErr, nil
	}

	var num int64
	strid := strconv.FormatInt(t.ID,10)
	err = sqlHumChkTeam.QueryRow(s.itemid,"%,"+strid+",%").Scan(&num)
	if err == nil{
		return  CmdRepeatJoin, nil
	}else if err != sql.ErrNoRows{
		m.Print(err)
		return CmdDbErr,err
	}

	var role int16
	var alw int8
	var code string
	err = sqlTeamAlw.QueryRow(t.ID).Scan(&role,&alw,&code)
	if err == sql.ErrNoRows {
		return CmdNoTeam,nil
	}else if err != nil{
		m.Print(err)
		return CmdDbErr,nil
	}

	switch alw {
	case JoinFree: //任意加入
		_,err = sqlHumAddTeam.Exec(t.ID,s.itemid)
		if err != nil{
			m.Print(err)
			return CmdTeamMax, nil
		}
		_,err = sqlTeamAddMem.Exec(t.ID,s.itemtype,s.itemid,role,"",t.De)
		if err != nil{
			m.Print(err)
			return CmdDbErr,nil
		}

		msg := MsgRTNew{Mt:MsgSysRTNewMem, St:s.itemtype, Si:s.itemid, Rt:ItemTeam, Ri:t.ID}
		msg.Bo = strconv.FormatInt(int64(s.itemtype),10) + "," +
			strconv.FormatInt(s.itemid,10)
		SendMsg(msg,5)

		return CmdJoinTeamOK+strconv.FormatInt(t.ID,10)+"\n",nil
	case JoinAllow:  //需要批准
		msg := MsgBdNew{Mt:MsgSysBdJoinTeam, St:s.itemtype, Rt:ItemTeam, Si:s.itemid, Ri:t.ID}
		msg.Bo = strconv.FormatInt(t.ID,10)+":"+t.De
		LeaveMsg(msg)

		return CmdJoinTeamWait+strconv.FormatInt(t.ID,10)+"\n",nil
	case JoinByCode: //通过加入码加入
		if t.Co != code{
			return CmdCodeErr,nil
		}
		_,err = sqlHumAddTeam.Exec(t.ID,s.itemid)
		if err != nil{
			m.Print(err)
			return CmdTeamMax, nil
		}

		_,err = sqlTeamAddMem.Exec(t.ID,s.itemtype,s.itemid,role,"",t.De)
		if err!= nil{
			m.Print(err)
			return CmdDbErr,nil
		}

		msg := MsgRTNew{Mt:MsgSysRTNewMem, St:s.itemtype, Rt:ItemTeam, Si:s.itemid, Ri:t.ID}
		msg.Bo = strconv.FormatInt(int64(s.itemtype),10) + "," +
			strconv.FormatInt(s.itemid,10)
		SendMsg(msg,5)

		return CmdJoinTeamOK,nil
	case JoinClose:
		return CmdNotAllow,nil
	default:
		m.Print("The Team.Join 's value is Unknown:", code)
		return CmdNotAllow,nil
	}
}

//查找Team
func cmdFindTeam(rd *bufio.Reader, hdr MsgHdr, s SessionInfo)(string, error){
	var ret string
	if s.itemtype != 1 {   //only human can search team
		ret = CmdNotAllow
		return ret, nil
	}

	jsonline,err := rd.ReadString('\n')
	if err != nil{
		m.Print(err)
		return CmdNetFault,ErrFatalErr
	}

	var t FindTeam
	dec := json.NewDecoder(strings.NewReader(jsonline))

	err = dec.Decode(&t)
	if err != nil {
		m.Print(err)    //MsgHdr解析出错
		return CmdFormatErr, nil
	}

	var rows *sql.Rows
	rows = nil
	
	if len(t.Wd) == 0 && len(t.Tg) == 0{
		rows,err = sqlTeamFi.Query(OpenHid,t.Nu,t.Nu+int16(t.To))
		if err != nil && err != sql.ErrNoRows{
			m.Print(err)    //数据库错误
			return CmdDbErr, nil
		}
	}else if len(t.Wd) == 0{
		tag := "%"+t.Tg+"%"
		rows,err = sqlTeamFiTg.Query(OpenHid,tag,t.Nu,t.Nu+int16(t.To))
		if err != nil && err != sql.ErrNoRows{
			m.Print(err)    //数据库错误
			return CmdDbErr, nil
		}
	}else if len(t.Tg) == 0{
		wd := "%"+t.Wd+"%"
		rows,err = sqlTeamFiWd.Query(OpenHid,wd,t.Nu,t.Nu+int16(t.To))
		if err != nil && err != sql.ErrNoRows{
			m.Print(err)    //数据库错误
			return CmdDbErr, nil
		}
	}else{
		tag := "%"+t.Tg+"%"
		wd  := "%"+t.Wd+"%"
		rows,err = sqlTeamFiWdTg.Query(OpenHid,wd,tag,t.Nu,t.Nu+int16(t.To))
		if err != nil && err != sql.ErrNoRows{
			m.Print(err)    //数据库错误
			return CmdDbErr, nil
		}
	}
	
	var buf []byte
	defer rows.Close()
	for rows.Next(){
		var g TeamPubInfo
		err = rows.Scan(&g.ID,&g.Pu,&g.Aw,&g.Na,&g.Tg,&g.De)
		if err != nil{
			m.Print(err)
			return CmdDbErr, nil
		}
		j, err := json.Marshal(g)
		if err != nil{
			m.Print(err)
			return InternalErr, nil
		}
		buf = AppendSlice(buf,j)
	}

	err = rows.Err()
	if err != nil{
		m.Print(err)  //查询结果时可能出错,只做记录，不影响返回给用户的数据
	}

	if len(buf) == 0{
		return CmdNotFind,nil
	}
	return CmdTeamGetOK+string(buf)+"\n",nil
}

//获取所属的Team的列表,cmd="TeamList",没有后续数据
func cmdTeamList(rd *bufio.Reader, hdr MsgHdr, s SessionInfo)(string, error){
	var ret string
	err := errors.New("")
	err = nil

	switch s.itemtype{
	case ItemHuman:
		err = sqlHumTeam.QueryRow(s.itemid).Scan(&ret)
	case ItemPhyTrackTypeA:
		//TODO:
		err = errUnFi
	case ItemVirtualTrackerA1:
		//TODO:
		err = errUnFi
	case ItemVirtualTrackerA2:
		//TODO
		err = errUnFi
	default:  //Item表的内容没有与这里的代码逻辑同步
		err = errNotSup
	}

	if err != nil{          //内部错误
		m.Print(err);
		return CmdDbErr, nil
	}
	return CmdTeamListOK+ret+"\n", nil
}

//创建新的Team, cmd="NewTeam", 后面跟随用json表示的设置内容
func cmdNewTeam(rd *bufio.Reader, hdr MsgHdr, s SessionInfo)(string, error){
	if s.itemtype != 1 {   //only human can create team
		return CmdNotAllow, nil
	}

	jsonline,err := rd.ReadString('\n')
	if err != nil{
		m.Print(err)
		return CmdNetFault, ErrFatalErr
	}

	var t NewTeam
	dec := json.NewDecoder(strings.NewReader(jsonline))

	err = dec.Decode(&t)
	if err != nil {
		m.Print(err)    //MsgHdr解析出错
		return CmdFormatErr, nil
	}

	if t.Aw == JoinByCode && len(t.Co) == 0{
		return CmdNoJoinCode, errors.New("No JoinCode")
	}

	res,err := sqlTeamNew.Exec(t.Na,s.itemtype,s.itemid,t.De,t.Tg,t.Pu,t.Aw,t.Pr,t.Co,TeamMember)
	if err != nil{
		m.Print(err)
		return CmdDbErr, nil
	}

	tid, err := res.LastInsertId()
	if err != nil{
		m.Print(err)
		return CmdDbErr, nil
	}

	switch s.itemtype{
	case ItemHuman:
		_,err = sqlHumAddTeam.Exec(tid,s.itemid)
		if err != nil{
			m.Print(err)
			return CmdTeamMax, nil
		}
	default:
		return CmdNotAllow, errors.New("CreatTeam is Not Allowed")
	}

	_,err = sqlTeamAddMem.Exec(tid,s.itemtype,s.itemid,TeamOwer,"","")
	if err!= nil{
		m.Print(err)
		strid := strconv.FormatInt(tid,10)
		_, err = sqlHumDelTeam.Exec(","+strid+",",s.itemid)
		if err != nil {
			m.Print(err)    //MsgHdr解析出错
			return CmdDbErr, nil
		}
		return CmdDbErr,nil
	}
	
	teaminfo := TeamPubInfo{ID:tid, Pu:t.Pu, Aw:t.Aw, Na:t.Na, Tg:t.Tg, De:t.De}
	j,err := json.Marshal(teaminfo)
	if err != nil{
		m.Print(err)
		return InternalErr,nil
	}

	fname := strconv.FormatInt(int64(ItemTeam),10)+"-"+strconv.FormatInt(tid,10);
	err = WriteStaticFile(fileTeamPubInfo, fname, j);
	if err != nil{
		m.Print(err)
		return InternalErr,nil
	}

	return CmdNewTeamOK+strconv.FormatInt(tid, 10)+"\n", nil
}

//Logout
func cmdLogout(rd *bufio.Reader, hdr MsgHdr, s SessionInfo)(string, error){
	_,err := sqlSesDel.Exec(hdr.Si)
	if err != nil{
		m.Print(err)
		return CmdDbErr, err
	}
	return CmdOK, ErrFatalErr  //用ErrFatalErr来关闭连接
}
