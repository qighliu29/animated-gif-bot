package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"regexp"

	"fmt"

	"encoding/hex"

	ct "github.com/daviddengcn/go-colortext"
	pq "github.com/lib/pq"
	b2b "github.com/minio/blake2b-simd"
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

type pgManager struct {
	db   *sql.DB
	size uint64
}

func error(format string, args ...interface{}) {
	// fmt.Printf(chalk.Red.Color(format), args...)
	ct.Foreground(ct.Red, false)
	fmt.Printf(format, args...)
	ct.ResetColor()
}

func success(format string, args ...interface{}) {
	// fmt.Printf(chalk.Green.Color(format), args...)
	ct.Foreground(ct.Green, false)
	fmt.Printf(format, args...)
	ct.ResetColor()
}

func (pg *pgManager) connect() {
	//get database config
	pgUsername := os.Getenv("PG_USERNAME")
	pgPassword := os.Getenv("PG_PASSWORD")
	pgDatabase := os.Getenv("PG_DATABASE")

	db, err := sql.Open("postgres", fmt.Sprintf("dbname=%s user=%s password=%s", pgDatabase, pgUsername, pgPassword))
	if err != nil {
		error(err.Error())
		return
	}
	pg.db = db
}

func (pg *pgManager) querySize() {
	err := pg.db.QueryRow(`SELECT COUNT(*) FROM $1`, "image").Scan(&pg.size)
	if err == nil {
		return
	}
	if err, ok := err.(*pq.Error); ok {
		error("pq error:", err.Code.Name())
		return
	}
	error(err.Error())

}

func (pg *pgManager) row(key []byte) {
	var match string
	pg.db.QueryRow(`SELECT match FROM $1`, "image").Scan(&match)
	//cond & union query return requested num of URLs
}

func handler(w http.ResponseWriter, r *http.Request) {
	//should be POST
	validPath := regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
	} else {
		//read JSON body
		imageFile, imageFileHeader, err := r.FormFile("image-data")
		if err != nil {
			error(err.Error())
			return
		}
		success(imageFileHeader.Filename)
		//calc file hash
		var imageData []byte
		imageBytesLength, err := imageFile.Read(imageData)
		if err != nil {
			error(err.Error())
			return
		}
		success("%d\n", imageBytesLength)
		imageHash := b2b.Sum256(imageData)
		success("%s\n", hex.EncodeToString(imageHash[:]))
		//select from database
		resData, _ := json.Marshal(resContent{ID: []string{"abc", "xyz"}, Size: 1})
		w.Write(resData)
	}
}

func main() {
	//launch http server
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
