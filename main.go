package main

import (
	"database/sql"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"

	"fmt"

	"encoding/hex"

	"strconv"

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
	db       *sql.DB
	gifCount uint64
}

var dataProvider pgManager

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

func (pg *pgManager) size() uint64 {
	if pg.gifCount == 0 {
		pg.queryGifCount()
	}

	return pg.gifCount
}

func (pg *pgManager) queryGifCount() {
	err := pg.db.QueryRow(`SELECT COUNT(*) FROM $1`, "gif").Scan(&pg.gifCount)
	if err == nil {
		return
	}
	//pq error
	if err, ok := err.(*pq.Error); ok {
		error("pq error:", err.Code.Name())
		return
	}
	error(err.Error())
}

func (pg *pgManager) matchURL(key []byte) (URLs []string, miss bool) {
	rows, err := pg.db.Query(`SELECT url FROM gif WHERE id = ANY (SELECT match FROM gif WHERE id = $1)`, key)
	if err != nil {
		//print to clear to error message of miss of item
		error(err.Error())
		return URLs, true
	}
	defer rows.Close()

	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			error(err.Error())
		} else {
			URLs = append(URLs, url)
		}
	}

	if err := rows.Err(); err != nil {
		error(err.Error())
	}
	return
}

func (pg *pgManager) newGIF(hash []byte /*url string, */) {

}

func (pg *pgManager) nthURL(nth uint64) (URL string) {
	err := pg.db.QueryRow(`SELECT url FROM gif LIMIT 1 OFFSET $1`, nth).Scan(&URL)
	if err == nil {
		return
	}
	//pq error
	if err, ok := err.(*pq.Error); ok {
		error("pq error:", err.Code.Name())
		return
	}
	error(err.Error())
	return
}

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

func intersection(l1, r1, l2, r2 int) (l, r int) {
	l = mustInRange(l1, r1, l2)
	r = mustInRange(l, r1, r2)
	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	//should be POST
	validPath := regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
	} else {
		//extract image data
		imageFile, imageFileHeader, err := r.FormFile("image-data")
		if err != nil {
			error(err.Error())
			return
		}
		success(imageFileHeader.Filename)
		//extract 'from' & 'size' field
		from, err := strconv.Atoi(r.FormValue("from"))
		if err != nil {
			error(err.Error())
			return
		}
		size, err := strconv.Atoi(r.FormValue("size"))
		if err != nil {
			error(err.Error())
			return
		}
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
		URLs, miss := dataProvider.matchURL(imageHash[:])
		if miss {
			//upload file to OSS & set callback
		}

		vFrom, vEnd := intersection(0, len(URLs), from, from+size)
		URLs = URLs[vFrom:vEnd]

		//less than requested
		if len(URLs) < size {
			ranIdx := make([]int64, 0, size-len(URLs))
			//only use low 63bit
			sizen := int64((dataProvider.size()) & ((uint64(1) << 63) - 1))
			s1 := rand.NewSource(time.Now().UnixNano())
			r1 := rand.New(s1)
			for i := 0; i < size-len(URLs); i++ {
				ranIdx = append(ranIdx, r1.Int63n(sizen))
			}
			//read URL from database
			for _, nth := range ranIdx {
				URLs = append(URLs, dataProvider.nthURL(uint64(nth)))
			}
		}

		resData, _ := json.Marshal(resContent{ID: []string{"abc", "xyz"}, Size: 1})
		w.Write(resData)
	}
}

func main() {
	//init database connection
	dataProvider.connect()
	//launch http server
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
