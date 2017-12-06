package routex

import (
	//. "dictx"
	. "dictx/client"
	//	"encoding/json"
	"fmt"
	"log"
	"net/http"
	//	"os"
	//"math"
	l "logx"
	"regexp"
	"strings"
)

var (
	reg      *regexp.Regexp
	routeMap []*routeRec
)

func init() {

	routeMap = make([]*routeRec, 0, 64)
	var err error
	reg, err = regexp.Compile("\\w+")
	if err != nil {
		log.Fatal(err)
	}
}

func RouteMethod(w http.ResponseWriter, r *http.Request, ctrl Controller) {
	log.Printf("Method: %s, Path: %s Url:\n", r.Method, r.URL.Path, r.URL.RawQuery)
	switch r.Method {
	case GET:
		ctrl.Get(w, r)
	case POST:
		ctrl.Post(w, r)
	case DELETE:
		ctrl.Delete(w, r)
	case PUT:
		ctrl.Put(w, r)
	}
}

//type Ctrl struct{
//  router map[string]HandlerFunc
//}

type HttpCtrl interface {
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type Controller interface {
	HttpCtrl
	New() (obj Controller)
	setCtx(ctx *Context)
}

type routeRec struct {
	path           string
	pathSegment    []string
	ctrl           Controller
	varPositionMap map[int]string
	varNameMap     map[string]int
}

func newRouteRec(path string, ctrl Controller) (rec *routeRec) {
	this := new(routeRec)
	this.path = path
	this.ctrl = ctrl
	this.varPositionMap = make(map[int]string)
	this.varNameMap = make(map[string]int)
	this.pathSegment = strings.Split(this.path, "/")
	for i, seg := range this.pathSegment {
		if strings.HasPrefix(seg, ":") {
			varName := strings.TrimPrefix(seg, ":")
			matched := reg.MatchString(varName)

			if matched {
				this.varPositionMap[i] = varName
				this.varNameMap[varName] = i
				this.pathSegment[i] = "/"
			} else {
				msg := fmt.Sprintf("Invalid var name: \"%s\" in path: \"%s\"",
					varName, this.path)
				log.Fatal(msg)
			}

		}
	}
	rec = this
	return
}

func (this *routeRec) matchPath(pathSegs []string) (matched bool) {
	matched = true
	segs := pathSegs

	if len(this.pathSegment)-len(segs) > 1 ||
		len(this.pathSegment)-len(segs) < 0 {
		matched = false
		return
	}
	for i, seg := range segs {
		if this.pathSegment[i] == "/" {
			continue
		}
		if this.pathSegment[i] != seg {
			matched = false
			return
		}
	}
	return
}

func RegRoute(path string, ctrl Controller) {
	rec := newRouteRec(path, ctrl)
	routeMap = append(routeMap, rec)
}

type RouteX struct {
}

func (x *RouteX) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	matched := false
	pathSegs := strings.Split(r.URL.Path, "/")
	l.Tr("Request: %s", r.URL)

	if r.URL.Path == "/" || r.URL.Path == "/index.html" {
		http.ServeFile(w, r, "index.html")
	} else {

		for _, item := range routeMap {
			if item.matchPath(pathSegs) {
				matched = true
				ctx := NewContext()
				ctx.regRoute = item
				for varName, ind := range item.varNameMap {
					if ind >= len(pathSegs) {
						ctx.PathParam[varName] = ""
					} else {
						ctx.PathParam[varName] = strings.TrimSpace(pathSegs[ind])
					}
				}
				ctrl := item.ctrl.New()
				ctrl.setCtx(ctx)
				log.Printf("Matched RouteRec: %v", item)
				log.Printf("Matched new ctrl: %v", ctrl)
				RouteMethod(w, r, ctrl)
				break
			}
		}
		if matched == false {
			http.ServeFile(w, r, "index.html")
		}
	}
}

func Serve(port string) {
	err := http.ListenAndServe(port, new(RouteX))
	if err != nil {
		l.Lg("Error while start listen, %s", err)
	}
}
