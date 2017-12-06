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
	table_name  string
	index_name  string
	fields      string
	field_names []string
	isUnique    bool
	action      string
	args        []string
)

func main() {
	flag.StringVar(&fields, "f", "", "Field name list splited by ','.")
	flag.StringVar(&fields, "field", "", "Field name list splited by ','.")
	flag.StringVar(&index_name, "index", "", "Index name")
	flag.StringVar(&index_name, "i", "", "Index name")
	flag.StringVar(&action, "action", "", "Index name")
	flag.StringVar(&action, "a", "", "Index name")
	flag.BoolVar(&isUnique, "unique", true, "Is is an unique index")
	flag.BoolVar(&isUnique, "u", true, "Is is an unique index")

	flag.Parse()
	args = flag.Args()
	if len(args) == 0 {
		log.Printf("Table name was missed.")
		os.Exit(1)
	}
	if len(args) > 1 {
		log.Printf("Too many table names.")
		os.Exit(1)
	}
	action = strings.ToLower(action)
	table_name = args[0]
	if index_name == "" && action == "c" {
		log.Printf("Index name was missed.")
		os.Exit(1)
	}

	switch strings.ToLower(action) {
	case "c":
		field_names = make([]string, 0)
		tmp := strings.Split(fields, ",")
		for _, fld := range tmp {
			fld = strings.TrimSpace(fld)
			if fld != "" {
				field_names = append(field_names, fld)
			}
		}
		defIndex(table_name, index_name, isUnique, field_names)
	case "q":
		queryIndexes(table_name, index_name)
	case "d":
		DeleteIndex(table_name, index_name)
	}

}

func defIndex(tableName, indexName string, isUnique bool,
	fieldNames []string) {
	url := fmt.Sprintf("http://localhost:9099/table/%s/index/%s",
		tableName, indexName)
	param := &IndexParam{
		IsUnique:   isUnique,
		FieldNames: fieldNames,
	}
	output := new(R)
	_, err := JsonReq(POST, url, param, output)
	if err == nil {
		log.Println(output.Msg)

	} else {
		log.Println(err)
	}
}

func queryIndexes(tableName, indexName string) {
	_, output, err := GetIndexes(tableName, indexName)
	if err == nil {
		if output.Code == 0 {
			printIndex(output.Rlt)
		} else {
			output.PMsg()
		}
	}

}

func printIndex(indexes []IndexParam) {
	println("No.\tNAME\tUNIQUE\tFIELD(s)")
	println("----------------------------------------------------------------")
	for i, ind := range indexes {
		f := strings.Join(ind.FieldNames, ",")
		fmt.Printf("%d\t%s\t%v\t%s\n", i, ind.Name, ind.IsUnique, f)
	}
}

func GetIndexes(tableName, indexName string) (indexes []IndexParam, output *RTQI, err error) {
	url := fmt.Sprintf("http://localhost:9099/table/%s/index/%s",
		tableName, indexName)
	output = new(RTQI)
	_, err = JsonReq(GET, url, nil, output)
	if err == nil {
		log.Println(output.Msg)
		indexes = output.Rlt
	} else {
		log.Println(err)
	}
	return
}

func DeleteIndex(tableName, indexName string) {
	url := fmt.Sprintf("http://localhost:9099/table/%s/index/%s",
		tableName, indexName)
	output := new(R)
	_, err := JsonReq(DELETE, url, nil, output)
	if err == nil {
		output.PMsg()
	} else {
		log.Println(err)
	}
	return
}
