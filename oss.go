package main

import (
	"bytes"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type ossClient struct {
	bucket *oss.Bucket
}

func (oc *ossClient) init(ak string, as string) {
	var err error
	var client *oss.Client
	if client, err = oss.New("agb-image.oss-cn-shenzhen-internal.aliyuncs.com", ak, as); err != nil {
		panic(err)
	}
	if oc.bucket, err = client.Bucket("agb-image"); err != nil {
		panic(err)
	}
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
