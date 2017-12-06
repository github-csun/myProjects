package {{.PkgName}}

import (
	. "database/sql"
	"container/list"
    "fmt"
	. "sqlx"
)

{{with .TblDef}}

type {{.TblName}}Note struct{
  {{range .Fields}}
  {{.|GetGoDefStr}}
  {{end}}
}

func (note *{{.TblName}}Note)Clone()(result *{{.TblName}}Note){
  result = new({{.TblName}}Note)
  {{range .Fields}}
  result.{{.Name}} = note.{{.Name}}
  {{end}}
  return
}

func (note *{{.TblName}}Note)Copy(src *{{.TblName}}Note){
  {{range .Fields}}
  note.{{.Name}} = src.{{.Name}}
  {{end}}
  return
}

type {{.TblName}}NoteList struct{
  list.List
}

type {{.TblName}} struct{

}

func New{{.TblName}}()(result *{{.TblName}}){
  result = new({{.TblName}})
  return
}

func (this *{{.TblName}})QueryByPK(conn *DB, {{.|GetPKParamDecl}})(note *{{.TblName}}Note, err error){
  if conn == nil{
	err = fmt.Errorf("Parameter conn is nil.")
  }else{
	  rows, err := conn.Query("select {{.|GetSelectList}} from {{.TblName}} where {{.|GetPKConditions}}", {{.|GetPKPExecParams}})
	  defer rows.Close()
	  if err == nil {
		note = new({{.TblName}}Note)
	    for rows.Next(){
	      err = rows.Scan({{.|GetScanParams}})
          if err != nil {
		    note = nil
	        CheckErr(err)
		  }
	    } 
	  }else{
		CheckErr(err)
	  }
  }
  return
}


func (this *{{.TblName}})TxQueryByPK(tx *Tx, {{.|GetPKParamDecl}})(note *{{.TblName}}Note, err error){
  if tx == nil{
	err = fmt.Errorf("Parameter tx is nil.")
  }else{
	  rows, err := tx.Query("select {{.|GetSelectList}} from {{.TblName}} where {{.|GetPKConditions}}", {{.|GetPKPExecParams}})
	  defer rows.Close()
	  if err == nil {
		note = new({{.TblName}}Note)
	    for rows.Next(){
	      err = rows.Scan({{.|GetScanParams}})
          if err != nil {
		    note = nil
	        CheckErr(err)
		  }
	    } 
	  }else{
		CheckErr(err)
	  }
  }
  return
}

func (this *{{.TblName}}) Insert(tx *Tx, note *{{.TblName}}Note) (err error) {

  if tx == nil{
	err = fmt.Errorf("Parameter tx is nil.")
  }else if note == nil{
    err = fmt.Errorf("Parameter note is nil.")
  } else{
    stmt, err := tx.Prepare("INSERT {{.TblName}} SET {{.|GetInsertPrepList}} ")
    if err == nil {
	  if _, err = stmt.Exec({{.|GetInsertExecParams}}); err != nil{
		CheckErr(err)
	  }
	} else {
	  CheckErr(err)
	}
  }

  return
}

func (this *{{.TblName}}) InsertBat(tx *Tx, notes *{{.TblName}}NoteList) (err error) {
  if tx == nil{
	err = fmt.Errorf("Parameter tx is nil.")
  }else if notes == nil{
    err = fmt.Errorf("Parameter notes is nil.")
  } else{
    stmt, err := tx.Prepare("INSERT {{.TblName}} SET {{.|GetInsertPrepList}} ")
    if err == nil{
      var note *{{.TblName}}Note
      for e := notes.Front(); e != nil; e = e.Next(){
	    note = e.Value.(*{{.TblName}}Note)
	    _, err = stmt.Exec({{.|GetInsertExecParams}})
        if err !=nil{
          CheckErr(err)
          goto RETURN
		}
      }
	}else{
	  CheckErr(err)
	}
  }
	//id, err := res.LastInsertId()
  RETURN:
  return
}

func (this *{{.TblName}}) Update(tx *Tx, newValue, originValue *{{.TblName}}Note)(err error) {
  if tx == nil{
	err = fmt.Errorf("Parameter tx is nil.")
  }else if newValue == nil{
    err = fmt.Errorf("Parameter newValue is nil.")
  }else if originValue == nil{
	err = fmt.Errorf("Parameter newValue is nil.")
  }else{
    stmt, err := tx.Prepare("UPDATE {{.TblName}} SET {{.|GetUpdatePrepList}} where {{.|GetPKConditions}}")
    if err == nil{
	  _, err = stmt.Exec({{.|GetUpdateExecParams}})	
      if err != nil{
		CheckErr(err)
	  }
	}else{
	  CheckErr(err)
	}
  }
  return 
}

func (this *{{.TblName}}) Delete(tx *Tx, value *{{.TblName}}Note)(err error) {
  if tx == nil{
	err = fmt.Errorf("Parameter tx is nil.")
  }else{
	stmt, err := tx.Prepare("delete from {{.TblName}} where {{.|GetPKConditions}}")
    if err == nil{
	  _, err = stmt.Exec({{.|GetDeleteExecParams}})
      if err !=nil{
	    CheckErr(err)
      }
	} else {
	  CheckErr(err)
    }
  }
  return 
}


{{range .Indexes}}
func (this *{{.Tbl.TblName}})QueryBy{{.IndexName}}(conn *DB, {{.|GetIndexParamDecl}})(notes *{{.Tbl.TblName}}NoteList,err error){
  if conn == nil{
	err = fmt.Errorf("Parameter conn is nil.")
  }else{
    rows, err := conn.Query("select {{.Tbl|GetSelectList}} from {{.Tbl.TblName}} where {{.|GetIndexConditions}}", {{.|GetIndexExecParams}})
    defer rows.Close()
	if err == nil{
	  notes = new({{.Tbl.TblName}}NoteList)
      var note *{{.Tbl.TblName}}Note
      for rows.Next(){
        note = new({{.Tbl.TblName}}Note)
        err = rows.Scan({{.Tbl|GetScanParams}})
        if err == nil{
          notes.PushBack(note)
        }else{
		  CheckErr(err)
          goto RETURN
		}
	  }
	}else{
	  CheckErr(err)
	}
  }
  RETURN:
  return
}

func (this *{{.Tbl.TblName}})TxQueryBy{{.IndexName}}(tx *Tx, {{.|GetIndexParamDecl}})(notes *{{.Tbl.TblName}}NoteList,err error){
  if tx == nil{
	err = fmt.Errorf("Parameter tx is nil.")
  }else{
    rows, err := tx.Query("select {{.Tbl|GetSelectList}} from {{.Tbl.TblName}} where {{.|GetIndexConditions}}", {{.|GetIndexExecParams}})
    defer rows.Close()
	if err == nil{
	  notes = new({{.Tbl.TblName}}NoteList)
      var note *{{.Tbl.TblName}}Note
      for rows.Next(){
        note = new({{.Tbl.TblName}}Note)
        err = rows.Scan({{.Tbl|GetScanParams}})
        if err == nil{
          notes.PushBack(note)
        }else{
		  CheckErr(err)
          goto RETURN
		}
	  }
	}else{
	  CheckErr(err)
	}
  }
  RETURN:
  return
}

{{end}}


{{end}}