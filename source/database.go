package source

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	// postgres as sql driver
	_ "github.com/lib/pq"
)

// ImageInfo describes property of image from database
type ImageInfo struct {
	ID  uuid.UUID
	URL string
}

var db *sql.DB
var rc uint64

func init() {
	var err error
	//panic if cannot connect to pg
	db, err = sql.Open("postgres", fmt.Sprintf("host=127.0.0.1 dbname=%s user=%s password=%s", os.Getenv("PG_DATABASE"), os.Getenv("PG_USERNAME"), os.Getenv("PG_PASSWORD")))
	if err = db.Ping(); err != nil {
		panic(err)
	}
}

// Size returns row count
func Size() uint64 {
	if rc == 0 {
		if err := db.QueryRow("SELECT COUNT(*) FROM gif").Scan(&rc); err != nil {
			// log.Fatal(err)
		}
	}

	return rc
}

// MatchImages returns matched gif url to a specific gif
func MatchImages(h []byte, c chan<- interface{}) {
	rows, err := db.Query(`SELECT id, url FROM gif WHERE id = ANY(SELECT match FROM gif WHERE img_hash = $1)`, h)
	if err != nil {
		c <- errors.New("Query match URL failed")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var img ImageInfo
		rows.Scan(&img.ID, &img.URL)
		c <- img
	}
	if err := rows.Err(); err != nil {
		c <- errors.New("Rows iteration failed in MatchURLs")
	} else {
		// iteration done without error
		c <- true
	}
}

// NthImages returns url of GIF from offset o with length l
func NthImages(o uint64, l int, c chan<- interface{}) {
	rows, err := db.Query(`SELECT id, url FROM gif LIMIT $2 OFFSET $1`, o, l)
	if err != nil {
		c <- errors.New("Query nth URL failed")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var img ImageInfo
		rows.Scan(&img.ID, &img.URL)
		c <- img
	}
	if err := rows.Err(); err != nil {
		c <- errors.New("Rows iteration failed in NthURLs")
	} else {
		// iteration done without error
		c <- true
	}
}

// NewMatch insert a new match to database
func NewMatch(h uuid.UUID, a uuid.UUID, s string, c chan<- interface{}) {
	_, err := db.Exec(`INSERT INTO submit_match VALUES ($1, $2, $3, $4, LOCALTIMESTAMP)`, uuid.New(), h, a, s)
	if err != nil {
		c <- errors.New("Exec failed in NewMatch")
	} else {
		c <- true
	}
}
