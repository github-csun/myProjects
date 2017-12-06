package main

import (
	. "dictx"
	. "dictx/client"
	"fmt"
	"log"
	"os"
)

var (
	table_name string
)

func main() {
	if len(os.Args) != 2 {
		log.Printf("Invalide program arguments.\n")
		os.Exit(1)
	} else {
		table_name = os.Args[1]
		url := fmt.Sprintf("http://localhost:9099/table/%s/go", table_name)
		o := new(RTSTR)
		_, err := JsonReq(GET, url, nil, o)
		if err == nil {
			println(o.Rlt)
		} else {
			log.Println(err)
		}
		return
	}
}
