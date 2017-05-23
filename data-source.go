package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type pgImageRepo struct {
	db    *sql.DB
	count uint64
}

type imageRow struct {
	ID  uuid.UUID
	URL string
}

func (repo *pgImageRepo) connect(h string, d string, u string, p string) {
	var err error
	//panic if cannot connect to pg
	repo.db, err = sql.Open("postgres", fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", h, d, u, p))
	if err = repo.db.Ping(); err != nil {
		panic(err)
	}
}

func (repo *pgImageRepo) size() uint64 {
	if repo.count == 0 {
		if err := repo.db.QueryRow("SELECT COUNT(*) FROM gif").Scan(&repo.count); err != nil {
			// log.Fatal(err)
		}
	}

	return repo.count
}

func (repo *pgImageRepo) matchImages(h []byte, c chan<- interface{}) {
	defer close(c)

	// 2 solutions combine ANY with ARRAY RETURNED by SUBSELECT: https://dba.stackexchange.com/questions/90460/how-to-do-an-anyselect-query-in-postgresql
	rows, err := repo.db.Query(`SELECT id, url FROM gif WHERE array[id] <@ (SELECT match FROM gif WHERE img_hash = $1)`, h)
	if err != nil {
		c <- errors.New("Query match URL failed")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var row imageRow
		rows.Scan(&row.ID, &row.URL)
		c <- row
	}
	if err := rows.Err(); err != nil {
		c <- errors.New("Rows iteration failed in MatchURLs")
	}
}

func (repo *pgImageRepo) newImage(id uuid.UUID, sz int, hb []byte, c chan<- interface{}) {
	defer close(c)

	_, err := repo.db.Exec(`INSERT INTO gif VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, LOCALTIMESTAMP, LOCALTIMESTAMP)`,
		id,
		"http://agb-image.oss-cn-shenzhen.aliyuncs.com/"+id.String()+".gif",
		pq.Array([]string{}),
		pq.Array([]int{}),
		pq.Array([]uuid.UUID{id}),
		sz,
		"gif",
		hb,
		"app-server")
	if err != nil {
		c <- errors.New("Exec failed in newImage: " + err.Error())
	}
}

func (repo *pgImageRepo) nthImages(o uint64, l int, c chan<- interface{}) {
	defer close(c)

	rows, err := repo.db.Query(`SELECT id, url FROM gif LIMIT $2 OFFSET $1`, o, l)
	if err != nil {
		c <- errors.New("Query nth URL failed")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var row imageRow
		rows.Scan(&row.ID, &row.URL)
		c <- row
	}
	if err := rows.Err(); err != nil {
		c <- errors.New("Rows iteration failed in NthURLs")
	}
}

func (repo *pgImageRepo) newMatch(h uuid.UUID, a uuid.UUID, s string, c chan<- interface{}) {
	defer close(c)

	_, err := repo.db.Exec(`INSERT INTO submit_match VALUES ($1, $2, $3, $4, LOCALTIMESTAMP)`, uuid.New(), h, a, s)
	if err != nil {
		c <- errors.New("Exec failed in NewMatch")
	}
}
