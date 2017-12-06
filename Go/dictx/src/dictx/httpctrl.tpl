package {{.PkgName}}

import (
	. "database/sql"
	"container/list"
    "fmt"
	. "sqlx"s
)

{{with .TblDef}}

type {{.TblName}}Note struct{
  {{range .Fields}}
  {{.|GetGoDefStr}}
  {{end}}
}

type OUT_Table{{.TblName}}Ctrl struct{
  R	
  Rlt []*{{.TblName}}Note
}

type IN_Table{{.TblName}}Ctrl struct{
  Param {{.TblName}}Note
}

type Table{{.TblName}}Ctrl struct {
	Context
}

func (this *Table{{.TblName}}Ctrl) New() (obj Controller) {
	obj = new(Table{{.TblName}}Ctrl)
	return
}

func (this *IndexCtrl) Get(w http.ResponseWriter, r *http.Request) {
	output := new(OUT_Table{{.TblName}}Ctrl)
	tabName := this.GetPathParam("table_name")
    conn, err = GetConn()
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

{{end}}