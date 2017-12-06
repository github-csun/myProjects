package routex

import (
	c "dictx/client"
	"log"
	l "logx"
	"net/http"
	"testing"
)

var path string = "/test/table"

type testCtrl struct {
}

func (t *testCtrl) Get(w http.ResponseWriter, r *http.Request) {
	logP(path, r)
}

func (t *testCtrl) Post(w http.ResponseWriter, r *http.Request) {

}
func (t *testCtrl) Put(w http.ResponseWriter, r *http.Request) {

}
func (t *testCtrl) Delete(w http.ResponseWriter, r *http.Request) {
}

func logP(path string, r *http.Request) {
	log.Printf("PathDef: \"%s.\"", path)
	log.Printf("PathReq: \"%s.\"", r.URL.Path)
}

func TestRoute(t *testing.T) {
	l.Lg("Started!")
	RegRoute(path+"/:xx", new(testCtrl))
	go Serve(":9989")
	c.JsonReq(c.GET, "http://localhost:9989/test/table/u", nil, nil)
}
