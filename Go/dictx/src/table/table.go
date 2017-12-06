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
)

func main() {

	var (
		//New, Update, Query, Delete bool
		newName, action string
	)
	flag.StringVar(&newName, "name", "", "New name for the table")
	flag.StringVar(&newName, "n", "", "New name for the table")
	flag.StringVar(&action, "action", "", "New name for the table")
	flag.StringVar(&action, "a", "", "New name for the table")
	flag.Parse()
	args := flag.Args()
	for i, arg := range args {
		args[i] = strings.ToUpper(arg)
	}

	switch strings.ToLower(action) {
	case "c", "create":
		NewTable(args)
	case "q", "query":
		QueryTable(args)
	case "d", "delete":

		DeleteTable(args)

	case "u", "update":
		if newName != "" {
			PutTable(args[0], newName)
		}
	}

}

func NewTable(args []string) {

	rt := new(RT)
	for _, arg := range args {

		var n = ItemName{Name: arg}
		_, err := JsonReq("POST", "http://localhost:9099/table", &n, rt)
		if err == nil {
			exit_code = rt.Code
			log.Println(rt.Msg)
		} else {
			exit_code = -1000
			log.Println(err)
		}
		os.Exit(exit_code)
	}
}

func QueryTable(args []string) {

	rt := new(RTQT)

	if len(args) == 0 {
		_, err := JsonReq(GET, "http://localhost:9099/table?tblName=", nil, rt)
		if err == nil {
			exit_code = rt.Code
			fmt.Println("No.\tTable Name")
			fmt.Println("__________________________________________")
			for i, tbl := range rt.Rlt {

				fmt.Printf("%d\t%v\n", i, tbl.TblName)
			}

		} else {
			exit_code = -1000

		}
		os.Exit(exit_code)
	} else {
		for _, arg := range args {
			url := fmt.Sprintf("http://localhost:9099/table?tblName=%s", arg)
			_, err := JsonReq(GET, url, nil, rt)
			if err == nil {
				exit_code = rt.Code
				fmt.Println("No.\tTable Name")
				fmt.Println("__________________________________________")
				for i, tbl := range rt.Rlt {
					fmt.Printf("%d\t%v\n", i, tbl.TblName)
				}

			} else {
				exit_code = -1000
				log.Println(err)
				os.Exit(exit_code)
			}
		}
	}

	//return
}

func DeleteTable(args []string) {
	if len(args) == 0 {
		exit_code = -2000
		log.Fatal("Missed table name of deletion.")
		return
	}
	var rt = new(RT)
	for _, arg := range args {
		var n = ItemName{Name: arg}
		_, err := JsonReq(DELETE, "http://localhost:9099/table", &n, rt)
		if err == nil {
			exit_code = rt.Code
			log.Println(rt.Msg)
		} else {
			exit_code = -1000
			log.Println(err)
		}
		os.Exit(exit_code)
	}

	return
}

func PutTable(origin, newName string) {
	newName = strings.ToUpper(newName)
	tblreq := TableUpdateRequest{OriginName: origin, NewName: newName}
	var rt = new(RT)

	_, err := JsonReq(PUT, "http://localhost:9099/table", &tblreq, rt)
	if err == nil {
		exit_code = rt.Code
		log.Println(rt.Msg)
	} else {
		exit_code = -1000
		log.Println(err)
	}
	os.Exit(exit_code)
	return
}
