//Copyright 2014. All rights reserved.
//Author: lja  
//Createtime: 2014/08/03 08:27:44  Lastchange: 2014/08/13 13:09:49
//changlog:  1. create by lja

package zainar

import (
	"os"
)


//将slice2添加到slice1上,返回slice1
//因为发现append只能添加成员，不能添加slice,和godoc中描述不一致
func AppendSlice(slice1,slice2 []byte)([]byte) {
	n := len(slice2)
	for i:=0; i <n; i++{
		slice1 = append(slice1, slice2[i])
	}
	return slice1
}

const(
	fileTeamPubInfo  int16 = 1
)

//在name中写入数据, 原先数据被替换
func WriteStaticFile(ftype int16, name string, content []byte) error{
	path:= "/opt/zainar/src/loginserver/www/static/"

	switch ftype {
	case fileTeamPubInfo:
		name = path + "/pubteaminfo/" + name;
	default:
		m.Print(errUnFi)
		return errUnFi
	}

	fp,err := os.Create(name)
	if err != nil{
		m.Print(err)
		return err
	}

	_, err = fp.Write(content)
	if err != nil{
		m.Print(err)
		return err
	}
	err = fp.Close()
	if err != nil{
		m.Print(err)
		return err
	}
	return nil
}
