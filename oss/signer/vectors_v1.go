package signer

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

var vectors_requiredSignedParametersMap = map[string]struct{}{
	"acl":                          {},
	"bucketInfo":                   {},
	"location":                     {},
	"stat":                         {},
	"delete":                       {},
	"append":                       {},
	"tagging":                      {},
	"objectMeta":                   {},
	"uploads":                      {},
	"uploadId":                     {},
	"partNumber":                   {},
	"security-token":               {},
	"position":                     {},
	"response-content-type":        {},
	"response-content-language":    {},
	"response-expires":             {},
	"response-cache-control":       {},
	"response-content-disposition": {},
	"response-content-encoding":    {},
	"restore":                      {},
	"callback":                     {},
	"callback-var":                 {},
	"versions":                     {},
	"versioning":                   {},
	"versionId":                    {},
	"sequential":                   {},
	"continuation-token":           {},
	"regionList":                   {},
	"cloudboxes":                   {},
	"symlink":                      {},
	"policy":                       {},
	"logging":                      {},
	"putVectorIndex":               {},
	"getVectorIndex":               {},
	"listVectorIndexes":            {},
	"deleteVectorIndex":            {},
	"putVectors":                   {},
	"getVectors":                   {},
	"listVectors":                  {},
	"deleteVectors":                {},
	"queryVectors":                 {},
}

type SignerVectorsV1 struct {
	AccountId *string
}

func (s *SignerVectorsV1) calcStringToSign(date string, signingCtx *SigningContext) string {
	/*
		SignToString =
			VERB + "\n"
			+ Content-MD5 + "\n"
			+ Content-Type + "\n"
			+ Date + "\n"
			+ CanonicalizedOSSHeaders
			+ CanonicalizedResource
		Signature = base64(hmac-sha1(AccessKeySecret, SignToString))
	*/
	request := signingCtx.Request

	contentMd5 := request.Header.Get(contentMd5Header)
	contentType := request.Header.Get(contentTypeHeader)

	//CanonicalizedOSSHeaders
	var headers []string
	for k := range request.Header {
		lowerCaseKey := strings.ToLower(k)
		if strings.HasPrefix(lowerCaseKey, ossHeaderPreifx) {
			headers = append(headers, lowerCaseKey)
		}
	}
	sort.Strings(headers)
	headerItems := make([]string, len(headers))
	for i, k := range headers {
		headerValues := make([]string, len(request.Header.Values(k)))
		for i, v := range request.Header.Values(k) {
			headerValues[i] = strings.TrimSpace(v)
		}
		headerItems[i] = k + ":" + strings.Join(headerValues, ",") + "\n"
	}
	canonicalizedOSSHeaders := strings.Join(headerItems, "")

	//CanonicalizedResource
	query := request.URL.Query()
	var params []string
	for k := range query {
		if _, ok := vectors_requiredSignedParametersMap[k]; ok {
			params = append(params, k)
		} else if strings.HasPrefix(k, ossHeaderPreifx) {
			params = append(params, k)
		} else if isSubResource(signingCtx.SubResource, k) {
			params = append(params, k)
		}
	}
	sort.Strings(params)
	paramItems := make([]string, len(params))
	for i, k := range params {
		v := query.Get(k)
		if len(v) > 0 {
			paramItems[i] = k + "=" + v
		} else {
			paramItems[i] = k
		}
	}
	subResource := strings.Join(paramItems, "&")
	arn := fmt.Sprintf("/acs:ossvector:%s:", toString(signingCtx.Region))
	if signingCtx.Bucket != nil {
		arn += fmt.Sprintf("%s:%s/", toString(s.AccountId), *signingCtx.Bucket)
	} else {
		// service's api without account id
		arn += ":/"
	}
	canonicalizedResource := arn

	if signingCtx.Key != nil {
		canonicalizedResource += *signingCtx.Key
	}
	if subResource != "" {
		canonicalizedResource += "?" + subResource
	}

	// string to Sign
	stringToSign :=
		request.Method + "\n" +
			contentMd5 + "\n" +
			contentType + "\n" +
			date + "\n" +
			canonicalizedOSSHeaders +
			canonicalizedResource

	//fmt.Printf("stringToSign:%s\n", stringToSign)
	return stringToSign
}

func (s *SignerVectorsV1) authHeader(ctx context.Context, signingCtx *SigningContext) error {
	request := signingCtx.Request
	cred := signingCtx.Credentials

	// Date
	if signingCtx.Time.IsZero() {
		signingCtx.Time = time.Now().Add(signingCtx.ClockOffset)
	}
	datetime := signingCtx.Time.UTC().Format(http.TimeFormat)
	request.Header.Set(dateHeader, datetime)

	// Credentials information
	if cred.SecurityToken != "" {
		request.Header.Set(securityTokenHeader, cred.SecurityToken)
	}

	// StringToSign
	stringToSign := s.calcStringToSign(datetime, signingCtx)
	signingCtx.StringToSign = stringToSign

	// Signature
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(cred.AccessKeySecret))
	if _, err := io.WriteString(h, stringToSign); err != nil {
		return err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Authorization header
	request.Header.Set(authorizationHeader, fmt.Sprintf("OSS %s:%s", cred.AccessKeyID, signature))

	return nil
}

func (s *SignerVectorsV1) authQuery(ctx context.Context, signingCtx *SigningContext) error {
	request := signingCtx.Request
	cred := signingCtx.Credentials

	// Date
	if signingCtx.Time.IsZero() {
		signingCtx.Time = time.Now().UTC().Add(defaultExpiresDuration)
	}
	datetime := fmt.Sprintf("%v", signingCtx.Time.UTC().Unix())

	// Credentials information
	query, _ := url.ParseQuery(request.URL.RawQuery)
	if cred.SecurityToken != "" {
		query.Add(securityTokenQuery, cred.SecurityToken)
		request.URL.RawQuery = query.Encode()
	}

	// StringToSign
	stringToSign := s.calcStringToSign(datetime, signingCtx)
	signingCtx.StringToSign = stringToSign

	// Signature
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(cred.AccessKeySecret))
	if _, err := io.WriteString(h, stringToSign); err != nil {
		return err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Authorization query
	query.Add(expiresQuery, datetime)
	query.Add(accessKeyIdQuery, cred.AccessKeyID)
	query.Add(signatureQuery, signature)
	request.URL.RawQuery = strings.Replace(query.Encode(), "+", "%20", -1)

	return nil
}

func (s *SignerVectorsV1) Sign(ctx context.Context, signingCtx *SigningContext) error {
	if signingCtx == nil {
		return fmt.Errorf("SigningContext is null.")
	}

	if signingCtx.Credentials == nil || !signingCtx.Credentials.HasKeys() {
		return fmt.Errorf("SigningContext.Credentials is null or empty.")
	}

	if signingCtx.Request == nil {
		return fmt.Errorf("SigningContext.Request is null.")
	}

	if signingCtx.AuthMethodQuery {
		return s.authQuery(ctx, signingCtx)
	}

	return s.authHeader(ctx, signingCtx)
}

func (*SignerVectorsV1) IsSignedHeader(additionalHeaders []string, h string) bool {
	lowerCaseKey := strings.ToLower(h)
	if strings.HasPrefix(lowerCaseKey, ossHeaderPreifx) ||
		lowerCaseKey == "date" ||
		lowerCaseKey == "content-type" ||
		lowerCaseKey == "content-md5" {
		return true
	}
	return false
}
