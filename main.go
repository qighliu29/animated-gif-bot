package main

import (
	"math/rand"
	"net/http"
	"time"

	"flag"

	b2b "github.com/minio/blake2b-simd"
)

var dbhost = flag.String("h", "127.0.0.1", "the host name of the machine on which the database server is running")
var dbname = flag.String("d", "agb", "the name of the image database")
var dbuser = flag.String("u", "agb", "the username which will be used to connect to the image database")
var dbpwd = flag.String("p", "agbpassword", "the password which will be used to connect to the image database")
var cors = flag.String("c", "*", "the \"Access-Control-Allow-Origin\" field will be set in response headers")

var repo pgImageRepo

func init() {
	repo.connect(*dbhost, *dbname, *dbuser, *dbpwd)
}

func imageInfo2IDURL(s []imageRow) []interface{} {
	res := make([]interface{}, 0, len(s))
	for _, i := range s {
		res = append(res, struct {
			ID  string
			URL string
		}{ID: i.ID.String(), URL: i.URL})
	}

	return res
}

func randomImageN(n int) ([]imageRow, error) {
	m := make([]imageRow, 0, n)
	// success("request %d images\n", n)
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	o := r.Int63n(int64(repo.size() - uint64(n) + 1))
	c := make(chan interface{})
	go repo.nthImages(uint64(o), n, c)
	if err := readChanUntilClose(c, func(arg interface{}) {
		m = append(m, arg.(imageRow))
	}); err != nil {
		return nil, err
	}
	// success("marigin %d images\n", len(m))
	return m, nil
}

func gifHandler(w http.ResponseWriter, r *http.Request) {
	var req imageReq
	var cur imageRow
	m := make([]imageRow, 0, 10)

	c := make(chan interface{})
	go parseGIFReq(r, c)
	if err := readChanUntilClose(c, func(arg interface{}) {
		req = arg.(imageReq)
	}); err != nil {
		resBadRequest(w)
		return
	}

	h := b2b.Sum512(req.Data)
	// success("%s\n", hex.EncodeToString(h[:]))
	c = make(chan interface{})
	go repo.matchImages(h[:], c)
	if err := readChanUntilClose(c, func(arg interface{}) {
		m = append(m, arg.(imageRow))
	}); err != nil {
		resInternalError(w)
		return
	}

	if len(m) == 0 {
		// gif does not exist, upload to OSS & set callback
		// wait it complete
	} else {
		cur = m[0]
		m = m[1:]
	}

	vf, ve := intersec(0, len(m), req.From, req.From+req.Length)
	m = m[vf:ve]
	mc := len(m)
	if mc < req.Length {
		if ri, err := randomImageN(req.Length - mc); err == nil {
			m = append(m, ri...)
		} else {
			resInternalError(w)
			return
		}
	}

	resData(w, struct {
		ID         string
		MatchArray []interface{}
		MatchCount int
	}{ID: cur.ID.String(), MatchArray: imageInfo2IDURL(m), MatchCount: mc})
}

func matchHandler(w http.ResponseWriter, r *http.Request) {
	var req matchReq

	c := make(chan interface{})
	go parseMatchReq(r, c)
	if err := readChanUntilClose(c, func(arg interface{}) {
		req = arg.(matchReq)
	}); err != nil {
		resBadRequest(w)
		return
	}

	c = make(chan interface{})
	go repo.newMatch(req.Home, req.Away, req.Submitter, c)
	if err := readChanUntilClose(c, func(arg interface{}) {}); err != nil {
		resInternalError(w)
		return
	}

	resOK(w)
}

func main() {
	flag.Parse()
	setCORSHeader(*cors)

	http.HandleFunc("/gif", handleWithMethod("POST", gifHandler))
	http.HandleFunc("/match", handleWithMethod("POST", matchHandler))
	http.ListenAndServe(":8080", nil)
}
