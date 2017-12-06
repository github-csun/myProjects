package dictx

import (
	//	"container/list"
	//"errors"
	"fmt"
	//"regexp"
	//"strconv"
	"strings"
)

type IndexDef struct {
	IndexName string
	Unique    bool
	ColNames  []string
	Cols      []*FieldDef `json:"-"`
	Tbl       *TblDef     `json:"-"`
}

type PKeyDef struct {
	IndexDef
}

func (ind *IndexDef) GetIndexDefStr() (inDef string) {
	var words []string = make([]string, 0)
	var cols []string = make([]string, 0)
	if ind.Unique {
		words = append(words, "UNIQUE")
	}
	words = append(words, "INDEX")
	words = append(words, ind.IndexName)
	words = append(words, "(")
	for _, fldname := range ind.ColNames {
		cols = append(cols, fldname)
	}
	words = append(words, strings.Join(cols, ", "))
	words = append(words, ")")
	inDef = strings.Join(words, " ")
	fmt.Println(inDef)
	return
}

func (ind *IndexDef) GetPKDefStr() (inDef string) {
	var words []string = make([]string, 0)
	var cols []string = make([]string, 0)
	words = append(words, "PRIMARY KEY")
	words = append(words, "(")
	for _, fld := range ind.Cols {
		cols = append(cols, fld.Name)
	}
	words = append(words, strings.Join(cols, ", "))
	words = append(words, ")")
	inDef = strings.Join(words, " ")
	fmt.Println(inDef)
	return
}
