package main

import (
	. "dictx"
	. "dictx/client"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"routex"

	"strings"
)

func NewTblDisp(tbl *TblDef) (rlt *TblDisp) {
	if tbl == nil {
		rlt = nil
	} else {
		rlt = new(TblDisp)
		rlt.TblName = tbl.TblName
		rlt.Fields = make([]FieldDef, len(tbl.Fields))
		rlt.Indexes = make([]IndexDef, len(tbl.Indexes))

	}
	return
}

type tableCtrl struct {
	routex.Context
}

func (this *tableCtrl) New() (obj routex.Controller) {
	obj = new(tableCtrl)
	//this.Ctx.SetCtx(nil)
	return
}

func (this *tableCtrl) Get(w http.ResponseWriter, r *http.Request) {
	GetTable(w, r)
	return
}
func (this *tableCtrl) Post(w http.ResponseWriter, r *http.Request) {
	PostTable(w, r)
	return
}
func (this *tableCtrl) Put(w http.ResponseWriter, r *http.Request) {
	PutTable(w, r)
	return
}
func (this *tableCtrl) Delete(w http.ResponseWriter, r *http.Request) {
	DeleteTable(w, r)
	return
}

func PostTable(w http.ResponseWriter, r *http.Request) {
	log.Println(r)
	log.Println(r.Body)
	rt := new(RT)
	name := new(ItemName)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&name)
	log.Println(name.Name)
	if _, ok := tbls[strings.ToUpper(name.Name)]; ok {
		rt.Code = -1
		rt.Msg = fmt.Sprintf("Failed, table \"%s\" already exists.", name.Name)
		rt.Rlt = nil
	} else {
		log.Println(name)
		tbl := NewTblDef(name.Name)
		tbls[tbl.TblName] = tbl
		rt.Code = 0
		rt.Msg = fmt.Sprintf("Table \"%s\" created.", name.Name)
		rt.Rlt = nil
	}
	log.Println(rt)
	encoder := json.NewEncoder(w)
	encoder.Encode(rt)
	return
}
func GetTable(w http.ResponseWriter, r *http.Request) {
	rt := new(RTQT)
	r.ParseForm()

	tn := r.FormValue("tblName")
	if tn != "" {
		tbl, ok := tbls[tn]
		if ok {
			rt.Code = 0
			rt.Msg = "OK."
			rt.Rlt = []TblDisp{*NewTblDisp(tbl)}
		} else {
			rt.Code = -1
			rt.Msg = fmt.Sprintf("Table \"%s\" doesn't exist.", tn)
		}
	} else {
		var tables []TblDisp
		tables = make([]TblDisp, len(tbls))
		i := 0
		for _, tbl := range tbls {
			tables[i] = *NewTblDisp(tbl)
			i = i + 1
		}
		rt.Rlt = tables
		rt.Code = 0
		rt.Msg = "OK"
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(rt)
	return
}

func DeleteTable(w http.ResponseWriter, r *http.Request) {
	rt := new(RT)
	name := new(ItemName)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&name)
	log.Println(name.Name)
	if _, ok := tbls[strings.ToUpper(name.Name)]; ok {
		delete(tbls, strings.ToUpper(name.Name))
		rt.Code = 0
		rt.Msg = fmt.Sprintf("Table \"%s\" deleted.", name.Name)

		rt.Rlt = nil
	} else {
		rt.Code = -1
		rt.Msg = fmt.Sprintf("Failed, table \"%s\" doesn't exists.", name.Name)
		rt.Rlt = nil
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(rt)
	return
}

func PutTable(w http.ResponseWriter, r *http.Request) {
	rt := new(RT)

	var req = TableUpdateRequest{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&req)
	log.Println(req)
	if tbl, ok := tbls[req.OriginName]; ok {
		delete(tbls, req.OriginName)
		tbl.TblName = req.NewName
		tbls[tbl.TblName] = tbl
		rt.Code = 0
		rt.Msg = fmt.Sprintf("Table name \"%s\" updated.", req.NewName)

		rt.Rlt = nil
	} else {

		rt.Code = -1
		rt.Msg = fmt.Sprintf("Failed, table \"%s\" doesn't exists.", req.OriginName)
		rt.Rlt = nil
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(rt)
	return
}
