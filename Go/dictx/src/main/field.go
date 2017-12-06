package main

import (
	. "dictx"
	. "dictx/client"
	"encoding/json"
	"fmt"
	//"log"
	"net/http"
	//"os"
	. "routex"
	//"strings"
	"log"
)

type FieldCtrl struct {
	Context
}

func (this *FieldCtrl) New() (obj Controller) {
	obj = new(FieldCtrl)
	return
}

func (this *FieldCtrl) Get(w http.ResponseWriter, req *http.Request) {
	r := make([]Field, 0)
	rt := new(RTFLD)

	tblName := this.GetPathParam("table_name")
	fldName := this.GetPathParam("field_Name")
	log.Printf("Context: %v", this.Context)
	if tbl, ok := tbls[tblName]; ok {
		rt.Code = 0
		rt.Msg = "OK."
		rt.Rlt = r
		if fldName == "" {
			for _, fld := range tbl.Fields {
				f := new(Field)
				f.CopyField(fld)
				r = append(r, *f)
			}
		} else {
			if fld, ok := tbl.Fields[fldName]; ok {
				f := new(Field)
				f.CopyField(fld)
				r = append(r, *f)
			} else {
				rt.Code = 1005
				rt.Msg = fmt.Sprintf("Field %s does not exist.", fldName)
				rt.Rlt = nil
			}
		}

	} else {
		rt.Code = 1005
		rt.Msg = fmt.Sprintf("Invalid table name \"%s\".", tblName)
		rt.Rlt = nil
	}
	if rt.Code == 0 {
		rt.Rlt = r
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(rt)
	return

}
func (this *FieldCtrl) Post(w http.ResponseWriter, r *http.Request) {
	this.PostField(w, r)
	return
}
func (this *FieldCtrl) Put(w http.ResponseWriter, r *http.Request) {
	this.PutField(w, r)
	return
}
func (this *FieldCtrl) Delete(w http.ResponseWriter, r *http.Request) {
	this.deleteField(w, r)
	return
}

func (this *FieldCtrl) PutField(w http.ResponseWriter, req *http.Request) {
	r := new(Field)
	rt := new(RT)
	err := JsonDecode(req.Body, r)
	tblName := this.GetPathParam("table_name")
	if err == nil {
		tbl, ok := tbls[tblName]
		if ok {
			err = tbl.UpdateField(r.Name, r.FieldType, r.Length, r.Decimals, r.Nullable,
				r.DefaultV, r.AutoIncre, r.GeneratedField, r.Comment)
			if err == nil {
				rt.Code = 0
				rt.Msg = "OK"
			} else {
				rt.Code = 1100
				rt.Msg = err.Error()
			}
		} else {
			rt.Code = 1100
			rt.Msg = fmt.Sprintf("Table \"%s\" does not exist.", tblName)
		}
	} else {
		rt.Code = 11000
		rt.Msg = err.Error()
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(rt)
	return
}

func (this *FieldCtrl) PostField(w http.ResponseWriter, req *http.Request) {
	r := new(Field)
	rt := new(RT)
	err := JsonDecode(req.Body, r)
	tblName := this.GetPathParam("table_name")
	rt.Code = 0
	rt.Msg = "OK."
	if err == nil {
		tbl, ok := tbls[tblName]
		if ok {
			err := tbl.AddField(r.Name, r.FieldType, r.Length, r.Decimals, r.Nullable,
				r.DefaultV, r.AutoIncre, r.GeneratedField, r.Comment)
			if err == nil {
				rt.Code = 0
				rt.Msg = fmt.Sprintf("Field \"%s\" created.", r.Name)
				rt.Rlt = nil
			} else {
				rt.Code = -777
				rt.Msg = fmt.Sprintln(err)
				rt.Rlt = nil
			}

		}

	}
	encoder := json.NewEncoder(w)
	encoder.Encode(rt)
	return
}

func (this *FieldCtrl) deleteField(w http.ResponseWriter, req *http.Request) {

	rt := new(RTFLD)

	tblName := this.GetPathParam("table_name")
	fldName := this.GetPathParam("field_Name")
	log.Printf("Context: %v", this.Context)
	if tbl, ok := tbls[tblName]; ok {
		rt.Code = 0
		rt.Msg = "OK."
		rt.Rlt = nil
		err := tbl.RemoveField(fldName)
		if err != nil {
			rt.Code = 1005
			rt.Msg = err.Error()
			rt.Rlt = nil
		}
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(rt)
	return

}
