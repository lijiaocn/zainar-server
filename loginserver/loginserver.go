//Copyright 2014. All rights reserved.
//Author: lja  
//Createtime: 2014/07/24 10:53:59  Lastchange: 2014/08/10 08:54:44
//changlog:  1. create by lja

package main

import (
	"net/http"
	"log"
	"html/template"
	"flag"
	"os"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"zainar"
)

func login(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET" {
		t, err := template.ParseFiles(zainar.Uri_human_login)
		if err != nil {
			log.Print(err)
			return
		}
		t.Execute(w, nil)
		return
	}else{
		r.ParseForm()
		pwds, ok := r.Form["pwd"]
		if ok == false || len(pwds) != 1{
			_, err := w.Write([]byte(zainar.IglForm))
			if err != nil{
				log.Print(err)
				return
			}
			return
		}
		mails, ok := r.Form["mail"]
		if ok == false || len(mails) != 1{
			_, err := w.Write([]byte(zainar.IglForm))
			if err != nil{
				log.Print(err)
				return
			}
			return
		}
		types, ok := r.Form["type"]
		if ok == false || len(types) != 1{
			_, err := w.Write([]byte(zainar.IglForm))
			if err != nil{
				log.Print(err)
				return
			}
			return
		}

		pwd := pwds[0]
		mail := mails[0]
		itemname := types[0]

		if zainar.RegxMail.MatchString(mail) == false{
			_, err := w.Write([]byte(zainar.IglMail))
			if err != nil{
				log.Print(err)
				return
			}
			return
		}

		var itemtype int16

		err := zainar.Db.QueryRow("select ID from Item where Name=?",itemname).Scan(&itemtype)

		switch {
		case err == sql.ErrNoRows:
			_, err = w.Write([]byte(zainar.ItemNotExist))
			if err != nil{
				log.Print(err)
				return
			}
		case err != nil:
			log.Print(err)
			return
		}

		switch {
		case itemname == "Human":
			err = zainar.HumanLogin(w, r, pwd, mail, itemtype)
			if err != nil{
				log.Print(err)
				_, err = w.Write([]byte(zainar.InternalErr))
				if err != nil{
					log.Print(err)
					return
				}
			}
			return

		case itemname == "VirtualTrackerA1":
			return

		default:
			_, err = w.Write([]byte(zainar.ItemNotExist))
			if err != nil{
				log.Print(err)
				return
			}
		}
	}
}

func register(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET" {
		t, err := template.ParseFiles(zainar.Uri_human_regist)
		if err != nil {
			log.Print(err)
			return
		}
		t.Execute(w, nil)
	}else{
		r.ParseForm()

		//TODO: Check Input
		pwd := r.Form["pwd"][0]
		nick := r.Form["nick"][0]
		mail    := r.Form["mail"][0]
		itemname := r.Form["type"][0]

		var itemtype int16

		err := zainar.Db.QueryRow("select ID from Item where Name=?", itemname).
			Scan(&itemtype)

		switch {
		case err == sql.ErrNoRows:
			_, err = w.Write([]byte(zainar.ItemNotExist))
			if err != nil{
				log.Print(err)
				return
			}
		case err != nil:
			log.Print(err)
		}

		switch {
		case itemname == "Human":
			err = zainar.HumanRegist(w, r, pwd, mail, nick, itemtype)
			if err != nil{
				log.Print(err)
				_, err = w.Write([]byte(zainar.InternalErr))
				if err != nil{
					log.Print(err)
					return
				}
			}
			return

		case itemname == "VirtualTrackerA1":
			//TODO:
			return
		case itemname == "VirtualTrackerA2":
			//TODO:
			return

		default:
			_, err = w.Write([]byte(zainar.ItemNotExist))
			if err != nil{
				log.Print(err)
				return
			}
		}
	}
}

func index(w http.ResponseWriter, r *http.Request){
	t, err := template.ParseFiles(zainar.Uri_index)
	if err != nil {
		log.Print(err)
		return
	}
	t.Execute(w, nil)
}

func main() {
	log.Print("Start ", os.Args)

	www  := flag.String("www","www"," web template path")
	addr := flag.String("addr","192.168.88.130:443","Server Listen Address")
	openmax := flag.Int("openmax",12000,"Database max open connections")
	idlemax := flag.Int("idlemax",10000,"Database max idle connections")
	datasource := flag.String("datasource","shangwei:123456@tcp(127.0.0.1:3306)/shangwei?"+
		"charset=utf8&collation=utf8_general_ci","Database Source Address")
	flag.Parse()

	//设置Web路径
	err := zainar.SetWWWPath(*www)
	if err != nil {
		log.Fatal(err)
	}

	//连接数据库
	defer zainar.CloseDatabase()
	err = zainar.SetDatabase("mysql", *datasource, *openmax, *idlemax)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(*www+"/static/"))))

	server := &http.Server{Addr: *addr, Handler: nil}
	server.SetKeepAlivesEnabled(false)

	//err = http.ListenAndServeTLS(*addr, "server.crt", "server.key", nil)
	err = server.ListenAndServe()

	if err != nil {
		log.Print(err)
	}
}
