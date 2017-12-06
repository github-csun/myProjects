// tableSql.go
package main

import (
	. "dictx"
	. "dictx/client"
	"log"
	"net/http"
	"routex"
)

type tableSqlCtrl struct {
	routex.Context
}

func (this *tableSqlCtrl) New() (obj routex.Controller) {
	obj = new(tableSqlCtrl)
	return
}

func (this *tableSqlCtrl) Get(w http.ResponseWriter, r *http.Request) {
	tabName := this.GetPathParam("table_name")
	o := new(RTCT)
	str, err := db.GetTableCreationSql(tabName)
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
func (this *tableSqlCtrl) Post(w http.ResponseWriter, r *http.Request) {
	//PostTable(w, r)
	return
}
func (this *tableSqlCtrl) Put(w http.ResponseWriter, r *http.Request) {
	//PutTable(w, r)
	return
}
func (this *tableSqlCtrl) Delete(w http.ResponseWriter, r *http.Request) {
	//DeleteTable(w, r)
	return
}
