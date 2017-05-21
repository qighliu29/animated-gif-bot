package main

import (
	"encoding/json"
	"net/http"
)

var acao string

func setCORSHeader(c string) {
	acao = c
}

func resJSON(w http.ResponseWriter, o interface{}) {
	var bs []byte
	var err error
	if bs, err = json.Marshal(o); err != nil {
		failed("[resJSON] Marshal failed: %s", err.Error())
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", acao)
	if _, err = w.Write(bs); err != nil {
		failed("[resJSON] Write failed: %s", err.Error())
	}
}

func resMessage(w http.ResponseWriter, msg string) {
	resJSON(w, struct{ Message string }{Message: msg})
}

func resMessageData(w http.ResponseWriter, msg string, dt interface{}) {
	resJSON(w, struct {
		Message string
		Data    interface{}
	}{Message: msg, Data: dt})
}

func resData(w http.ResponseWriter, dt interface{}) {
	resMessageData(w, "OK", dt)
}

func resUnsupportedMethod(w http.ResponseWriter) {
	resMessage(w, "Unsupported Method")
}

func resBadRequest(w http.ResponseWriter) {
	resMessage(w, "Bad Request")
}

func resInternalError(w http.ResponseWriter) {
	resMessage(w, "Internal Error")
}

func resOK(w http.ResponseWriter) {
	resMessage(w, "OK")
}

func handleWithMethod(method string, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			resUnsupportedMethod(w)
			return
		}
		fn(w, r)
	}
}
