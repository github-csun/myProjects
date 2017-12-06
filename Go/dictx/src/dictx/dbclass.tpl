package dbx

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	
)

type {{.TblName}}Note struct{
  {{range .Fields}}
  {{.|GetGoDefStr}}
  {{end}}
}

func (note {{.TblName}}Note)Clone()(result {{.TblName}}Note){
  result = new({{.TblName}}Note)
  {{range .Fields}}
  result.{{.Name}} = rec.{{.Name}}
  {{end}}
  return
}

type {{.TblName}}Record struct{
  note []{{.TblName}}Note
  conn *DB
}

func New{{.TblName}}Record(conn *DB)(result *{{.TblName}}Record){
  result = new({{.TblName}}Record)
  return
}

func (rec *{{.TblName}}Record)Clone()(result *{{.TblName}}Record){
  result = new({{.TblName}}Record)
  {{with .Fields}}
  {{range .}}
  result.{{.Name}} = rec.{{.Name}}
  {{end}}
  {{end}}
  return
}

func QueryByPK( {{.|GetPKParamDecl}})(rec *{{.TblName}}Record,err){
  stmt, err := conn.Prepare("select {{.|GetSelectList}} from {{.TblName}} where {{.|GetPKConditions}} "")
  checkErr(err)

  rows, err := stmt.Exec({{.|GetPKPExecParams}})
  for rows.Next(){
    err = rows.Scan({{.|GetScanParams}})
  }
}

func (rec *{{.TblName}}Record) Insert(tx *Tx) (err error) {

	stmt, err := tx.Prepare("INSERT {{.TblName}} SET {{.|GetInsertPrepList}} ")
	checkErr(err)

	res, err := stmt.Exec({{.|GetInsertExecParams}})

	id, err := res.LastInsertId()
	
	return
}

func (rec *{{.TblName}}Record) Update(tx *Tx, newValue *{{.TblName}}Record) error {

	stmt, err = tx.Prepare("UPDATE {{.TblName}} SET {{.|GetUpdatePrepList}} where {{.|GetPKConditions}}")
	checkErr(err)

	_, err = stmt.Exec({{.|GetUpdateExecParams}})
	checkErr(err)

	return err
}

func (rec *{{.TblName}}Record) Delete(tx *Tx) error {

	stmt, err = tx.Prepare("delete from {{.TblName}} where {{.|GetPKConditions}}")
	checkErr(err)

	_, err = stmt.Exec({{.|GetDeleteExecParams}})
	checkErr(err)

	return err
}


{{range .Indexes}}
func QueryBy{{.IndexName}}(conn *DB, {{.|GetIndexParamDecl}})(recs []*{{.Tbl.TblName}}Record,err error){
  
  stmt, err := conn.Prepare("select {{.Tbl|GetSelectList}} from {{.Tbl.TblName}} where {{.|GetIndexConditions}}")
  checkErr(err)

  rows, err := stmt.Exec({{.|GetIndexExecParams}})
  for rows.Next(){
    rec = new({{.Tbl.TblName}}Record)
    err = rows.Scan({{.Tbl|GetScanParams}})
    recs=append(recs,&rec)
  }
}
{{end}}