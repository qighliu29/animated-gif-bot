package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type imageReq struct {
	Data   []byte
	From   int
	Length int
	Format string
}

type matchReq struct {
	Home      uuid.UUID
	Away      uuid.UUID
	Submitter string
}

func parseGIFReq(r *http.Request, c chan<- interface{}) {
	var ir imageReq
	var err error

	defer close(c)

	for {
		if r.ContentLength > (1 << 19) {
			err = errors.New("Body too large")
			break
		}
		if err = r.ParseMultipartForm(1 << 19); err != nil {
			err = errors.New("Body too large")
			break
		}
		if ir.From, err = strconv.Atoi(r.FormValue("from")); err != nil {
			err = errors.New("No field from")
			break
		}
		if ir.Length, err = strconv.Atoi(r.FormValue("length")); err != nil {
			err = errors.New("No field length")
			break
		}
		f, _, e := r.FormFile("image-file")
		if e != nil {
			err = errors.New("No field file")
			break
		}
		defer f.Close()
		ir.Data, err = ioutil.ReadAll(f)
		if err != nil {
			err = errors.New("Load file failed")
			break
		}
		if http.DetectContentType(ir.Data) != "image/gif" {
			err = errors.New("Image is not GIF")
			break
		}
		ir.Format = "image/gif"
		c <- ir
		return
	}
	c <- err
}

func parseMatchReq(r *http.Request, c chan<- interface{}) {
	var mr matchReq

	defer close(c)

	var err error
	for {
		if err = r.ParseMultipartForm(1 << 19); err != nil {
			err = errors.New("Body too large")
			break
		}
		if mr.Home, err = uuid.Parse(r.FormValue("home")); err != nil {
			err = errors.New("No field home")
			break
		}
		if mr.Away, err = uuid.Parse(r.FormValue("away")); err != nil {
			err = errors.New("No field away")
			break
		}
		if mr.Submitter = r.FormValue("user-identifier"); mr.Submitter == "" {
			err = errors.New("No field user-identifier")
			break
		}
		c <- mr
		return
	}
	c <- err
}
