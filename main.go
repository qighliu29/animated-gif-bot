package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"

	_ "github.com/lib/pq"
)

// type ReqContent struct {
// 	id   string
// 	from int8
// 	size int8
// }

type resContent struct {
	ID   []string
	Size int8
}

func handler(w http.ResponseWriter, r *http.Request) {
	validPath := regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
	} else {
		resData, _ := json.Marshal(resContent{ID: []string{"abc", "xyz"}, Size: 1})
		w.Write(resData)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)

	db, err := sql.Open("postgres", "dbname=postgres user=postgres password=lq0729 host=172.19.10.1")
	if err != nil {
		return
	}
	db.Query("")
}
