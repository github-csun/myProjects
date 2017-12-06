package dictx

import (
	"fmt"
)

const (
	OK = "OK."
)

type ItemName struct {
	Name string
}

type R struct {
	Code int
	Msg  string
}

func (this *R) PMsg() {
	fmt.Printf("Code=%d, %s\n", this.Code, this.Msg)
}

type RT struct {
	R
	Rlt interface{}
}

type RTQT struct {
	R
	Rlt []TblDisp
}

type RTFLD struct {
	R
	Rlt []Field
}

type Err struct {
	code    int
	message string
}

func (e *Err) Error() string {
	return e.message
}

func NewErr(code int, msg string) (err *Err) {
	err = new(Err)
	err.code = code
	err.message = msg
	return
}

type Field struct {
	TableName      string
	Name           string
	FieldType      string
	Length         uint8
	Decimals       uint8
	Nullable       bool
	DefaultV       string
	AutoIncre      bool
	GeneratedField bool
	Comment        string
	//GoType         string
	//reserved
	//flag fieldFlag
}

func (f *Field) CopyField(fld *FieldDef) {
	f.Name = fld.Name
	f.FieldType, _ = GetSQLTypeName(fld.FieldType)
	f.Length = fld.Length
	f.Decimals = fld.Decimals
	f.Nullable = fld.Nullable
	f.DefaultV = fld.DefaultV
	f.AutoIncre = fld.AutoIncre
	f.GeneratedField = fld.GeneratedField
	f.Comment = fld.Comment
}

type TblDisp struct {
	TblName string
	Indexes []IndexDef
	Fields  []FieldDef
	Pk      IndexDef
}

/*
Used for define Index. Client use it send parameters to server
*/
type IndexParam struct {
	Name       string
	IsUnique   bool
	FieldNames []string
}

type RTQI struct {
	R
	Rlt []IndexParam
}

type RTCT struct {
	R
	Rlt string
}

type RTSTR struct {
	R
	Rlt string
}
