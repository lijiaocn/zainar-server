//Copyright 2014. All rights reserved.
//Author: lja  
//Createtime: 2014/07/30 05:15:36  Lastchange: 2014/08/10 14:09:56
//changlog:  1. create by lja

package main

import (
	"net"
	"flag"
	"os"
	"log"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"zainar"
)

func main() {
	log.Print("Start ", os.Args)

	addr := flag.String("addr","192.168.88.130:8000","Server Listen Address")
	datasource := flag.String("datasource","shangwei:123456@tcp(127.0.0.1:3306)/shangwei?"+
		"charset=utf8&collation=utf8_general_ci","Database Source Address")
	sec := flag.Int("polltime",5,"msg pool time, seconds")
	openmax := flag.Int("openmax",12000,"Database max open connections")
	idlemax := flag.Int("idlemax",10000,"Database max idle connections")
	flag.Parse()

	//连接数据库
	defer zainar.CloseDatabase()
	err := zainar.SetDatabase("mysql", *datasource, *openmax, *idlemax)
	if err != nil {
		log.Fatal(err)
	}

	laddr, err := net.ResolveTCPAddr("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	ln, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Fatal(err)
	}

	go zainar.TimerMsgRTCast(time.Duration(*sec))  //实时消息
	go zainar.TimerMsgBdCast(time.Duration(*sec))  //留言板消息

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Print(err)
			continue
		}

		go zainar.MsgHandler(conn)
	}
}
