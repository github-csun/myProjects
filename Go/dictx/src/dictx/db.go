package dictx

import (
	//	"container/list"
	//"errors"
	"fmt"
	//"regexp"
	//"strconv"
	//"strings"
)

type Database struct {
	tables map[string]*TblDef
}

func (this *Database) SetTables(tables map[string]*TblDef) {
	this.tables = tables
}

func (this *Database) GetTable(name string) (table *TblDef, err error) {
	if t, ok := this.tables[name]; ok == false {
		err = fmt.Errorf("Table \"%s\" does not exist.")
	} else {
		table = t
	}
	return
}

func (this *Database) DefineIndex(tableName, indexName string, isUnique bool,
	fieldNames []string) (err error) {

	var tbl *TblDef
	if tbl, err = this.GetTable(tableName); err == nil {
		_, err = tbl.AddIndex(indexName, isUnique, fieldNames)
	}
	return
}
func (this *Database) GetIndexes(tableName,
	indexName string) (indexes []IndexDef, err error) {

	var tbl *TblDef
	if tbl, err = this.GetTable(tableName); err == nil {
		indexes, err = tbl.GetIndexes(indexName)
	}
	return
}

func (this *Database) DeleteIndex(tableName,
	indexName string) (err error) {

	var tbl *TblDef
	if tbl, err = this.GetTable(tableName); err == nil {
		err = tbl.DeleteIndex(indexName)
	}
	return
}

func (this *Database) GetTableCreationSql(tableName string) (str string, err error) {
	var tbl *TblDef
	if tbl, err = this.GetTable(tableName); err == nil {
		str = tbl.GetCreateTableStr()
	}
	return

}

func (this *Database) GetGoStr(tableName string) (str string, err error) {
	var tbl *TblDef
	if tbl, err = this.GetTable(tableName); err == nil {
		str, err = tbl.GetGoStr()
	}
	return
}
