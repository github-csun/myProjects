// serv.go
package main

import (
	. "dictx"
	. "dictx/client"
	//"encoding/json"
	//"fmt"
	"log"
	"net/http"
	"os"
	. "routex"
	//"strings"
)

var (
	tbls map[string]*TblDef
	db   *Database
)

func main() {
	tbls = make(map[string]*TblDef)
	err := Read()
	if err != nil {
		log.Println(err)
	}
	db = new(Database)
	db.SetTables(tbls)
	for _, tbl := range tbls {
		for _, ind := range tbl.Indexes {
			ind.Tbl = tbl
			ind.Cols = make([]*FieldDef, len(ind.ColNames))
			for i, colName := range ind.ColNames {
				ind.Cols[i] = tbl.Fields[colName]
			}
		}
		if tbl.Pk != nil {
			tbl.Pk.Tbl = tbl
		}
		CheckTable(tbl)
	}
	//http.HandleFunc("/table", TableCtrl)
	RegRoute("/table/:table_name", new(tableCtrl))
	RegRoute("/data", new(dataCtrl))
	RegRoute("/table/:table_name/field/:field_Name", new(FieldCtrl))
	RegRoute("/table/:table_name/index/:index_name", new(IndexCtrl))
	RegRoute("/table/:table_name/CreationSql", new(tableSqlCtrl))
	RegRoute("/table/:table_name/go", new(GoTplCtrl))

	//err = http.ListenAndServe(":9099", nil)
	//if err != nil {
	//	print(err)
	//	log.Fatal("ListenAndServe: ", err)
	//}
	Serve(":9099")
}

type dataCtrl struct {
	Context
}

func (this *dataCtrl) New() (obj Controller) {
	obj = new(dataCtrl)
	return
}

func (this *dataCtrl) Get(w http.ResponseWriter, r *http.Request) {
	Read()
	return
}
func (this *dataCtrl) Post(w http.ResponseWriter, r *http.Request) {
	Save()
	return
}
func (this *dataCtrl) Put(w http.ResponseWriter, r *http.Request) {

	return
}
func (this *dataCtrl) Delete(w http.ResponseWriter, r *http.Request) {

	return
}

func Save() (err error) {
	file, err := os.Create("db")
	//file, err := os.OpenFile("db", os.O_RDWR&os.O_CREATE&os.O_TRUNC, os.ModePerm)
	defer file.Close()
	if err == nil {
		_, err = JsonEncode(file, tbls)
	}
	if err != nil {
		log.Println(err)
	}
	return
}

func Read() (err error) {
	file, err := os.Open("db")
	defer file.Close()
	if err == nil {
		err = JsonDecode(file, &tbls)
	}
	if err != nil {
		log.Println(err)
	}
	return
}
