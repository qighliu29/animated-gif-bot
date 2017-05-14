package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	b2b "github.com/minio/blake2b-simd"
	ps "github.com/qighliu29/animated-gif-bot/parser"
	src "github.com/qighliu29/animated-gif-bot/source"
)

type imgIDURL struct {
	ID  string
	URL string
}

type resGIF struct {
	ID         string
	MatchArray []imgIDURL
	MatchCount int
}

type resMatch struct {
	Message string
}

func imageInfo2IDURL(s []src.ImageInfo) []imgIDURL {
	res := make([]imgIDURL, 0, len(s))
	for _, i := range s {
		res = append(res, imgIDURL{ID: i.ID.String(), URL: i.URL})
	}

	return res
}

func gifHandler(w http.ResponseWriter, r *http.Request) {
	rc := make(chan interface{})
	var img ps.Image
	var cur src.ImageInfo
	var vc int
	m := make([]src.ImageInfo, 0, 10)

	go ps.ParseGIF(r, rc)
	switch t := <-rc; t.(type) {
	case ps.Image:
		img = t.(ps.Image)
	case error:
		badRequest(w, r)
		return
	}

	h := b2b.Sum256(img.Data)
	// success("%s\n", hex.EncodeToString(h[:]))
	go src.MatchImages(h[:], rc)
	for t := range rc {
		switch t.(type) {
		case src.ImageInfo:
			m = append(m, t.(src.ImageInfo))
		case bool:
			break
		case error:
			internalError(w, r)
			return
		}
	}

	if len(m) == 0 {
		// gif does not exist, upload to OSS & set callback
		// wait it complete
	} else {
		cur = m[1]
		m = m[1:]
	}

	vf, ve := intersec(0, len(m), img.From, img.From+img.Length)
	m = m[vf:ve]
	vc = len(m)
	if vc < img.Length {
		mr := img.Length - vc
		s := rand.NewSource(time.Now().UnixNano())
		ran := rand.New(s)
		n := ran.Int63n(int64(src.Size()/uint64(mr-1))) * int64(mr)
		go src.NthImages(uint64(n), mr, rc)
		for t := range rc {
			switch t.(type) {
			case src.ImageInfo:
				m = append(m, t.(src.ImageInfo))
			case bool:
				break
			case error:
				internalError(w, r)
				return
			}
		}
	}
	rd, err := json.Marshal(resGIF{ID: cur.ID.String(), MatchArray: imageInfo2IDURL(m), MatchCount: vc})
	if err != nil {
		internalError(w, r)
	} else {
		w.Write(rd)
	}
}

func matchHandler(w http.ResponseWriter, r *http.Request) {
	c := make(chan interface{})
	var m ps.Match

	go ps.ParseMatch(r, c)
	switch t := <-c; t.(type) {
	case ps.Match:
		m = t.(ps.Match)
	case error:
		badRequest(w, r)
		return
	}

	go src.NewMatch(m.Home, m.Away, m.Submitter, c)
	switch t := <-c; t.(type) {
	case error:
		internalError(w, r)
		return
	}
	rd, _ := json.Marshal(resMatch{Message: "OK"})
	w.Write(rd)
}

func main() {
	// launch http server
	http.HandleFunc("/gif", handleWithMethod("POST", gifHandler))
	http.HandleFunc("/match", handleWithMethod("POST", matchHandler))
	http.ListenAndServe(":8080", nil)
}
