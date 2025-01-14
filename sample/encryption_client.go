package main

import (
	"bufio"
	"context"
	"flag"
	"io"
	"log"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/crypto"
)

var (
	rsaPublicKey string = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCokfiAVXXf5ImFzKDw+XO/UByW
6mse2QsIgz3ZwBtMNu59fR5zttSx+8fB7vR4CN3bTztrP9A6bjoN0FFnhlQ3vNJC
5MFO1PByrE/MNd5AAfSVba93I6sx8NSk5MzUCA4NJzAUqYOEWGtGBcom6kEF6MmR
1EKib1Id8hpooY5xaQIDAQAB
-----END PUBLIC KEY-----`

	rsaPrivateKey string = `-----BEGIN PRIVATE KEY-----
MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBAKiR+IBVdd/kiYXM
oPD5c79QHJbqax7ZCwiDPdnAG0w27n19HnO21LH7x8Hu9HgI3dtPO2s/0DpuOg3Q
UWeGVDe80kLkwU7U8HKsT8w13kAB9JVtr3cjqzHw1KTkzNQIDg0nMBSpg4RYa0YF
yibqQQXoyZHUQqJvUh3yGmihjnFpAgMBAAECgYA49RmCQ14QyKevDfVTdvYlLmx6
kbqgMbYIqk+7w611kxoCTMR9VMmJWgmk/Zic9mIAOEVbd7RkCdqT0E+xKzJJFpI2
ZHjrlwb21uqlcUqH1Gn+wI+jgmrafrnKih0kGucavr/GFi81rXixDrGON9KBE0FJ
cPVdc0XiQAvCBnIIAQJBANXu3htPH0VsSznfqcDE+w8zpoAJdo6S/p30tcjsDQnx
l/jYV4FXpErSrtAbmI013VYkdJcghNSLNUXppfk2e8UCQQDJt5c07BS9i2SDEXiz
byzqCfXVzkdnDj9ry9mba1dcr9B9NCslVelXDGZKvQUBqNYCVxg398aRfWlYDTjU
IoVVAkAbTyjPN6R4SkC4HJMg5oReBmvkwFCAFsemBk0GXwuzD0IlJAjXnAZ+/rIO
ItewfwXIL1Mqz53lO/gK+q6TR585AkB304KUIoWzjyF3JqLP3IQOxzns92u9EV6l
V2P+CkbMPXiZV6sls6I4XppJXX2i3bu7iidN3/dqJ9izQK94fMU9AkBZvgsIPCot
y1/POIbv9LtnviDKrmpkXgVQSU4BmTPvXwTJm8APC7P/horSh3SVf1zgmnsyjm9D
hO92gGc+4ajL
-----END PRIVATE KEY-----`

	region     string
	bucketName string
	objectName string
	letters    = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
	flag.StringVar(&objectName, "object", "", "The name of the object.")
}

func main() {
	flag.Parse()
	if len(bucketName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name required")
	}

	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	if len(objectName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, object name required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := oss.NewClient(cfg)

	mc, err := crypto.CreateMasterRsa(map[string]string{"tag": "value"}, rsaPublicKey, rsaPrivateKey)
	if err != nil {
		log.Fatalf("failed to create master rsa %v", err)
	}
	eclient, err := oss.NewEncryptionClient(client, mc)

	// case 1: Simple Upload
	putObjRequest := &oss.PutObjectRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
		Body:   strings.NewReader("hi, simple put object"),
	}

	putObjRequestResult, err := eclient.PutObject(context.TODO(), putObjRequest)
	if err != nil {
		log.Fatalf("failed to put object with encryption client %v", err)
	}
	log.Printf("put object with encryption client result:%#v\n", putObjRequestResult)

	// case 2: Mutlipart Upload
	length := 3 * 100 * 1024
	partSize := 100 * 1024
	count := length / partSize

	initRequest := &oss.InitiateMultipartUploadRequest{
		Bucket:      oss.Ptr(bucketName),
		Key:         oss.Ptr(objectName),
		CSEDataSize: oss.Ptr(int64(length)),
		CSEPartSize: oss.Ptr(int64(partSize)),
	}
	initResult, err := eclient.InitiateMultipartUpload(context.TODO(), initRequest)
	if err != nil {
		log.Fatalf("failed to initiate multi part upload %v", err)
	}
	var wg sync.WaitGroup
	var parts oss.UploadParts
	body := randStr(length)
	reader := strings.NewReader(body)
	bufReader := bufio.NewReader(reader)
	content, _ := io.ReadAll(bufReader)
	var mu sync.Mutex
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(partNumber int, partSize int, i int) {
			defer wg.Done()
			partRequest := &oss.UploadPartRequest{
				Bucket:              oss.Ptr(bucketName),
				Key:                 oss.Ptr(objectName),
				PartNumber:          int32(partNumber),
				UploadId:            oss.Ptr(*initResult.UploadId),
				Body:                strings.NewReader(string(content[i*partSize : (i+1)*partSize])),
				CSEMultiPartContext: initResult.CSEMultiPartContext,
			}
			partResult, err := eclient.UploadPart(context.TODO(), partRequest)
			if err != nil {
				log.Fatalf("failed to upload part %d: %v", partNumber, err)
			}
			part := oss.UploadPart{
				PartNumber: partRequest.PartNumber,
				ETag:       partResult.ETag,
			}
			mu.Lock()
			parts = append(parts, part)
			mu.Unlock()
		}(i+1, partSize, i)
	}
	wg.Wait()
	sort.Sort(parts)
	request := &oss.CompleteMultipartUploadRequest{
		Bucket:   oss.Ptr(bucketName),
		Key:      oss.Ptr(objectName),
		UploadId: oss.Ptr(*initResult.UploadId),
		CompleteMultipartUpload: &oss.CompleteMultipartUpload{
			Parts: parts,
		},
	}
	result, err := eclient.CompleteMultipartUpload(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to complete multipart upload %v", err)
	}
	log.Printf("complete multipart upload result:%#v\n", result)
}

func randStr(n int) string {
	b := make([]rune, n)
	randMarker := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[randMarker.Intn(len(letters))]
	}
	return string(b)
}
