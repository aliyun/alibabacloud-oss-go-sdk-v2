package oss

import (
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_ListCloudBoxes(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *ListCloudBoxesRequest
	var input *OperationInput
	var err error

	request = &ListCloudBoxesRequest{}
	input = &OperationInput{
		OpName: "ListCloudBoxes",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"cloudboxes": "",
		},
	}

	input.OpMetadata.Set(signer.SubResource, []string{"cloudboxes"})

	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)

	request = &ListCloudBoxesRequest{
		Marker:  Ptr(""),
		MaxKeys: 10,
		Prefix:  Ptr("/"),
	}
	input = &OperationInput{
		OpName: "ListCloudBoxes",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"cloudboxes": "",
		},
	}

	input.OpMetadata.Set(signer.SubResource, []string{"cloudboxes"})

	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["cloudboxes"], "")
	assert.Equal(t, input.Parameters["marker"], "")
	assert.Equal(t, input.Parameters["max-keys"], "10")
	assert.Equal(t, input.Parameters["prefix"], "/")
}
