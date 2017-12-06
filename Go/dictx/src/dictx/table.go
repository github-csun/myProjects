package dictx

import (
	//	"container/list"
	"errors"
	"fmt"
	"regexp"
	//"strconv"
	"bytes"
	"log"
	"strings"
	"text/template"
)

type TblDef struct {
	TblName string
	Indexes map[string]*IndexDef
	Fields  map[string]*FieldDef
	Pk      *IndexDef
}

func NewTblDef(table_name string) (tbl *TblDef) {
	matched, _ := regexp.MatchString("^[a-zA-Z]\\w*$", table_name)
	if matched {
		tbl = new(TblDef)
		tbl.TblName = strings.ToUpper(table_name)
		tbl.Fields = make(map[string]*FieldDef)
		tbl.Indexes = make(map[string]*IndexDef)
		//tbl.Pk = new(PKeyDef)
	} else {
		tbl = nil
	}

	return
}

func (tb *TblDef) AddField(fieldName string, dataType string, lenth,
	decimal uint8, nullable bool, defaultValue string, autoIncrement,
	isGeneratedField bool, comment string) (err error) {

	f := NewField(fieldName, dataType, lenth, decimal, nullable, defaultValue,
		autoIncrement, isGeneratedField, comment)
	if _, ok := tb.Fields[fieldName]; ok {
		err = fmt.Errorf("Duplicated field \"%v\".", fieldName)
	} else {
		tb.Fields[fieldName] = f
	}
	return
}

func (tb *TblDef) UpdateField(fieldName string, dataType string, lenth,
	decimal uint8, nullable bool, defaultValue string, autoIncrement,
	isGeneratedField bool, comment string) (err error) {

	f := NewField(fieldName, dataType, lenth, decimal, nullable, defaultValue,
		autoIncrement, isGeneratedField, comment)
	if _, ok := tb.Fields[fieldName]; ok {
		tb.Fields[fieldName] = f
	} else {
		err = fmt.Errorf("Field \"%s\" does not exist", fieldName)
	}
	return
}

func (tb *TblDef) RemoveField(fieldName string) (err error) {
	if _, ok := tb.Fields[fieldName]; ok {
		delete(tb.Fields, fieldName)
	} else {
		err = fmt.Errorf("Field \"%s\" does not exist.", fieldName)
	}
	return
}

func (tb *TblDef) AddIndex(name string, isUnique bool, fieldNames []string) (ind *IndexDef, err error) {
	if strings.ToUpper(name) != "PRIMARY" {
		t := new(IndexDef)
		t.IndexName = name
		t.Unique = isUnique
		t.Tbl = tb
		t.ColNames = fieldNames
		t.Cols = make([]*FieldDef, len(fieldNames))
		for i, fld := range fieldNames {
			if col, ok := tb.Fields[fld]; ok {
				t.Cols[i] = col
			} else {
				err = errors.New("Field: " + fld + " doesn't exist")
			}
		}
		if err == nil {
			tb.Indexes[name] = t
			ind = t
		}
	} else {
		_, err = tb.AddPK(fieldNames)
	}
	return
}

func (tb *TblDef) GetIndexes(name string) (indexes []IndexDef, err error) {
	indexes = make([]IndexDef, 0)
	if name == "*" {
		for _, ind := range tb.Indexes {
			indexes = append(indexes, *ind)
		}
	} else {
		if ind, ok := tb.Indexes[name]; ok {
			indexes = append(indexes, *ind)
		} else {
			err = fmt.Errorf("Index \"%s\" does not exist.", name)
		}
	}
	if tb.Pk != nil {
		indexes = append(indexes, *tb.Pk)
	}

	return
}

func (tb *TblDef) DeleteIndex(name string) (err error) {
	if name == "*" {
		for n, _ := range tb.Indexes {
			delete(tb.Indexes, n)
		}
	} else {
		if _, ok := tb.Indexes[name]; ok {
			delete(tb.Indexes, name)
		} else {
			err = fmt.Errorf("Index \"%s\" does not exist.", name)
		}
	}

	return
}

func (tb *TblDef) AddPK(fieldNames []string) (ind *IndexDef, err error) {
	pk := new(IndexDef)
	pk.Unique = true
	pk.IndexName = "PRIMARY"
	pk.Tbl = tb
	pk.ColNames = fieldNames
	pk.Cols = make([]*FieldDef, len(fieldNames))
	for i, fld := range fieldNames {
		if col, ok := tb.Fields[fld]; ok {
			if col.Nullable == false {
				pk.Cols[i] = col
			} else {
				err = errors.New("Field: " + fld +
					" is nullable, cannot be a primary key field.")
			}
		} else {
			err = errors.New("Field: " + fld + " doesn't exist")
		}
	}
	if err == nil {
		tb.Pk = pk
	}
	return
}

func (tb *TblDef) DeletePK() (err error) {

	if tb.Pk != nil {
		tb.Pk = nil
	} else {

	}

	return
}

func (tb *TblDef) GetCreateTableStr() (str string) {
	var lines []string = make([]string, 0)
	var create_definition []string = make([]string, 0)
	//var column_def = make([]string, 0)
	//var index_def = make([]string, 0)
	for _, fld := range tb.Fields {
		create_definition = append(create_definition, fld.GetColDefStr())
	}
	if tb.Pk != nil {
		create_definition = append(create_definition, tb.Pk.GetPKDefStr())
	}
	for _, ind := range tb.Indexes {
		create_definition = append(create_definition, ind.GetIndexDefStr())
	}

	lines = append(lines, "CREATE TABLE "+tb.TblName)
	lines = append(lines, "(\n"+strings.Join(create_definition, ",\n")+"\n)")
	str = strings.Join(lines, "\n")
	//var lines []string = make([]string, 0)
	//var create_definition []string = make([]string, 0)
	//var column_def = make([]string, 0)
	//var index_def = make([]string, 0)
	//for _, fld := range tb.Fields {
	//	column_def = append(column_def, fld.GetColDefStr())
	//}

	//for _, ind := range tb.Indexes {
	//	index_def = append(index_def, ind.GetIndexDefStr())
	//}

	//create_definition = append(create_definition, column_def[])
	//create_definition = append(create_definition, tb.Pk.GetPKDefStr())
	//create_definition = append(create_definition, index_def)

	//lines = append(lines, "CREATE TABLE")
	//lines = append(lines, tb.TblName)
	//lines = append(lines, "(")
	//lines = append(lines, strings.Join(column_def, ",\n"))
	//lines = append(lines, tb.Pk.GetPKDefStr())
	//lines = append(lines, strings.Join(index_def, ",\n"))
	//lines = append(lines, ")")
	//str = strings.Join(lines, "\n")
	return
}

func (tb *TblDef) GetGoStr() (str string, err error) {
	buf := new(bytes.Buffer)
	tpl := template.New("v2.tpl")
	tpl = tpl.Funcs(FuncMap)
	tpl, err = tpl.ParseFiles("v2.tpl")
	if err == nil {
		e := new(ExecuteObj)
		e.TblDef = tb
		CheckTable(tb)
		err = tpl.Execute(buf, e)
		if err == nil {
			str = buf.String()
		}
	}
	return
}

func CheckTable(tb *TblDef) {
	for name, ind := range tb.Indexes {
		log.Printf("Check index %s\n", name)
		if name != ind.IndexName {
			log.Printf("Index Key name = \"%s\", Index name = \"%s\"\n", name, ind.IndexName)
		}
		if ind.Tbl == nil {
			log.Printf("Index Tbl pointer is nil.\n")
		}
		if ind.Tbl.TblName != tb.TblName {
			log.Printf("Index Tbl TblName is \"%s\", table name = \"%s\"\n", ind.Tbl.TblName, tb.TblName)
		}
	}
}
