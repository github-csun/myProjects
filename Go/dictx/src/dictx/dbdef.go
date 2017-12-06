// Package dbutil
package dictx

import (
	//	"container/list"
	//"errors"
	"fmt"
	//"regexp"
	"strconv"
	"strings"
	//	"text/template"
)

const (
	SERIAL uint8 = iota + 1
	TINYINT
	SMALLINT
	INT
	BIGINT
	DECIMAL
	UTINYINT
	USMALLINT
	UINT
	UBIGINT
	DATETIME
	TIMESTAMP
	DATE
	CHAR
	VARCHAR
)

var (
	sqlTypeMap map[uint8]string
	goTypeMap  map[uint8]string
	fldTypeMap map[string]uint8
)

func init() {
	initMaps()
}
func initMaps() {
	goTypeMap = make(map[uint8]string)
	sqlTypeMap = make(map[uint8]string)
	fldTypeMap = map[string]uint8{
		"TINYINT":   TINYINT,
		"SMALLINT":  SMALLINT,
		"INT":       INT,
		"BIGINT":    BIGINT,
		"DECIMAL":   DECIMAL,
		"UTINYINT":  UTINYINT,
		"USMALLINT": USMALLINT,
		"UINT":      UINT,
		"UBIGINT":   UBIGINT,
		"DATETIME":  DATETIME,
		"TIMESTAMP": TIMESTAMP,
		"DATE":      DATE,
		"CHAR":      CHAR,
		"VARCHAR":   VARCHAR,
	}
	goTypeMap[TINYINT] = "int8"
	goTypeMap[SMALLINT] = "int16"
	goTypeMap[INT] = "int32"
	goTypeMap[BIGINT] = "int64"
	goTypeMap[DECIMAL] = "Decimal"
	goTypeMap[UTINYINT] = "uint8"
	goTypeMap[USMALLINT] = "uint16"
	goTypeMap[UINT] = "uint32"
	goTypeMap[UBIGINT] = "uint64"
	goTypeMap[DATETIME] = "XTime"
	goTypeMap[TIMESTAMP] = "XTime"
	goTypeMap[DATE] = "XTime"
	goTypeMap[CHAR] = "string"
	goTypeMap[VARCHAR] = "string"

	sqlTypeMap[TINYINT] = "TINYINT"
	sqlTypeMap[SMALLINT] = "SMALLINT"
	sqlTypeMap[INT] = "INT"
	sqlTypeMap[BIGINT] = "BIGINT"
	sqlTypeMap[DECIMAL] = "DECIMAL"
	sqlTypeMap[UTINYINT] = "TINYINT UNSIGNED"
	sqlTypeMap[USMALLINT] = "SMALLINT UNSIGNED"
	sqlTypeMap[UINT] = "INT UNSIGNED"
	sqlTypeMap[UBIGINT] = "BIGINT UNSIGNED"
	sqlTypeMap[DATETIME] = "DATETIME"
	sqlTypeMap[TIMESTAMP] = "TIMESTAMP"
	sqlTypeMap[DATE] = "DATE"
	sqlTypeMap[CHAR] = "CHAR"
	sqlTypeMap[VARCHAR] = "VARCHAR"
}

type FieldDef struct {
	Name           string
	FieldType      uint8
	Length         uint8
	Decimals       uint8
	Nullable       bool
	DefaultV       string
	AutoIncre      bool
	GeneratedField bool
	Comment        string
	GoType         string
	//reserved
	//flag fieldFlag
}

func GetSQLTypeName(typeId uint8) (typeName string, err *Err) {
	typeName, ok := sqlTypeMap[typeId]
	if ok == false {
		typeName = ""
		err = NewErr(1100, fmt.Sprintf("Invalid type: %d.", typeId))
	}
	return
}

func GetSQLTypeId(typeName string) (typeId uint8, err *Err) {
	typeId, ok := fldTypeMap[typeName]
	if ok == false {
		typeId = 255
		err = NewErr(1101, fmt.Sprintf("Invalid type: %s.", typeName))
	}
	return
}

func (f *FieldDef) SetName(name string) (err error) {
	f.Name = name
	return
}

func (f *FieldDef) SetFieldType(typeName string) (err error) {
	if t, ok := fldTypeMap[typeName]; ok == false {
		fmt.Print(t, ok)
		err = NewErr(1004, "Invalid field type")
	} else {
		f.FieldType = t
	}
	return
}

func (f *FieldDef) SetLength(lenth uint8) (err error) {
	return
}

func (f *FieldDef) SetDecimal(lenth uint8) (err error) {
	if f.FieldType != DECIMAL {
		err = NewErr(1005, "Decimal is not available for type "+sqlTypeMap[f.FieldType]+".")
	} else {
		if lenth >= f.Length {
			err = NewErr(1005, "Decimal is too big.")
		}
		f.Decimals = lenth
	}
	return
}

func (f *FieldDef) SetNullable(nullable bool) (err error) {
	f.Nullable = nullable
	return
}

func (f *FieldDef) SetAutoIncre(isAutoIncrease bool) (err error) {
	f.AutoIncre = isAutoIncrease
	// If AutoIncre == true, it must be a db generated value field.
	if isAutoIncrease {
		f.GeneratedField = true
	}
	return
}

func (f *FieldDef) SetGeneratedField(isGeneratedField bool) (err error) {
	if f.AutoIncre != true {
		f.GeneratedField = isGeneratedField
	} else {
		err = fmt.Errorf("Field \"%s\" is an AutoIncre field, cannot set isGeneratedField to false.", f.Name)
	}

	return
}

func (f *FieldDef) SetDefaultV(defaultValue string) (err error) {
	f.DefaultV = defaultValue
	return
}

func (f *FieldDef) SetComment(comment string) (err error) {
	f.Comment = comment
	return
}
func (f *FieldDef) Update(name, dataType string, lenth, decimal uint8, nullable bool,
	defaultValue string, autoIncrement, isGeneratedField bool, comment string) (err error) {
	if err = f.SetName(name); err != nil {
		goto RETURN
	}
	if err = f.SetFieldType(dataType); err != nil {
		goto RETURN
	}
	if err = f.SetLength(lenth); err != nil {
		goto RETURN
	}
	if err = f.SetDecimal(decimal); err != nil {
		goto RETURN
	}
	if err = f.SetNullable(nullable); err != nil {
		goto RETURN
	}
	if err = f.SetDefaultV(defaultValue); err != nil {
		goto RETURN
	}
	if err = f.SetAutoIncre(autoIncrement); err != nil {
		goto RETURN
	}
	if err = f.SetGeneratedField(isGeneratedField); err != nil {
		goto RETURN
	}
	if err = f.SetComment(comment); err != nil {
		goto RETURN
	}

RETURN:
	return
}

func NewField(name, dataType string, lenth, decimal uint8, nullable bool,
	defaultValue string, autoIncrement, isGeneratedField bool, comment string) (fld *FieldDef) {

	fld = new(FieldDef)
	fld.Name = strings.ToUpper(name)
	//if fld.FieldType, ok := fldTypeMap[dataType]; !ok {
	if t, ok := fldTypeMap[dataType]; ok == false {
		fmt.Print(t, ok)
		panic("Invalid field type :" + dataType)
	} else {
		fld.FieldType = t
	}
	fld.Length = lenth
	fld.Decimals = decimal
	fld.Nullable = nullable
	fld.DefaultV = defaultValue
	fld.AutoIncre = autoIncrement
	fld.Comment = comment
	fld.GeneratedField = isGeneratedField
	fld.GoType = goTypeMap[fld.FieldType]
	return
}

func (f *FieldDef) GetColDefStr() (colDef string) {
	var words []string = make([]string, 0)

	words = append(words, f.Name)

	if t, ok := sqlTypeMap[f.FieldType]; ok {
		switch f.FieldType {
		case CHAR, VARCHAR:
			t = t + "(" + strconv.Itoa(int(f.Length)) + ")"
		case DECIMAL:
			t = fmt.Sprintf("%s(%v,%v)", t, f.Length, f.Decimals)
			//case UTINYINT, USMALLINT, UINT, UBIGINT:
			//	t = t + "UNSIGNED"
		}
		words = append(words, t)
	} else {
		panic(strconv.FormatUint(uint64(f.FieldType), 10) + "is incorrect")
	}
	if f.Nullable == false {
		words = append(words, "NOT")
	}
	words = append(words, "NULL")
	if len(f.DefaultV) != 0 {
		words = append(words, "DEFAULT "+f.DefaultV)
	}
	if f.AutoIncre {
		words = append(words, "AUTO_INCREMENT")
	}
	if len(f.Comment) != 0 {
		words = append(words, fmt.Sprintf("COMMENT '%s'", f.Comment))
	}
	colDef = strings.Join(words, " ")
	return
}

type Decimal struct {
	n        int64
	decimals uint8
}
