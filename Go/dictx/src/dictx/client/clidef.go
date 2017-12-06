package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	PATCH  = "PATCH"
	DELETE = "DELETE"
)

type TableUpdateRequest struct {
	OriginName string
	NewName    string
}

func JsonReq(method, url string, input, output interface{}) (
	resp *http.Response, err error) {

	var (
		buf     = new(bytes.Buffer)
		client  = &http.Client{}
		req     *http.Request
		decoder *json.Decoder
	)

	switch method {
	case PUT, POST, DELETE, GET:
		err = nil
	default:
		err = errors.New("Invalid method: \"" + method + "\".")
		goto RETURN
	}

	if input != nil {

		encoder := json.NewEncoder(buf)
		err = encoder.Encode(input)
		if err != nil {
			goto RETURN
		}
	}
	req, err = http.NewRequest(method, url, buf)
	if err != nil {
		goto RETURN
	}
	resp, err = client.Do(req)
	if err != nil {
		goto RETURN
	}
	if output != nil {
		decoder = json.NewDecoder(resp.Body)
		err = decoder.Decode(output)
	}
RETURN:
	if err != nil {
		log.Println(err)
	}
	return
}

func JsonEncode(w io.Writer, value interface{}) (w1 io.Writer, err error) {
	if w == nil {
		w = new(bytes.Buffer)
	}
	err = json.NewEncoder(w).Encode(value)
	w1 = w
	return
}

func JsonDecode(r io.Reader, value interface{}) (err error) {
	err = json.NewDecoder(r).Decode(value)
	return
}
