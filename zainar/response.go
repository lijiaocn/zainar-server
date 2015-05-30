//Copyright 2014. All rights reserved.
//Author: lja  
//Createtime: 2014/08/05 16:18:33  Lastchange: 2014/08/27 03:01:07
//changlog:  1. create by lja

package zainar

import (
)

//回应的数据有分为消息头和消息体两部分,用'\n'分隔
//消息头表示消息的执行结果, 成功用"OK "开头，失败用"ER "开头
//消息体json表示
//如果带有Json消息体，"OK "后面跟随消息体类型标志，例如"OK TL\n"
const(
	//消息头
	CmdReconnect  string = "ER Connect\n"     //通知客户端重新连接
	CmdReLogin    string = "ER Login\n"       //通知客户端重新登陆
	CmdInvalid    string = "ER Invalid\n"     //无效的消息
	CmdNotAllow   string = "ER NotAllow\n"    //不允许操作
	CmdNetFault   string = "ER NetFault\n"   //网络故障
	CmdFormatErr  string = "ER FomatErr\n"    //格式错误
	CmdDbErr      string = "ER Internal\n"    //格式错误
	CmdTeamMax    string = "ER TeamMax\n"     //参与的Team数量太多
	CmdNoJoinCode string = "ER NoJoinCode\n"  //新建Team选择了加入码方式，但是没有设置加入码
	CmdRepeatJoin string = "ER RepeatJoin\n"  //重复加入同一个Team
	CmdNoTeam     string = "ER NoTeam\n"      //Team不存在
	CmdNotEnter   string = "ER NotEnter\n"    //未进入Team
	CmdNotFind    string = "ER NotFind\n"     //没找到
	CmdCodeErr    string = "ER CodeErr\n"     //加入码错误
	CmdNotMem     string = "ER NotMem\n"      //不是Team成员
	CmdNoMsgBd    string = "ER NoMsgBd\n"     //没有留言
	CmdOwnerDeny  string = "ER OwnDeny\n"     //创始人不能执行这个操作
	CmdTeamDis    string = "ER TeamDis\n"     //团队已经解散

	//Log登陆结果
	LogNotExist string = "ER NotExist\n"     //登陆账号不存在
	LogTryMax   string = "ER MaxTry\n"       //登陆错误次数达到上线
	LogWrong    string = "ER Wrong\n"        //账户或密码错误
	LogLock     string = "ER Lock\n"         //账户被锁定, 后面跟一行等待时间

	//注册结果
	RegMailUsed   string = "ER MailUsed\n"    //注册邮箱已经被占用
	RegNickUsed   string = "ER NickUsed\n"    //昵称已经被占用

	//共用的结果
	ItemNotExist string = "ER NoItem\n"      //物品类型不存在
	InternalErr  string = "ER Internal\n"    //内部错误

	//输入检查不通过
	IglMail      string = "ER IllegelMail\n" //邮件格式错误
	IglForm      string = "ER IllegeForm\n"  //非法表单

	CmdOK         string = "OK \n"            //客户消息正确,后面跟随消息结果
	CmdJoinTeamWait   string = "OK WT\n"      //消息已被接收等待处理
	CmdJoinTeamOK string = "OK JT\n"          //Team加入成功
	CmdTeamGetOK  string = "OK TG\n"          //消息体是TeamGet
	CmdTeamListOK string = "OK TL\n"          //消息体是用","分隔的TeamID
	CmdNewTeamOK  string = "OK NT\n"          //消息体是新建立的TeamID
	CmdGetMsgBdOK string = "OK MB\n"          //消息体是查询的留言板消息
	CmdLookMsgBd  string = "OK LM\n"          //消息体是未读的留言的数量
	CmdRejectJoinOK string = "OK RJ\n"        //Reject Ok
	CmdTeamSensiOK string = "OK TS\n"         //消息体是Team敏感信息
	PushMsgRT     string = "OK MS\n"          //消息体是推送的实时消息
	PushTeamStat  string = "OK IS\n"          //消息体是推送的ItemStat
	PushMsgBdNum  string = "OK MBN\n"         //消息体是留言数量

	//消息类型
	MsgRTTinyText    int8  = 1   //实时短消息
	MsgRTLongText    int8  = 2   //实时长消息
	MsgRTFile        int8  = 3   //实时文件
	MsgRTVideo       int8  = 4   //实时视频流
	MsgRTVoice       int8  = 5   //实时语音流

	MsgBdTinyText    int8  = 11   //留言短消息
	MsgBdLongText    int8  = 12   //留言长消息
	MsgBdFile        int8  = 13   //留言文件
	MsgBdVideo       int8  = 14   //留言视频流
	MsgBdVoice       int8  = 15   //留言语音流

	MsgSysRTPosUp      int8  = 21    //位置更新信息的实时消息,Bo=经度,纬度,海拔
	MsgSysRTInTeam     int8  = 22    //成员进入Team的实时消息,Bo=nil
	MsgSysRTOutTeam    int8  = 23    //成员离开Team的实时消息,Bo=nil
	MsgSysRTNewMem     int8  = 24    //新成员加入的实时消息,Bo=itemtype,itemid
	MsgSysRTDisTeam    int8  = 25    //Team被解散, Bo=nil

	MsgSysBdJoinTeam   int8  = 31    //加入Team请求的留言,Bo=加入说明
	MsgSysBdAdded      int8  = 32    //通知被加入的成员留言,Bo=TeamID
	MsgSysBdExitTeam   int8  = 33    //成员退出team的留言,Bo=nil
	MsgSysBdDisTeam    int8  = 34    //Team被解散, Bo=nil
	MsgSysBdRejectJoin int8  = 35    //拒绝加入Team的申请
)

//查找到的Team信息
type TeamPubInfo struct{
	ID    int64   //TeamID
	Pu    int8    //public公开度
	Aw    int8     //Alw加入许可
	Na    string   //Name
	Tg    string   //Tag
	De    string   //Description
}

//Team中Item的状态
type ItemStat struct{
	Mt  int16    //Mem Type
	Md  int64    //Mem ID
	Lt  int16    //直属上级的类型Type
	Ld  int64    //直属上级的ID
	Rl  int16    //Role
	St  int8     //Stat 
	Ar  int8     //Attr
	Tg  string    //Tag
	De  string    //Des
	Lo  float32   //经度
	La  float32   //纬度
	He  float32   //海拔
	Up  string    //更新时间
}

//Msg表中的实时消息
type MsgRT struct{
	Si    int64         //发送者ID
	Ri    int64         //接受者ID
	St    int16         //发送者类型
	Rt    int16         //接收者类型
	Mt    int8          //消息类型
	Bo    string         //消息体
}

//留言板(MsgBoard)上的消息
type MsgBd struct{
	ID   int64     //消息ID, 用来确定消息已读
	Si   int64     //发送者ID
	St   int16     //发送者类型
	Mt   int8      //消息类型
	Bo   string     //消息体
	Ti   string     //留言时间
}

//Team的敏感信息
type TeamSensi struct{
	ID int64  //TeamID
	Pr string //JoinCode Prompt
	Co string //JoinCode
}
