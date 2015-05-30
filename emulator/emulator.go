//Copyright 2014. All rights reserved.
//Author: lja  
//Createtime: 2014/07/30 10:28:59  Lastchange: 2014/08/26 19:03:50
//changlog:  1. create by lja

package main

import (
	"net"
	"net/http"
	"net/url"
	"log"
	"fmt"
	"os"
	"flag"
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"time"
	"crypto/tls"
	"zainar"
)

var(
	//domain name
	logindn  *string
	msgdn    *string

	//Current itemstat
	action string

	itemname string
	account  string
	status   string  //"on", "off", "trylogout"
	start  time.Time  //指令发送出去的时间, 用来计算返回延迟
)

type CmdInput func(hdr *zainar.MsgHdr)([]byte, error)

var(
	Cmds ="TL:TeamList NT:NewTeam FT:FindTeam JT:JoinTeam DS:DissTeam"+
		"IT:InTeam UP:UpPos LT:LeaveTeam ET:ExitTeam LO:Logout GM:GetMsgBd AM:AddMem"

	CmdInputs = map[string]CmdInput{
		"TL":  inputTeamList,
		"NT":  inputNewTeam,
		"FT":  inputFindTeam,
		"JT":  inputJoinTeam,
		"IT":  inputEnterTeam,
		"UP":  inputUpPos,
		"LT":  inputLeaveTeam,
		"ET":  inputExitTeam,
		"LO":  inputLogout,
		"GM":  inputGetMsgBd,
		"LM":  inputLookMsgBd,
		"SM":  inputSetMsgBd,
		"AM":  inputAddMem,
		"DS":  inputDisTeam,
	}
)

func inputDisTeam(hdr *zainar.MsgHdr)([]byte, error){
	var mem zainar.DisTeam
	fmt.Printf("TeamID:")
	fmt.Scanf("%d", &mem.ID)

	json_t_stat, err := json.Marshal(mem)
	if err != nil {
		log.Fatal(err)
	}

	hdr.Da = 0
	return json_t_stat, err

}
func inputAddMem(hdr *zainar.MsgHdr)([]byte, error){
	var mem zainar.MemAdd

	fmt.Printf("MemType:")
	fmt.Scanf("%d", &mem.Tp)

	fmt.Printf("MemID:")
	fmt.Scanf("%d", &mem.ID)

	fmt.Printf("TeamID:")
	fmt.Scanf("%d", &mem.TI)

	fmt.Printf("De:")
	fmt.Scanf("%s", &mem.De)

	fmt.Printf("Tg:")
	fmt.Scanf("%s", &mem.Tg)

	fmt.Printf("Role:")
	fmt.Scanf("%d", &mem.Rl)

	json_t_stat, err := json.Marshal(mem)
	if err != nil {
		log.Fatal(err)
	}

	hdr.Da = 0
	return json_t_stat, err
}

func inputSetMsgBd(hdr *zainar.MsgHdr)([]byte, error){
	var mbd zainar.SetMsgBd

	fmt.Printf("ID:")
	fmt.Scanf("%d", &mbd.ID)

	fmt.Printf("Read(1)/UnRead(0):")
	fmt.Scanf("%d", &mbd.Rd)

	json_t_stat, err := json.Marshal(mbd)
	if err != nil {
		log.Fatal(err)
	}

	hdr.Da = 0
	return json_t_stat, err
}

func inputLookMsgBd(hdr *zainar.MsgHdr)([]byte, error){
	var mbd zainar.LookMsgBd
	fmt.Printf("Start ID:")
	fmt.Scanf("%d", &mbd.ID)

	json_t_stat, err := json.Marshal(mbd)
	if err != nil {
		log.Fatal(err)
	}

	hdr.Da = 0
	return json_t_stat, err
}

func inputGetMsgBd(hdr *zainar.MsgHdr)([]byte, error){
	var bd zainar.GetMsgBd

	fmt.Printf("Start ID:")
	fmt.Scanf("%d", &bd.ID)

	fmt.Printf("Tot:")
	fmt.Scanf("%d", &bd.To)

	json_t_stat, err := json.Marshal(bd)
	if err != nil {
		log.Fatal(err)
	}

	hdr.Da = 0
	return json_t_stat, err
}

func inputUpPos(hdr *zainar.MsgHdr)([]byte, error){
	var pos zainar.Pos

	fmt.Printf("Longtitude:")
	fmt.Scanf("%f", &pos.Lo)

	fmt.Printf("Latitude:")
	fmt.Scanf("%f", &pos.La)

	fmt.Printf("Height:")
	fmt.Scanf("%f", &pos.He)

	json_t_stat, err := json.Marshal(pos)
	if err != nil {
		log.Fatal(err)
	}

	hdr.Da = 0
	return json_t_stat, err
}

func inputExitTeam(hdr *zainar.MsgHdr)([]byte, error){
	var exit zainar.ExitTeam

	fmt.Printf("TeamID:")
	fmt.Scanf("%d",&exit.ID)

	json_t_stat, err := json.Marshal(exit)
	if err != nil {
		log.Fatal(err)
	}

	hdr.Da = 0
	return json_t_stat, err
}

func inputLeaveTeam(hdr *zainar.MsgHdr)([]byte, error){
	return nil,nil
}

func inputEnterTeam(hdr *zainar.MsgHdr)([]byte, error){
	var enterteam zainar.EnterTeam

	fmt.Printf("TeamID:")
	fmt.Scanf("%d",&enterteam.ID)

	json_t_stat, err := json.Marshal(enterteam)
	if err != nil {
		log.Fatal(err)
	}

	hdr.Da = 0
	return json_t_stat, err
}

func inputJoinTeam(hdr *zainar.MsgHdr)([]byte, error){
	var jointeam zainar.JoinTeam

	fmt.Printf("TeamID:")
	fmt.Scanf("%d",&jointeam.ID)

	fmt.Printf("JoinCode:")
	fmt.Scanf("%s",&jointeam.Co)

	fmt.Printf("SelfDescripe:")
	fmt.Scanf("%s", &jointeam.De)

	json_t_stat, err := json.Marshal(jointeam)
	if err != nil {
		log.Fatal(err)
	}

	hdr.Da = 0
	return json_t_stat, err
}

func inputFindTeam(hdr *zainar.MsgHdr)([]byte, error){
	var findteam zainar.FindTeam

	fmt.Printf("TeamName:")
	fmt.Scanf("%s",&findteam.Wd)

	fmt.Printf("TeamTag:")
	fmt.Scanf("%s",&findteam.Tg)

	fmt.Printf("Start Num:")
	fmt.Scanf("%d",&findteam.Nu)
	
	fmt.Printf("Tot Num:")
	fmt.Scanf("%d",&findteam.To)

	json_t_stat, err := json.Marshal(findteam)
	if err != nil {
		log.Fatal(err)
	}

	hdr.Da = 0
	return json_t_stat, err
}

func inputTeamList(hdr *zainar.MsgHdr)([]byte, error){
	hdr.Da = 0
	return nil,nil
}

func inputNewTeam(hdr *zainar.MsgHdr)([]byte, error){
	var t zainar.NewTeam

	fmt.Printf("Name:")
	fmt.Scanf("%s", &t.Na)

	fmt.Printf("Des:")
	fmt.Scanf("%s", &t.De)

	fmt.Printf("Tag:")
	fmt.Scanf("%s", &t.Tg)

	fmt.Printf("Pub(%d:信息完全公开 %d:信息仅对Team成员公开 %d:隐藏Team):",
		zainar.OpenAll, zainar.OpenMem, zainar.OpenHid)
	fmt.Scanf("%d", &t.Pu)

	fmt.Printf("Join(%d 随意加入 %d 加入需要审批 %d 禁止加入 %d 凭加入码加入):",
		zainar.JoinFree,zainar.JoinAllow,zainar.JoinClose,zainar.JoinByCode)
	fmt.Scanf("%d", &t.Aw)

	switch t.Aw{
	case zainar.JoinByCode: //凭加入码加入
		fmt.Printf("JoinPromt(加入码提示):")
		fmt.Scanf("%s", &t.Pr)
		fmt.Printf("JoinCode(加入码):")
		fmt.Scanf("%s", &t.Co)
	}

	json_t_stat, err := json.Marshal(t)
	if err != nil {
		log.Fatal(err)
	}

	hdr.Da = 0
	return json_t_stat, err
}

func inputLogout(hdr *zainar.MsgHdr)([]byte, error){
	hdr.Da = 0
	status = "trylogout"
	return nil,nil
}

func send_msg(conn *net.TCPConn, buf []byte){
	fmt.Printf("Send: %s", buf)
	_, err := conn.Write(buf)
	start = time.Now()
	if err != nil {
		log.Fatal(err)
	}
	//Wait for Response
	time.Sleep(1 *time.Second)
}

func automsg(humaninfo *zainar.HumanInfo)([]byte, error){
	var hdr zainar.MsgHdr

	//Sid
	hdr.Si = humaninfo.SI

	//Msg Num
	hdr.Nu = humaninfo.SW

	//Cmd
	fmt.Printf("cmd(%s):",Cmds)
	fmt.Scanf("%s", &hdr.Cm)

	input, ok := CmdInputs[hdr.Cm]
	if ok == false{
		return nil, errors.New("Cmd Wrong: Don't have this cmd")
	}

	data,err := input(&hdr)
	if err != nil{
		log.Fatal(err)
	}

	jsonhdr, err := json.Marshal(hdr)
	if err != nil {
		log.Fatal(err)
	}

	//Combine the Msg
	send := append(jsonhdr, '\n')
	send = zainar.AppendSlice(send, data)
	send = append(send, '\n')

	return send, nil
}

//Login成功后使用的消息发送程序，自动填充消息中非交互内容
func login_msg(humaninfo zainar.HumanInfo){
	if status == "off" {
		fmt.Printf("You must Login!\n")
	}
	addr, err := net.ResolveTCPAddr("tcp", *msgdn)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	go receiver(conn)
	for {
		if status == "off" {
			fmt.Printf("You have Logout.\n")
			return
		}
		send,err := automsg(&humaninfo) //构造msg
		if err != nil{
			log.Print(err)
			continue
		}

		send_msg(conn, send)
		humaninfo.SW += int32(humaninfo.SP)
	}
}

func human_login()(zainar.HumanInfo, error){
	var humaninfo zainar.HumanInfo
	var mail string
	var pwd string

	fmt.Printf("Mail:")
	fmt.Scanf("%s", &mail)
	fmt.Printf("Password:")
	fmt.Scanf("%s", &pwd)
	
	//emualtor status
	account = mail
	//end

	//跳过证书检查
	tr := http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify:true}}
	client := http.Client{Transport: &tr}

	resp, err := client.PostForm("http://"+*logindn+"/login",
		url.Values{"mail":{mail},"pwd":{pwd},"type":{itemname}})
	if err != nil{
		log.Fatal(err)
	}
	defer resp.Body.Close()

	buf := make([]byte, 1024)
	n, err := resp.Body.Read(buf)
	if err != nil && err != io.EOF{
		log.Fatal(err)
	}
	fmt.Printf(">%s\n",buf)

	if buf[0]=='E' && buf[1]=='R' && buf[2] ==' '{
		return humaninfo, errors.New(string(buf))
	}
	
	err = json.Unmarshal(buf[0:n], &humaninfo)
	if err != nil{
		log.Fatal(err)
		return humaninfo, err
	}
	return humaninfo, nil
}

func login(){
	fmt.Printf("ItemName(Human, VirtualTrackerA1):")
	fmt.Scanf("%s", &itemname)
	
	switch itemname {
	case "Human":
		for {
			status = "off"
			humaninfo, err := human_login()
			if err != nil{
				continue
			}
			status = "on"
			login_msg(humaninfo)
			return   //You have logout
		}
	case "VirtualTrackerA1":
		fmt.Printf("Unfinished")
		return
	default:
		fmt.Printf("Don't have thie item\n")
		return
	}
}


//手动填充消息的所有内容
func manualmsg()([]byte, error){
	var hdr zainar.MsgHdr

	//Sid
	fmt.Printf("Sid:")
	fmt.Scanf("%d", &hdr.Si)

	//Msg Num
	fmt.Printf("Num:")
	fmt.Scanf("%d", &hdr.Nu)

	//Cmd
	fmt.Printf("cmd(")
	for k,_ := range CmdInputs{
		fmt.Printf("%s, ",k)
	}
	fmt.Printf("):")
	fmt.Scanf("%s", &hdr.Cm)

	//Data
	var data string
	fmt.Printf("dat:")
	fmt.Scanf("%s", &data)
	hdr.Da = len(data)

	jsonhdr, err := json.Marshal(hdr)
	if err != nil {
		log.Fatal(err)
	}

	//Combine the Msg
	send := append(jsonhdr, '\n')
	n := len(data)
	for i:=0; i < n; i++{
		send = append(send, data[i])
	}
	send = append(send, '\n')

	return send, nil
}

//接收goroutine
func receiver(conn *net.TCPConn){
	rd := bufio.NewReader(conn)
	for {
		//Wait response, read cmd status
		res, err := rd.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		if(res == "\n"){
			continue
		}
		fmt.Printf(">%s", res)

		//Logout执行成功, 这个纯粹是为emualtor服务的, 方便emuator的使用
		if status == "trylogout" && res[0:3] == "OK "{
			status = "off"
			return
		}

		//Wait response, read cmd result data
		res, err = rd.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(">%s", res)
		ns := time.Now().Sub(start).Nanoseconds()
		ms := ns/1000000
		fmt.Printf(">use: %dns (%dms)\n", ns, ms)
	}
}

//干预消息通信，手动填充非用户的填充的内容，用于构造异常数据
func inter_msg(){

	addr, err := net.ResolveTCPAddr("tcp", *msgdn)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	go receiver(conn)
	for {
		send,err := manualmsg() //构造msg
		if err != nil{
			log.Print(err)
			continue
		}
		send_msg(conn, send)
	}
}

func regist(){
	var mail string
	var pwd1 string
	var pwd2 string
	var nick string

	tr := http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify:true}}
	client := http.Client{Transport: &tr}

	for {
		fmt.Printf("ItemName(Human, VirtualTrackerA1):")
		fmt.Scanf("%s", &itemname)

		switch itemname{
		case "Human":
			fmt.Printf("Email:")
			fmt.Scanf("%s", &mail)
			fmt.Printf("Password:")
			fmt.Scanf("%s", &pwd1)
			fmt.Printf("Password(again):")
			fmt.Scanf("%s", &pwd2)

			if pwd1 != pwd2{
				fmt.Printf("passwd isn't consisent\n")
				continue
			}

			fmt.Printf("NickName:")
			fmt.Scanf("%s",&nick)

			resp, err := client.PostForm("http://"+*logindn+"/register",
				url.Values{"mail":{mail},"pwd":{pwd1},"nick":{nick},"type":{itemname}})
			if err != nil{
				log.Fatal(err)
			}
			defer resp.Body.Close()

			buf := make([]byte, 1024)
			_, err = resp.Body.Read(buf)
			if err != nil && err != io.EOF{
				log.Fatal(err)
			}
			fmt.Printf(">%s\n",buf)

			if buf[0]=='E' && buf[1]=='R' && buf[2] ==' '{
				fmt.Printf("Register Error!\n")
				continue
			}
			
			fmt.Printf("Reigster Sucess! You can Login in Now\n")
			return
		case "VirtualTrackerA1":
			fmt.Printf("Unfinished!\n")
			continue
		default:
			fmt.Printf("Don't have this item\n")
			continue
		}
	}
}

func main() {
	log.Print("Start ", os.Args)

	log.SetFlags(log.Ldate|log.Lmicroseconds|log.Llongfile)
	log.SetPrefix("["+os.Args[0]+"]")

	logindn = flag.String("logindn","192.168.88.130:443","Login Server Domain Name")
	msgdn = flag.String("msgdn","192.168.88.130:8000","Msg Server Domain Name")
	flag.Parse()

	for {
		fmt.Printf("Action(Login, Regist, Msg, Quit):")
		fmt.Scanf("%s", &action)
		switch action{
		case "Login":
			login()
		case "Regist":
			regist()
		case "Msg":
			inter_msg()
			continue
		case "Quit":
			fmt.Printf("Bye.\n")
			return
		default:
			fmt.Printf("Don't have this Action.\n")
			continue
		}
	}
}
