//Copyright 2014. All rights reserved.
//Author: lja  
//Createtime: 2014/07/28 05:46:07  Lastchange: 2014/08/27 06:44:12
//changlog:  1. create by lja

package zainar

import (
	"log"
	"os"
	"database/sql"
	"runtime"
	"errors"
	"regexp"
)

var (
	//static file
	Uri_index string

	//gtpl file
	Uri_human_login string
	Uri_human_regist string
	
	www string   //web path
	Db *sql.DB   //Database

	l  *log.Logger    //login server的log
	m  *log.Logger    //msgserver的log

	//错误类型
	errUnFi        = errors.New("Unfinish")
	errNotSup      = errors.New("Not Support!")
	errNotFd       = errors.New("Not Found!")

	//正则匹配
	RegxMail   *regexp.Regexp    //邮箱

	//sql预处理语句
	sqlHumWait			*sql.Stmt
	sqlHumMailUsed		*sql.Stmt
	sqlHumNickUsed		*sql.Stmt
	sqlHumNew			*sql.Stmt
	sqlHumChkPwd		*sql.Stmt
	sqlHumLock			*sql.Stmt
	sqlHumStatNew		*sql.Stmt
	sqlHumDelTeam		*sql.Stmt
	sqlHumChkTeam		*sql.Stmt
	sqlHumAddTeam		*sql.Stmt
	sqlHumTeam			*sql.Stmt
	sqlTeamAlw			*sql.Stmt
	sqlTeamFi			*sql.Stmt
	sqlTeamFiWd			*sql.Stmt
	sqlTeamFiTg			*sql.Stmt
	sqlTeamFiWdTg		*sql.Stmt
	sqlTeamNew			*sql.Stmt
	sqlTeamOwner		*sql.Stmt
	sqlTeamSensi		*sql.Stmt
	sqlTeamAddMem		*sql.Stmt
	sqlTeamMemID		*sql.Stmt
	sqlTeamDel			*sql.Stmt
	sqlTeamUpdate		*sql.Stmt
	sqlTeamDelMem		*sql.Stmt
	sqlTeamDelMemAll	*sql.Stmt
	sqlTeamUpMemPos		*sql.Stmt
	sqlTeamUpMemSt		*sql.Stmt
	sqlTeamUpMemAr		*sql.Stmt
	sqlTeamFiMem		*sql.Stmt
	sqlTeamRol			*sql.Stmt
	sqlTeamStat			*sql.Stmt
	sqlMsgNew			*sql.Stmt
	sqlMsgLast			*sql.Stmt
	sqlMsgLastRec		*sql.Stmt
	sqlMsg				*sql.Stmt
	sqlSesNew			*sql.Stmt
	sqlSesTry			*sql.Stmt
	sqlSesUpTeam		*sql.Stmt
	sqlSesUpErr			*sql.Stmt
	sqlSesUpOk			*sql.Stmt
	sqlSesNewSW			*sql.Stmt
	sqlSesDel			*sql.Stmt
	sqlSesInfo			*sql.Stmt
	sqlMsgBdUpRd		*sql.Stmt
	sqlMsgBd			*sql.Stmt
	sqlMsgBdNew			*sql.Stmt
	sqlMsgBdLast		*sql.Stmt
	sqlMsgBdNu			*sql.Stmt
	sqlMsgBdNewNu		*sql.Stmt
	sqlHisLocNew		*sql.Stmt
	sqlHisLocFi			*sql.Stmt
	sqlHisLocDel		*sql.Stmt
	sqlHisLocDelSp		*sql.Stmt
)

const(
	//版本信息
	Version string = "v1.0"

	//预设大小
	MsgMaxSize int16 = 128
	LogTryNum int8 = 5                      //登陆尝试次数
	NoTeam  int64 = 0                       //不存在的Team的ID


	//Item类型
	ItemHuman              int16 = 1     //Human
	ItemPhyTrackTypeA      int16 = 2     //A型物理跟踪仪器
	ItemVirtualTrackerA1   int16 = 3     //A型虚拟跟踪仪，Android平台
	ItemVirtualTrackerA2   int16 = 4     //A型虚拟跟踪仪，IOS平台
	ItemTeam               int16 = 5     //Team
	ItemPublic             int16 = 6     //Public

	//Team的加入许可
	JoinFree   int8  =  0    //自由加入，不需要审批
	JoinAllow  int8  =  1    //加入需要审批
	JoinByCode int8  =  2    //通过加入码加入
	JoinClose  int8  =  3    //关闭加入

	//Team的开放程度
	OpenAll  int8   =  0    //信息完全公开
	OpenMem  int8   =  1    //信息对Team公开
	OpenHid  int8   =  2    //信息对Team公开，同时禁止被搜索到

	//Team Member的状态
	MemOut    int8  = 0    //没有进入Team
	MemIn     int8  = 1    //进入Team

	//Team Member的属性
	MemAttDefaut   int8  = 0
	MemAttrBlock   int8  = 1   //用户被封锁
	MemAttrGag     int8  = 2   //用户被禁言

	//角色
	TeamOwer       int16 = 1   //
	TeamLeader     int16 = 2   //
	TeamManage     int16 = 3   //
	TeamMember     int16 = 4   //
	SysAdmin       int16 = 5   //
	SilverUser     int16 = 6   //
	GoldenUser     int16 = 7   //
	DiamondUser    int16 = 8   //
	NormalUser     int16 = 9   //

	//留言状态
	MsgBdUnRd      int8 = 0 //未读
	MsgBdRd        int8 = 1 //已读
)

//身份标识
type Identi struct{
	ItemType int16
	ItemID   int64
}

//初始化, 在被import的时候自动运行
func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	fp, err := os.OpenFile("/opt/log/zainar_login.log."+hostname+Version,
			os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	l = log.New(fp, "["+os.Args[0]+"] ", log.Ldate | log.Lmicroseconds | log.Llongfile)

	fp, err = os.OpenFile("/opt/log/zainar_msg.log."+hostname,
			os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	m = log.New(fp, "["+os.Args[0]+"] ", log.Ldate | log.Lmicroseconds | log.Llongfile)

	log.SetFlags(log.Ldate|log.Lmicroseconds|log.Llongfile)
	//log.SetFlags(log.Ldate|log.Lmicroseconds)
	//log.SetOutput(fp)
	//log.SetPrefix("["+os.Args[0]+"]")


	//正则初始化
	RegxMail,err = regexp.Compile(".*@.*\\..*")
	if err != nil{
		log.Fatal(err)
	}

}

//对用户的密码加盐
func AddSalt(pwd, name string) string {
	return pwd+"+-xdaDFaAel^a*&)9"+name
}

//设置web资源路径
func SetWWWPath(path string) error{

	//TODO 检查文件是否存在
	www = path+"/"

	Uri_index     =  path + "/index.html"
	Uri_human_login  = path + "/human_login.gtpl"
	Uri_human_regist = path + "/human_regist.gtpl"

	return nil
}

//输入数据取决于所用的驱动的类型, go-sql-driver/mysql的输入如下:
//mysql shangwei:123456@tcp(127.0.0.1:3306)/shangwei?charset=utf8&collation=utf8_general_ci
func SetDatabase(driverName, dataSourceName string, openmax,idlemax int) error {
	tmpdb,err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}
	Db = tmpdb

	Db.SetMaxIdleConns(idlemax)
	Db.SetMaxOpenConns(openmax)

	//Human   
	//登陆时,查看账户的解锁时间
	sqlHumWait, err = Db.Prepare("select TIMEDIFF(Allow,NOW()) from Human where Mail=?")
	if err != nil{
		log.Fatal(err)
	}
	//注册和修改邮箱时, 检查邮箱是否被占用
	sqlHumMailUsed, err = Db.Prepare("select ID from Human where Mail=?")
	if err != nil{
		log.Fatal(err)
	}
	//注册和修改昵称时，检查昵称是否被占用
	sqlHumNickUsed, err = Db.Prepare("select ID from Human where Nick=?")
	if err != nil{
		log.Fatal(err)
	}
	//注册时，增加记录
	sqlHumNew, err = Db.Prepare("insert into Human set Mail=?,Pwd=SHA(?),Nick=?,MailIdenti=UUID(), IdentiFail=ADDDATE(NOW(),2),Allow=NOW()")
	if err != nil{
		log.Fatal(err)
	}
	//登陆时，检查密码是否正确
	sqlHumChkPwd, err = Db.Prepare("select Nick,RealMail from Human where ID=? and Pwd=SHA(?)")   //检查用户密码
	if err != nil{
		log.Fatal(err)
	}
	//错误登陆次数太多时，锁定用户
	sqlHumLock, err = Db.Prepare("update Human set Allow=ADDDATE(NOW(),INTERVAL 30 MINUTE) where ID=?")
	if err != nil{
		log.Fatal(err)
	}

	//HumanStat
	//注册时，增加记录
	sqlHumStatNew, err = Db.Prepare("insert into HumanStat set ID=?,Teams=',',Friends=','")
	if err != nil{
		log.Fatal(err)
	}
	//用户删除team时，更新Teams
	sqlHumDelTeam, err = Db.Prepare("update HumanStat set Teams=REPLACE(Teams,?,',') where ID=?")
	if err != nil{
		log.Fatal(err)
	}
	//用户加入team时，检查是否已经加入了
	sqlHumChkTeam, err = Db.Prepare("select ID from HumanStat where ID=? and Teams like ?")
	if err != nil{
		log.Fatal(err)
	}
	//用户加入team时，更新Teams
	sqlHumAddTeam, err = Db.Prepare("update HumanStat set Teams=concat(Teams,?,',') where ID=? and CHAR_LENGTH(Teams)<65500")
	if err != nil{
		log.Fatal(err)
	}
	//查看用户加入的team
	sqlHumTeam, err = Db.Prepare("select Teams from HumanStat where ID=?")
	if err != nil{
		log.Fatal(err)
	}

	//Team
	//用户加入team时，查看team的加入许可
	sqlTeamAlw, err = Db.Prepare("select Rol,Alw,Cod from Team where ID=?")
	if err != nil{
		log.Fatal(err)
	}
	//查找team
	sqlTeamFi, err = Db.Prepare("select ID,Pub,Alw,Na,Tg,De from Team where Pub!=? limit ?,?")
	if err != nil{
		log.Fatal(err)
	}
	//根据team的名称查找
	sqlTeamFiWd,err = Db.Prepare("select ID,Pub,Alw,Na,Tg,De from Team where Pub!=? and Na like ? limit ?,?")
	if err != nil{
		log.Fatal(err)
	}
	//根据team的标签查找
	sqlTeamFiTg,err = Db.Prepare("select ID,Pub,Alw,Na,Tg,De from Team where Pub!=? and Tg like ? limit ?,?")
	if err != nil{
		log.Fatal(err)
	}
	//根据team的名称和标签查找
	sqlTeamFiWdTg,err = Db.Prepare("select ID,Pub,Alw,Na,Tg,De from Team where Pub!=? and Na like? and Tg like ? limit ?,?")
	if err != nil{
		log.Fatal(err)
	}
	//创建team时, 插入记录
	sqlTeamNew, err = Db.Prepare("insert into Team set Na=?,CT=?,CID=?,De=?,Tg=?,Pub=?,Alw=?,Pro=?,Cod=?,Rol=?")
	if err != nil{
		log.Fatal(err)
	}

	//TeamStat
	//增加team成员
	sqlTeamAddMem, err = Db.Prepare("insert into TeamStat set TID=?,MT=?,MID=?,Rol=?,Tg=?,De=?,Birth=NOW(),Up=NOW()")
	if err != nil{
		log.Fatal(err)
	}

	//查看Team创始人
	sqlTeamOwner, err = Db.Prepare("select CT,CID from Team where ID=?")
	if err != nil{
		log.Fatal(err)
	}

	//获取Team的敏感信息,JoinCode等
	sqlTeamSensi, err = Db.Prepare("select Pro, Cod from Team where ID=?")
	if err != nil{
		log.Fatal(err)
	}

	//查看指定类型Team成员的ID
	sqlTeamMemID, err = Db.Prepare("select MID from TeamStat where TID=? and MT=?")
	if err != nil{
		log.Fatal(err)
	}

	//更新Team信息
	sqlTeamUpdate, err = Db.Prepare("update Team set Na=?,De=?,Tg=?,Pub=?,Alw=?,Pro=?,Cod=? where ID=?")
	if err != nil{
		log.Fatal(err)
	}

	//解散team
	sqlTeamDel, err = Db.Prepare("delete from Team where ID=? and CT=? and CID=?")
	if err != nil{
		log.Fatal(err)
	}
	//删除team成员
	sqlTeamDelMem, err = Db.Prepare("delete from TeamStat where TID=? and MT=? and MID=?")
	if err != nil{
		log.Fatal(err)
	}
	//删除Team所有成员
	sqlTeamDelMemAll, err = Db.Prepare("delete from TeamStat where TID=?")
	if err != nil{
		log.Fatal(err)
	}
	//更新team成员的位置
	sqlTeamUpMemPos, err = Db.Prepare("update TeamStat set Lo=?,La=?,He=? where TID=? and MT=? and MID=?")
	if err != nil{
		log.Fatal(err)
	}
	//更新team成员的状态
	sqlTeamUpMemSt, err = Db.Prepare("update TeamStat set Stat=? where TID=? and MT=? and MID=?")
	if err != nil{
		log.Fatal(err)
	}
	//更新team成员的属性
	sqlTeamUpMemAr, err = Db.Prepare("update TeamStat set Attr=? where TID=? and MT=? and MID=?")
	if err != nil{
		log.Fatal(err)
	}
	//检查team成员是否存在
	sqlTeamFiMem, err = Db.Prepare("select MID from TeamStat where TID=? and MT=? and MID=?")
	if err != nil{
		log.Fatal(err)
	}
	//获取在team中的角色
	sqlTeamRol, err = Db.Prepare("select Rol from TeamStat where TID=? and MT=? and MID=?")
	if err != nil{
		log.Fatal(err)
	}
	//获取team当前状态
	sqlTeamStat, err = Db.Prepare("select MT,MID,LT,LID,Rol,Stat,Attr,De,Tg,Lo,La,He,Up from TeamStat where TID=?")
	if err != nil{
		log.Fatal(err)
	}

	//Msg
	//生成新的MsgRT时，插入记录
	sqlMsgNew, err = Db.Prepare("insert into MsgRT set MT=?,ST=?,SID=?,RT=?,RID=?,Bo=?,Exp=DATE_ADD(NOW(),INTERVAL ? MINUTE)")
	if err != nil{
		log.Fatal(err)
	}
	//最后一条MsgRT的ID
	sqlMsgLast, err = Db.Prepare("select * from (select MAX(ID) as ID from MsgRT) as t where t.ID>0")
	if err != nil{
		log.Fatal(err)
	}
	//指定接收对象的最后一条MsgRT的ID
	sqlMsgLastRec, err = Db.Prepare("select * from (select MAX(ID) as ID from MsgRT where RT=? and RID=?) as t where t.ID>0")
	if err != nil{
		log.Fatal(err)
	}
	//读取MsgRT
	sqlMsg, err = Db.Prepare("select MT,ST,SID,RT,RID,Bo from MsgRT where ID>? and ID<=? and RT=? and RID=? and Exp>NOW()")
	if err != nil{
		log.Fatal(err)
	}

	//Session
	//新建session
	sqlSesNew, err = Db.Prepare("insert into Session set IID=?,IT=?,Ok=0,Err=0,Birth=NOW(),Up=NOW()")
	if err != nil{
		log.Fatal(err)
	}
	//查看session的登陆尝试
	sqlSesTry, err = Db.Prepare("select ID,Ok,Err from Session where IT=? and IID=?")
	if err != nil{
		log.Fatal(err)
	}
	//更新session当前的Team
	sqlSesUpTeam, err = Db.Prepare("update Session set TID=?,Up=NOW() where ID=?")
	if err != nil{
		log.Fatal(err)
	}
	//更新session错误登陆次数
	sqlSesUpErr, err = Db.Prepare("update Session set Err=?,Up=NOW() where ID=?")
	if err != nil{
		log.Fatal(err)
	}
	//更新session登陆成功
	sqlSesUpOk, err = Db.Prepare("update Session set Ok=1,Err=0,CW=?,CP=?,SW=?,SP=?,Up=NOW() where ID=?")
	if err != nil{
		log.Fatal(err)
	}
	//设置新的SW
	sqlSesNewSW, err = Db.Prepare("update Session set SW=? where ID=?")
	if err != nil{
		log.Fatal(err)
	}

	//删除session
	sqlSesDel, err = Db.Prepare("delete from Session where ID=?")
	if err != nil{
		log.Fatal(err)
	}
	//读取session信息
	sqlSesInfo, err =  Db.Prepare("select IID,IT,TID,CW,CP,SW,SP from Session where ID=?")
	if err != nil{
		log.Fatal(err)
	}

	//添加一条留言
	sqlMsgBdNew, err = Db.Prepare("insert into MsgBd set Rd=?,MT=?,ST=?,SID=?,RT=?,RID=?,Bo=?,Birth=NOW()")
	if err != nil{
		log.Fatal(err)
	}

	//读取留言
	sqlMsgBd, err = Db.Prepare("select ID,MT,ST,SID,Bo,Birth from MsgBd where RT=? and RID=? and Rd=? and ID>? limit ?")
	if err != nil{
		log.Fatal(err)
	}

	//更新留言状态
	sqlMsgBdUpRd, err = Db.Prepare("update MsgBd set Rd=? where ID=? and RT=? and RID=?")
	if err != nil{
		log.Fatal(err)
	}

	//读取所有未读留言数量
	sqlMsgBdNu, err = Db.Prepare("select count(*) from MsgBd where ID>? and Rd=? and RT=? and RID=?")
	if err != nil{
		log.Fatal(err)
	}

	//读取最新未读留言的数量
	sqlMsgBdNewNu, err = Db.Prepare("select count(*) from MsgBd where Rd=? and ID>? and ID<= ? and RT=? and RID=?")
	if err != nil{
		log.Fatal(err)
	}


	//最后一条留言的ID
	sqlMsgBdLast, err = Db.Prepare("select * from (select MAX(ID) as ID from MsgBd where Rd=?) as t where t.ID>0")
	if err != nil{
		log.Fatal(err)
	}

	//插入一条历史位置记录
	sqlHisLocNew, err = Db.Prepare("insert into HisLoc set IT=?,IID=?,Lo=?,La=?,He=?,Up=NOW()")
	if err != nil{
		log.Fatal(err)
	}

	//查找特定时间范围内，特定成员的位置记录
	sqlHisLocFi, err = Db.Prepare("select Lo,La,He,Up from HisLoc where IT=? and IID=? and Up>? and Up<?")
	if err != nil{
		log.Fatal(err)
	}
	
	//删除制定时间范围内的数据
	sqlHisLocDel, err = Db.Prepare("delete from HisLoc where Up>? and Up <?")
	if err != nil{
		log.Fatal(err)
	}

	//删除制定时间范围内特定用户的数据
	sqlHisLocDelSp, err = Db.Prepare("delete from HisLoc where IT=? and IID=? and Up>? and Up<?")
	if err != nil{
		log.Fatal(err)
	}

	return nil
}

//关闭数据库
func CloseDatabase() error{

	sqlHumWait.Close()
	sqlHumMailUsed.Close()
	sqlHumNickUsed.Close()
	sqlHumNew.Close()
	sqlHumChkPwd.Close()
	sqlHumLock.Close()
	sqlHumStatNew.Close()
	sqlHumDelTeam.Close()
	sqlHumChkTeam.Close()
	sqlHumAddTeam.Close()
	sqlHumTeam.Close()
	sqlTeamAlw.Close()
	sqlTeamFi.Close()
	sqlTeamFiWd.Close()
	sqlTeamFiTg.Close()
	sqlTeamFiWdTg.Close()
	sqlTeamNew.Close()
	sqlTeamOwner.Close()
	sqlTeamAddMem.Close()
	sqlTeamMemID.Close()
	sqlTeamDel.Close()
	sqlTeamDelMem.Close()
	sqlTeamDelMemAll.Close()
	sqlTeamUpMemPos.Close()
	sqlTeamUpMemSt.Close()
	sqlTeamUpMemAr.Close()
	sqlTeamFiMem.Close()
	sqlTeamRol.Close()
	sqlTeamStat.Close()
	sqlMsgNew.Close()
	sqlMsgLast.Close()
	sqlMsgLastRec.Close()
	sqlMsg.Close()
	sqlSesNew.Close()
	sqlSesTry.Close()
	sqlSesUpTeam.Close()
	sqlSesUpErr.Close()
	sqlSesUpOk.Close()
	sqlSesNewSW.Close()
	sqlSesDel.Close()
	sqlSesInfo.Close()
	sqlMsgBdUpRd.Close()
	sqlMsgBd.Close()
	sqlMsgBdNew.Close()
	sqlMsgBdLast.Close()
	sqlMsgBdNu.Close()
	sqlMsgBdNewNu.Close()
	sqlHisLocNew.Close()
	sqlHisLocFi.Close()
	sqlHisLocDel.Close()
	sqlHisLocDelSp.Close()

	err := Db.Close()
	if err != nil {
		return err
	}


	return nil
}
