package main

import (
	"fmt"
	"net/http"

	ct "github.com/daviddengcn/go-colortext"
)

func mustInRange(l, r, num int) int {
	switch {
	case num < l:
		return l
	case num > r:
		return r
	default:
		return num
	}
}

func intersec(l1, r1, l2, r2 int) (l, r int) {
	l = mustInRange(l1, r1, l2)
	r = mustInRange(l, r1, r2)
	return
}

func unsupportedMethod(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("UNSUPPORTED METHOD"))
}

func badRequest(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("BAD REQUEST"))
}

func internalError(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("INTERNAL ERROR"))
}

func handleWithMethod(method string, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			unsupportedMethod(w, r)
			return
		}
		fn(w, r)
	}
}

func failed(format string, args ...interface{}) {
	ct.Foreground(ct.Red, false)
	fmt.Printf(format, args...)
	ct.ResetColor()
}

func success(format string, args ...interface{}) {
	ct.Foreground(ct.Green, false)
	fmt.Printf(format, args...)
	ct.ResetColor()
}
