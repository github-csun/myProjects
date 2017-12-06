package dictx

import (
	"os"
	. "testing"
	"text/template"
)

//func Test_GetColDefStr(t *T) {
//	//func NewField(name, dataType string, lenth, decimal uint8, nullable bool,
//	//	defaultValue string, autoIncrement bool, comment string) (fld *FieldDef)
//	tmp := NewField("id", "INT", 0, 0, false, "", true, "test")
//	str := tmp.GetColDefStr()
//	if len(str) == 0 {
//		t.Error("Failed!")
//	} else {
//		t.Log(str)
//	}

//	tmp = NewField("CreateTime", "TIMESTAMP", 0, 0, true, "CURRENT_TIMESTAMP", false, "test")
//	str = tmp.GetColDefStr()
//	if len(str) == 0 {
//		t.Error("Failed!")
//	} else {
//		t.Log(str)
//	}

//	tmp = NewField("TEXT", "CHAR", 250, 0, true, "", false, "test")
//	str = tmp.GetColDefStr()
//	if len(str) == 0 {
//		t.Error("Failed!")
//	} else {
//		t.Log(str)

//	}

//}

func Test_TBL(t *T) {
	tbl := NewTblDef("test")
	tbl.AddField("id1", "INT", 0, 0, false, "", true, true, "test")
	tbl.AddField("id2", "UINT", 0, 0, false, "", false, false, "test")
	tbl.AddField("id3", "INT", 0, 0, false, "", false, false, "test")
	tbl.AddField("CreateTime", "TIMESTAMP", 0, 0, true, "CURRENT_TIMESTAMP", false, true, "test")
	tbl.AddField("TEXT", "CHAR", 250, 0, true, "", false, false, "test")
	tbl.AddPK([]string{"id1", "TEXT"})
	_, err := tbl.AddIndex("INDEX1", false, []string{"id1", "CreateTime"})
	if err == nil {
		t.Log("Add Index pass.")
	} else {
		t.Error("Add Index failed.")
	}
	tbl.AddIndex("INDEX2", true, []string{"id2", "id3"})
	t.Log(tbl.GetCreateTableStr())
	tpl := template.New("v2.tpl")
	tpl = tpl.Funcs(FuncMap)
	tpl, err = tpl.ParseFiles("v2.tpl")
	if err != nil {
		t.Error("ParseFiles Failed.", err)
	} else {
		e := new(ExecuteObj)
		e.TblDef = tbl
		f, _ := os.Create("test.go")
		err = tpl.Execute(f, e)
		f.Close()
		if err != nil {
			t.Error("Execute Failed.", err)
		}
	}
	if err = Setting("/Users/zuora/Store/dbf/ttt", "test"); err != nil {
		t.Error(err)
	} else {
		GenerateGo(tbl)
	}
}
