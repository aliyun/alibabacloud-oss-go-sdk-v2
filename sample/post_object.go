package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"hash"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region     string
	bucketName string
	objectName string
	product    = "oss"
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

	credentialsProvider := credentials.NewEnvironmentVariableCredentialsProvider()
	cred, err := credentialsProvider.GetCredentials(context.TODO())
	if err != nil {
		log.Fatalf("GetCredentials fail, err:%v", err)
	}

	content := "hi oss"

	// build policy
	utcTime := time.Now().UTC()
	date := utcTime.Format("20060102")
	expiration := utcTime.Add(1 * time.Hour)
	policyMap := map[string]any{
		"expiration": expiration.Format("2006-01-02T15:04:05.000Z"),
		"conditions": []any{
			map[string]string{"bucket": bucketName},
			map[string]string{"x-oss-signature-version": "OSS4-HMAC-SHA256"},
			map[string]string{"x-oss-credential": fmt.Sprintf("%v/%v/%v/%v/aliyun_v4_request",
				cred.AccessKeyID, date, region, product)},
			map[string]string{"x-oss-date": utcTime.Format("20060102T150405Z")},
			//other condition
			[]any{"content-length-range", 1, 1024},
			//[]any{"eq", "$success_action_status", "201"},
			//[]any{"starts-with", "$key", "user/eric/"},
			//[]any{"in", "$content-type", []string{"image/jpg", "image/png"}},
			//[]any{"not-in", "$cache-control", []string{"no-cache"}},
		},
	}

	policy, err := json.Marshal(policyMap)
	if err != nil {
		log.Fatalf("json.Marshal fail, err:%v", err)
	}

	// sign policy
	stringToSign := base64.StdEncoding.EncodeToString([]byte(policy))

	// signing key
	hmacHash := func() hash.Hash { return sha256.New() }

	signingKey := "aliyun_v4" + cred.AccessKeySecret
	h1 := hmac.New(hmacHash, []byte(signingKey))
	io.WriteString(h1, date)
	h1Key := h1.Sum(nil)

	h2 := hmac.New(hmacHash, h1Key)
	io.WriteString(h2, region)
	h2Key := h2.Sum(nil)

	h3 := hmac.New(hmacHash, h2Key)
	io.WriteString(h3, product)
	h3Key := h3.Sum(nil)

	h4 := hmac.New(hmacHash, h3Key)
	io.WriteString(h4, "aliyun_v4_request")
	h4Key := h4.Sum(nil)

	// Signature
	h := hmac.New(hmacHash, h4Key)
	io.WriteString(h, stringToSign)
	signature := hex.EncodeToString(h.Sum(nil))

	// Post Request
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// object info, key & metadata
	bodyWriter.WriteField("key", objectName)
	// meta-data
	//bodyWriter.WriteField("x-oss-", value)
	// Policy
	bodyWriter.WriteField("policy", stringToSign)
	// Signature
	bodyWriter.WriteField("x-oss-signature-version", "OSS4-HMAC-SHA256")
	bodyWriter.WriteField("x-oss-credential", fmt.Sprintf("%v/%v/%v/%v/aliyun_v4_request", cred.AccessKeyID, date, region, product))
	bodyWriter.WriteField("x-oss-date", utcTime.Format("20060102T150405Z"))
	bodyWriter.WriteField("x-oss-signature", signature)

	// Data
	w, _ := bodyWriter.CreateFormField("file")
	w.Write([]byte(content))

	bodyWriter.Close()

	req, _ := http.NewRequest("POST", fmt.Sprintf("http://%v.oss-%v.aliyuncs.com/", bucketName, region), bodyBuf)
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	req.WithContext(context.Background())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Do fail, err:%v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		log.Fatalf("Post Object Fail, status code:%v, reson:%v", resp.StatusCode, resp.Status)
	}

	log.Printf("post object done, status code:%v, request id:%v\n", resp.StatusCode, resp.Header.Get("X-Oss-Request-Id"))
}
