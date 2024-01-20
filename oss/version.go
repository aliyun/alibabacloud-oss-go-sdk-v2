package oss

import (
	"bytes"
	"fmt"
	"sync"
)

const (
	major = "1"
	minor = "0"
	patch = "0"
	tag   = "-devpreview"

	SdkName = "alibabacloud-go-sdk-v2"
)

var once sync.Once
var version string

func Version() string {
	once.Do(func() {
		ver := fmt.Sprintf("%s.%s.%s", major, minor, patch)
		verBuilder := bytes.NewBufferString(ver)
		if tag != "" && tag != "-" {
			_, err := verBuilder.WriteString(tag)
			if err != nil {
				verBuilder = bytes.NewBufferString(ver)
			}
		}
		version = verBuilder.String()
	})
	return version
}
