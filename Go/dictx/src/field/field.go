package main

import (
	. "dictx"
	. "dictx/client"
	"flag"
	"fmt"
	"log"

	"os"
	"strings"
)

var (
	exit_code int
	fi        *Field
)
var (
	//New, Update, Query, Delete bool
	action, dataType, defaultValue, comment, fieldName string
	isNullable, isAutoInreace, isGenerate              bool
	lenth, decimal                                     uint
)

func main() {

	flag.StringVar(&fieldName, "fieldName", "", "New name for the table")
	flag.StringVar(&fieldName, "f", "", "New name for the table")
	flag.StringVar(&action, "action", "", "New name for the table")
	flag.StringVar(&action, "a", "", "New name for the table")

	flag.StringVar(&dataType, "type", "", "New name for the table")
	flag.StringVar(&dataType, "t", "", "New name for the table")
	flag.StringVar(&defaultValue, "defaultValue", "", "New name for the table")
	flag.StringVar(&defaultValue, "v", "", "New name for the table")

	flag.StringVar(&comment, "comment", "", "New name for the table")
	flag.StringVar(&comment, "c", "", "New name for the table")
	flag.BoolVar(&isNullable, "nullable", true, "New name for the table")
	flag.BoolVar(&isNullable, "n", true, "New name for the table")

	flag.BoolVar(&isAutoInreace, "autoInreace", false, "New name for the table")
	flag.BoolVar(&isAutoInreace, "auto", false, "New name for the table")
	flag.BoolVar(&isGenerate, "generate", false, "New name for the table")
	flag.BoolVar(&isGenerate, "g", false, "New name for the table")

	flag.UintVar(&lenth, "lenth", 0, "")
	flag.UintVar(&lenth, "l", 0, "")
	flag.UintVar(&decimal, "decimal", 0, "")
	flag.UintVar(&decimal, "d", 0, "")

	flag.Parse()
	args := flag.Args()
	for i, arg := range args {
		args[i] = strings.ToUpper(arg)
	}

	fi = new(Field)
	fi.TableName = args[0]
	fi.Name = fieldName
	fi.FieldType = dataType
	fi.Length = uint8(lenth)
	fi.Decimals = uint8(decimal)
	fi.Nullable = isNullable
	fi.DefaultV = defaultValue
	fi.AutoIncre = isAutoInreace
	fi.GeneratedField = isGenerate

	switch strings.ToLower(action) {
	case "c", "create":
		CreateField(fi)
	case "q", "query":
		QueryField(args[0], fieldName)
	case "d", "delete":
		deleteField(args[0], fieldName)
		//DeleteTable(args)

	case "u", "update":
		fields, err := getFields(args[0], fieldName)
		if err == nil {
			if len(fields) > 0 {
				fi = &fields[0]
				flag.Visit(setField)
				fi.TableName = args[0]
				updateField(fi)
			} else {
				log.Println("Fields is empty. ")
			}
		} else {
			log.Println(err)
		}
	}

}

func updateField(f *Field) {
	rt := new(RT)
	url := fmt.Sprintf("http://localhost:9099/table/%s/field", f.TableName)
	_, err := JsonReq("PUT", url, &f, rt)
	if err == nil {
		exit_code = rt.Code
		log.Println(rt.Msg)

	} else {
		exit_code = -1000
		log.Println(err)
	}
	os.Exit(exit_code)
}

func deleteField(table_name, field_name string) {
	rt := new(RT)
	url := fmt.Sprintf("http://localhost:9099/table/%s/field/%s", table_name, field_name)
	_, err := JsonReq(DELETE, url, nil, rt)
	if err == nil {
		exit_code = rt.Code
		log.Println(rt.Msg)

	} else {
		exit_code = -1000
		log.Println(err)
	}
	os.Exit(exit_code)
}

func CreateField(f *Field) {
	rt := new(RT)
	url := fmt.Sprintf("http://localhost:9099/table/%s/field", f.TableName)
	_, err := JsonReq("POST", url, &f, rt)
	if err == nil {
		exit_code = rt.Code
		log.Println(rt.Msg)

	} else {
		exit_code = -1000
		log.Println(err)
	}
	os.Exit(exit_code)
}

func QueryField(table_name, fild_name string) {

	flds, _ := getFields(table_name, fild_name)
	fmt.Printf("Name\tType\tLen\tDeci\tNull\tAutoInc\tGener\tDefault\tComment\n")
	fmt.Printf("----------------------------------------------------------------------------\n")

	for _, fld := range flds {
		printField(&fld)
	}
}

func getFields(table_name, field_name string) (fields []Field, err error) {
	fields = make([]Field, 0)
	if table_name != "" {
		rt := new(RTFLD)
		url := fmt.Sprintf("http://localhost:9099/table/%s/field/%s", table_name, field_name)
		_, err = JsonReq("GET", url, nil, rt)
		if err == nil {
			if rt.Code == 0 {
				fields = rt.Rlt
			}
			log.Println(rt.Code)
			log.Println(rt.Msg)
			log.Println(rt.Rlt)
		}
	}
	return
}

func printField(f *Field) {
	fmt.Printf("%s\t%s\t%d\t%d\t%v\t%v\t%v\t%s\t%s\n", f.Name, f.FieldType, f.Length, f.Decimals,
		f.Nullable, f.AutoIncre, f.GeneratedField, f.DefaultV, f.Comment)
}

func setField(f *flag.Flag) {
	switch f.Name {
	case "type", "t":
		fi.FieldType = dataType
	case "defaultValue", "v":
		fi.DefaultV = defaultValue
	case "comment", "c":
		fi.Comment = comment
	case "nullable", "n":
		fi.Nullable = isNullable
	case "autoInreace", "auto":
		fi.AutoIncre = isAutoInreace
	case "generate", "g":
		fi.GeneratedField = isGenerate
	case "lenth", "l":
		fi.Length = uint8(lenth)
	case "decimal", "d":
		fi.Decimals = uint8(decimal)

	}
}

//type Field struct {
//	TableName      string
//	Name           string
//	FieldType      string
//	Length         uint8
//	Decimals       uint8
//	Nullable       bool
//	DefaultV       string
//	AutoIncre      bool
//	GeneratedField bool
//	Comment        string
//	//GoType         string
//	//reserved
//	//flag fieldFlag
//}

//func NewField(name, dataType string, lenth, decimal uint8, nullable bool,
//defaultValue string, autoIncrement, isGeneratedField bool, comment string)
