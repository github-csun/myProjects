package dictx

import (
	//	"container/list"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

func init() {
	switch runtime.GOOS {
	case "darwin", "freebsd", "linux":
		path = "/tmp"
	default:
		path = "C:\\"
	}
	packageName = "test"
}

var (
	path        string
	packageName string
	FuncMap     = template.FuncMap{
		"GetGoDefStr":         GetGoDefStr,
		"GetInsertPrepList":   GetInsertPrepList,
		"GetInsertExecParams": GetInsertExecParams,
		"GetUpdatePrepList":   GetUpdatePrepList,
		"GetPKConditions":     GetPKConditions,
		"GetUpdateExecParams": GetUpdateExecParams,
		"GetDeleteExecParams": GetDeleteExecParams,
		"GetPKParamDecl":      GetPKParamDecl,
		"GetSelectList":       GetSelectList,
		"GetPKPExecParams":    GetPKPExecParams,
		"GetScanParams":       GetScanParams,
		"GetIndexParamDecl":   GetIndexParamDecl,
		"GetIndexConditions":  GetIndexConditions,
		"GetIndexExecParams":  GetIndexExecParams,
		"NeedMysqlPackage":    NeedMysqlPackage,
	}
)

type ExecuteObj struct {
	*TblDef
	PkgName string
}

func NeedMysqlPackage(args ...*TblDef) string {
	no := ""
	yes := "github.com/go-sql-driver/mysql"
	f := args[0]
	for _, fld := range f.Fields {
		switch fld.FieldType {
		case DATETIME, TIMESTAMP, DATE:
			return yes
		default:
			continue
		}
	}
	return no
}

func GetGoDefStr(args ...*FieldDef) (fieldDef string) {
	f := args[0]
	if t, ok := goTypeMap[f.FieldType]; ok {
		fieldDef = f.Name + " " + t
	} else {
		panic(strconv.FormatUint(uint64(f.FieldType), 10) + "is incorrect")
	}
	return
}

func GetInsertPrepList(args ...*TblDef) (str string) {
	f := args[0]
	cols := make([]string, 0)
	for _, fld := range f.Fields {
		if fld.GeneratedField {
			continue
		} else {
			cols = append(cols, fld.Name+"=?")
		}
	}
	sort.Strings(cols)
	str = strings.Join(cols, ", ")
	return
}

func GetInsertExecParams(args ...*TblDef) (str string) {
	f := args[0]
	cols := make([]string, 0)
	for _, fld := range f.Fields {
		if fld.GeneratedField {
			continue
		} else {
			cols = append(cols, "note."+fld.Name)
		}
	}
	sort.Strings(cols)
	str = strings.Join(cols, ", ")
	return
}

func GetUpdatePrepList(args ...*TblDef) (str string) {
	f := args[0]
	cols := make([]string, 0)
	for _, fld := range f.Fields {
		if fld.GeneratedField {
			continue
		} else {
			cols = append(cols, fld.Name+"=?")
		}
	}
	sort.Strings(cols)
	str = strings.Join(cols, ", ")
	return
}

func GetPKConditions(args ...*TblDef) (str string) {
	f := args[0]
	cols := make([]string, 0)
	for _, fld := range f.Pk.Cols {
		cols = append(cols, fld.Name+"=?")
	}
	str = strings.Join(cols, " AND ")
	return
}

func GetUpdateExecParams(args ...*TblDef) (str string) {
	f := args[0]
	cols := make([]string, 0)
	for _, fld := range f.Fields {
		if fld.GeneratedField {
			continue
		} else {
			cols = append(cols, "newValue."+fld.Name)
		}
	}
	sort.Strings(cols)
	for _, fld := range f.Pk.Cols {
		cols = append(cols, "originValue."+fld.Name)
	}
	str = strings.Join(cols, ", ")
	return
}

func GetDeleteExecParams(args ...*TblDef) (str string) {
	f := args[0]
	cols := make([]string, 0)
	for _, fld := range f.Pk.Cols {
		cols = append(cols, "value."+fld.Name)
	}
	str = strings.Join(cols, ", ")
	return
}

func GetPKParamDecl(args ...*TblDef) (str string) {
	f := args[0]
	cols := make([]string, 0)
	for _, fld := range f.Pk.Cols {
		s := fmt.Sprintf("%v %v", strings.ToLower(fld.Name),
			goTypeMap[fld.FieldType])
		cols = append(cols, s)
	}
	str = strings.Join(cols, ", ")
	return
}

func GetSelectList(args ...*TblDef) (str string) {
	f := args[0]
	cols := make([]string, 0)
	for _, fld := range f.Fields {
		cols = append(cols, fld.Name)
	}
	sort.Strings(cols)
	str = strings.Join(cols, ", ")
	return
}
func GetScanParams(args ...*TblDef) (str string) {
	f := args[0]
	cols := make([]string, 0)
	for _, fld := range f.Fields {
		cols = append(cols, "&note."+fld.Name)
	}
	sort.Strings(cols)
	str = strings.Join(cols, ", ")
	return
}

func GetPKPExecParams(args ...*TblDef) (str string) {
	f := args[0]
	cols := make([]string, 0)
	for _, fld := range f.Pk.Cols {
		cols = append(cols, strings.ToLower(fld.Name))
	}
	str = strings.Join(cols, ", ")
	return
}

func GetIndexParamDecl(args ...*IndexDef) (str string) {
	f := args[0]
	cols := make([]string, 0)
	for _, fld := range f.Cols {
		s := fmt.Sprintf("%v %v", strings.ToLower(fld.Name),
			goTypeMap[fld.FieldType])
		cols = append(cols, s)
	}
	str = strings.Join(cols, ", ")
	return
}

func GetIndexConditions(args ...*IndexDef) (str string) {
	f := args[0]
	cols := make([]string, 0)
	for _, fld := range f.Cols {
		cols = append(cols, fld.Name+"=?")
	}
	str = strings.Join(cols, " AND ")
	return
}

func GetIndexExecParams(args ...*IndexDef) (str string) {
	f := args[0]
	cols := make([]string, 0)
	for _, fld := range f.Cols {
		cols = append(cols, strings.ToLower(fld.Name))
	}
	str = strings.Join(cols, ", ")
	return
}

func GenerateGo(tableDef *TblDef) (err error) {
	tpl := template.New("v2.tpl")
	tpl = tpl.Funcs(FuncMap)
	tpl, err = tpl.ParseFiles("v2.tpl")
	if err != nil {
		return
	} else {
		e := new(ExecuteObj)
		e.TblDef = tableDef
		e.PkgName = packageName
		pathName := fmt.Sprintf("%s%c%s%s", path, os.PathSeparator, tableDef.TblName, ".go")
		file, _ := os.Create(pathName)
		defer file.Close()
		err = tpl.Execute(file, e)
	}
	return
}

/*
Set generated file's destination path and package name
*/
func Setting(pathName string, pkgName string) (err error) {
	err = os.MkdirAll(pathName, os.ModeDir|os.ModePerm)
	if err == nil {
		path = pathName
	}
	packageName = pkgName
	return
}
