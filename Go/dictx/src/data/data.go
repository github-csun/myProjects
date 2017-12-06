package main

import (
	. "dictx"
	. "dictx/client"
	"flag"
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
		action string
	)

	flag.StringVar(&action, "action", "", "New name for the table")
	flag.StringVar(&action, "a", "", "New name for the table")
	flag.Parse()

	switch strings.ToLower(action) {
	case "s", "save":
		Save()
	case "l", "load":
		load()
	}

}

func Save() {
	rt := new(RT)

	_, err := JsonReq("POST", "http://localhost:9099/data", rt, rt)
	if err == nil {
		exit_code = rt.Code
		log.Println(rt.Msg)
	} else {
		exit_code = -1000
		log.Println(err)
	}
	os.Exit(exit_code)

}

func load() {

}
