package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"hash"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type ossClient struct {
	bucket          *oss.Bucket
	accessKeyId     string
	accessKeySecret string
}

func (oc *ossClient) init(ak string, as string) {
	var err error
	var client *oss.Client
	if client, err = oss.New("oss-cn-shenzhen-internal.aliyuncs.com", ak, as); err != nil {
		panic(err)
	}
	if oc.bucket, err = client.Bucket("agb-image"); err != nil {
		panic(err)
	}
	oc.accessKeyId = ak
	oc.accessKeySecret = as
}

func (oc *ossClient) upload(fn string, fb []byte, c chan<- interface{}) {
	defer close(c)

	if err := oc.bucket.PutObject(fn, bytes.NewReader(fb)); err != nil {
		c <- err
	}
}

func (oc *ossClient) delete(fn string) {
	_ = oc.bucket.DeleteObject(fn)
}

func (oc *ossClient) ossSignature(m string, cm string, ct string, e int64, r string) string {
	signStr := m + "\n" + cm + "\n" + ct + "\n" + strconv.FormatInt(e, 10) + "\n" + r
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(oc.accessKeySecret))
	io.WriteString(h, signStr)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func getResource(b string, o string) string {
	r, _ := regexp.Compile("^/+")
	return "/" + b + "/" + r.ReplaceAllString(o, "")
}

func (oc *ossClient) signURL(u string, o string) string {
	expire := time.Now().Unix() + 3*60

	su, _ := url.Parse(u)
	q := su.Query()
	q.Set("OSSAccessKeyId", oc.accessKeyId)
	q.Set("Expires", strconv.FormatInt(expire, 10))
	q.Set("Signature", oc.ossSignature("GET", "", "", expire, getResource("agb-image", o)))
	su.RawQuery = q.Encode()

	return su.String()
}
