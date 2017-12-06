package main

import (
	. "dictx"
	. "dictx/client"
	//"fmt"
	"log"
	"net/http"
	. "routex"
)

type GoTplCtrl struct {
	Context
}

func (this *GoTplCtrl) New() (obj Controller) {
	obj = new(GoTplCtrl)
	return
}

func (this *GoTplCtrl) Get(w http.ResponseWriter, r *http.Request) {
	tabName := this.GetPathParam("table_name")
	o := new(RTSTR)
	str, err := db.GetGoStr(tabName)
	if err == nil {
		o.Code = 0
		o.Msg = OK
		o.Rlt = str
	} else {
		o.Code = 1202
		o.Msg = err.Error()
	}
	log.Println(o.Msg)
	JsonEncode(w, o)
	return
}
func (this *GoTplCtrl) Post(w http.ResponseWriter, r *http.Request) {
	//PostTable(w, r)
	return
}
func (this *GoTplCtrl) Put(w http.ResponseWriter, r *http.Request) {
	//PutTable(w, r)
	return
}
func (this *GoTplCtrl) Delete(w http.ResponseWriter, r *http.Request) {
	//DeleteTable(w, r)
	return
}
