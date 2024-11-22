package main

import (
	"bytes"
	"context"
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"strings"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region     string
	bucketName string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
}

func main() {
	flag.Parse()
	var (
		route = "fc output route"
		token = "fc output token"
	)
	if len(bucketName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name required")
	}

	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := oss.NewClient(cfg)

	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	red := color.RGBA{255, 0, 0, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{red}, image.Point{}, draw.Src)

	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		log.Fatalf("failed to put access point policy for object process %v", err)
	}

	request := &oss.WriteGetObjectResponseRequest{
		RequestRoute: oss.Ptr(route),
		RequestToken: oss.Ptr(token),
		FwdStatus:    oss.Ptr("200"),
		Body:         strings.NewReader(string(buf.Bytes())),
	}
	result, err := client.WriteGetObjectResponse(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to write get object response %v", err)
	}

	log.Printf("write get object response result:%#v\n", result)
}
