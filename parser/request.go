package parser

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

// Image represents request image information
type Image struct {
	Data   []byte
	From   int
	Length int
	Format string
}

// Match represents request match information
type Match struct {
	Home      uuid.UUID
	Away      uuid.UUID
	Submitter string
}

// ParseGIF parses the 'gif' request
func ParseGIF(r *http.Request, c chan<- interface{}) {
	var img Image
	var err error
	for {
		if r.ContentLength > (1 << 19) {
			err = errors.New("Body too large")
			break
		}
		if err = r.ParseMultipartForm(1 << 19); err != nil {
			err = errors.New("Body too large")
			break
		}
		if img.From, err = strconv.Atoi(r.FormValue("from")); err != nil {
			err = errors.New("No field from")
			break
		}
		if img.Length, err = strconv.Atoi(r.FormValue("length")); err != nil {
			err = errors.New("No field length")
			break
		}
		f, _, e := r.FormFile("image-file")
		if e != nil {
			err = errors.New("No field file")
			break
		}
		defer f.Close()
		img.Data, err = ioutil.ReadAll(f)
		if err != nil {
			err = errors.New("Load file failed")
			break
		}
		if http.DetectContentType(img.Data) != "image/gif" {
			err = errors.New("Image is not GIF")
			break
		}
		img.Format = "image/gif"
		c <- img
	}
	c <- err
}

// ParseMatch parses the 'match' request
func ParseMatch(r *http.Request, c chan<- interface{}) {
	var m Match
	var err error
	for {
		if err = r.ParseMultipartForm(1 << 19); err != nil {
			err = errors.New("Body too large")
			break
		}
		if m.Home, err = uuid.Parse(r.FormValue("home")); err != nil {
			err = errors.New("No field home")
			break
		}
		if m.Away, err = uuid.Parse(r.FormValue("away")); err != nil {
			err = errors.New("No field away")
			break
		}
		if m.Submitter = r.FormValue("user-identifier"); m.Submitter == "" {
			err = errors.New("No field user-identifier")
			break
		}
		c <- m
	}
	c <- err
}
