package main

import (
	. "dictx"
	. "dictx/client"
	"fmt"
	"log"
	"net/http"
	. "routex"
)

type IndexCtrl struct {
	Context
}

func (this *IndexCtrl) New() (obj Controller) {
	obj = new(IndexCtrl)
	return
}
func (this *IndexCtrl) Get(w http.ResponseWriter, r *http.Request) {
	output := new(RTQI)
	tabName := this.GetPathParam("table_name")
	indName := this.GetPathParam("index_name")
	inds, err := db.GetIndexes(tabName, indName)
	if err != nil {
		output.Code = 1200
		output.Msg = err.Error()
	} else {
		output.Code = 0
		output.Msg = "OK"
		tmp := make([]IndexParam, len(inds))
		for i, ind := range inds {
			tmp[i].Name = ind.IndexName
			tmp[i].IsUnique = ind.Unique
			tmp[i].FieldNames = ind.ColNames
		}
		output.Rlt = tmp
	}
	JsonEncode(w, output)
	return
}
func (this *IndexCtrl) Post(w http.ResponseWriter, r *http.Request) {
	log.Println("Define Index")
	input := new(IndexParam)
	output := new(R)
	tabName := this.GetPathParam("table_name")
	indName := this.GetPathParam("index_name")
	err := JsonDecode(r.Body, input)
	if err == nil {
		err = db.DefineIndex(tabName, indName, input.IsUnique, input.FieldNames)
		if err == nil {
			output.Code = 0
			output.Msg = "Index was defined successfully."

		} else {
			output.Code = 1100
			output.Msg = fmt.Sprintf("Defined index failed. %s", err.Error())
		}
	} else {
		output.Code = 1101
		output.Msg = err.Error()
	}
	log.Println(output.Msg)
	JsonEncode(w, output)
	return
}
func (this *IndexCtrl) Put(w http.ResponseWriter, r *http.Request) {

	return
}
func (this *IndexCtrl) Delete(w http.ResponseWriter, r *http.Request) {

	output := new(R)
	tabName := this.GetPathParam("table_name")
	indName := this.GetPathParam("index_name")
	log.Printf("Request to delete index %s from %s", indName, tabName)
	err := db.DeleteIndex(tabName, indName)
	if err == nil {
		output.Code = 0
		output.Msg = fmt.Sprintf("Index \"%s\" was deleted successfully.", indName)

	} else {
		output.Code = 1100
		output.Msg = fmt.Sprintf("Delete index failed. %s", err.Error())
	}
	log.Println(output.Msg)
	JsonEncode(w, output)
	return
}
