package routex

import (
	//. "dictx"
	// "dictx/client"
	//	"encoding/json"
	// "fmt"
	"log"
	// "net/http"
	//	"os"
	//"math"
	//l "logx"
	//"regexp"
	//"strings"
)
import ()

type Context struct {
	PathParam map[string]string
	regRoute  *routeRec
}

func NewContext() (context *Context) {
	context = new(Context)
	context.PathParam = make(map[string]string)

	return
}

func (this *Context) setCtx(ctx *Context) {
	this.PathParam = ctx.PathParam
	this.regRoute = ctx.regRoute
	log.Printf("Set path param: %v", this.PathParam)
	return
}

func (this Context) GetPathParam(varName string) (value string) {
	value, ok := this.PathParam[varName]
	if ok == false {
		value = ""
	}
	return
}
