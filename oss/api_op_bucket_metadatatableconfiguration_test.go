package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_CreateBucketMetadataTableConfiguration(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *CreateBucketMetadataTableConfigurationRequest
	var input *OperationInput
	var err error

	request = &CreateBucketMetadataTableConfigurationRequest{}
	input = &OperationInput{
		OpName: "CreateBucketMetadataTableConfiguration",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"metadataConfiguration": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &CreateBucketMetadataTableConfigurationRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "CreateBucketMetadataTableConfiguration",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"metadataConfiguration": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, MetadataConfiguration.")

	request = &CreateBucketMetadataTableConfigurationRequest{
		Bucket: Ptr("oss-demo"),
		MetadataConfiguration: &MetadataConfiguration{
			JournalTableConfiguration: &JournalTableConfiguration{
				RecordExpiration: &RecordExpiration{
					Expiration: Ptr("ENABLED"),
					Days:       Ptr(7),
				},
			},
			InventoryTableConfiguration: &InventoryTableConfiguration{
				ConfigurationState: Ptr("ENABLED"),
				EncryptionConfiguration: &MetadataTableEncryptionConfiguration{
					SseAlgorithm: Ptr("AES256"),
				},
			},
		},
	}
	input = &OperationInput{
		OpName: "CreateBucketMetadataTableConfiguration",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"metadataConfiguration": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.Nil(t, err)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t,
		"<MetadataConfiguration><JournalTableConfiguration><RecordExpiration><Expiration>ENABLED</Expiration><Days>7</Days></RecordExpiration></JournalTableConfiguration><InventoryTableConfiguration><ConfigurationState>ENABLED</ConfigurationState><EncryptionConfiguration><SseAlgorithm>AES256</SseAlgorithm></EncryptionConfiguration></InventoryTableConfiguration></MetadataConfiguration>",
		string(body))
}

func TestUnmarshalOutput_CreateBucketMetadataTableConfiguration(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &CreateBucketMetadataTableConfigurationResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "OK", result.Status)
	assert.Equal(t, "534B371674E88A4D8906****", result.Headers.Get("X-Oss-Request-Id"))

	body := `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>MetadataConfigurationAlreadyExists</Code>
  <Message>The metadata table configuration already exists.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
</Error>`
	output = &OperationOutput{
		StatusCode: 409,
		Status:     "Conflict",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &CreateBucketMetadataTableConfigurationResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, 409, result.StatusCode)
	assert.Equal(t, "Conflict", result.Status)
}

func TestMarshalInput_GetBucketMetadataTableConfiguration(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetBucketMetadataTableConfigurationRequest
	var input *OperationInput
	var err error

	request = &GetBucketMetadataTableConfigurationRequest{}
	input = &OperationInput{
		OpName: "GetBucketMetadataTableConfiguration",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"metadataConfiguration": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetBucketMetadataTableConfigurationRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetBucketMetadataTableConfiguration",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"metadataConfiguration": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetBucketMetadataTableConfiguration(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error

	body := `<?xml version="1.0" encoding="UTF-8"?>
<GetBucketMetadataConfigurationResult>
  <MetadataConfigurationResult>
    <DestinationResult>
      <TableBucketType>oss</TableBucketType>
      <TableBucketArn>arn:acs:oss:cn-hangzhou:123456789012:bucket/examplebucket</TableBucketArn>
      <TableNamespace>b_examplebucket</TableNamespace>
    </DestinationResult>
    <JournalTableConfigurationResult>
      <TableStatus>ACTIVE</TableStatus>
      <TableName>journal</TableName>
      <TableArn>arn:acs:oss:cn-hangzhou:123456789012:bucket/examplebucket/journal</TableArn>
      <RecordExpiration>
        <Expiration>ENABLED</Expiration>
        <Days>7</Days>
      </RecordExpiration>
      <EncryptionConfiguration>
        <SseAlgorithm>AES256</SseAlgorithm>
      </EncryptionConfiguration>
    </JournalTableConfigurationResult>
    <InventoryTableConfigurationResult>
      <ConfigurationState>ENABLED</ConfigurationState>
      <TableStatus>ACTIVE</TableStatus>
      <TableName>inventory</TableName>
      <TableArn>arn:acs:oss:cn-hangzhou:123456789012:bucket/examplebucket/inventory</TableArn>
      <EncryptionConfiguration>
        <SseAlgorithm>AES256</SseAlgorithm>
      </EncryptionConfiguration>
    </InventoryTableConfigurationResult>
  </MetadataConfigurationResult>
</GetBucketMetadataConfigurationResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &GetBucketMetadataTableConfigurationResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "OK", result.Status)
	assert.NotNil(t, result.GetBucketMetadataConfigurationResult)
	assert.Equal(t, "oss", *result.GetBucketMetadataConfigurationResult.MetadataConfigurationResult.DestinationResult.TableBucketType)
	assert.Equal(t, "b_examplebucket", *result.GetBucketMetadataConfigurationResult.MetadataConfigurationResult.DestinationResult.TableNamespace)
	assert.Equal(t, "ACTIVE", *result.GetBucketMetadataConfigurationResult.MetadataConfigurationResult.JournalTableConfigurationResult.TableStatus)
	assert.Equal(t, "journal", *result.GetBucketMetadataConfigurationResult.MetadataConfigurationResult.JournalTableConfigurationResult.TableName)
	assert.Equal(t, "ENABLED", *result.GetBucketMetadataConfigurationResult.MetadataConfigurationResult.JournalTableConfigurationResult.RecordExpiration.Expiration)
	assert.Equal(t, int(7), *result.GetBucketMetadataConfigurationResult.MetadataConfigurationResult.JournalTableConfigurationResult.RecordExpiration.Days)
	assert.Equal(t, "AES256", *result.GetBucketMetadataConfigurationResult.MetadataConfigurationResult.JournalTableConfigurationResult.EncryptionConfiguration.SseAlgorithm)
	assert.Equal(t, "ENABLED", *result.GetBucketMetadataConfigurationResult.MetadataConfigurationResult.InventoryTableConfigurationResult.ConfigurationState)
	assert.Equal(t, "ACTIVE", *result.GetBucketMetadataConfigurationResult.MetadataConfigurationResult.InventoryTableConfigurationResult.TableStatus)
	assert.Equal(t, "inventory", *result.GetBucketMetadataConfigurationResult.MetadataConfigurationResult.InventoryTableConfigurationResult.TableName)

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>MetadataConfigurationNotFound</Code>
  <Message>The metadata table configuration does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
</Error>`
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "Not Found",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetBucketMetadataTableConfigurationResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, 404, result.StatusCode)
}

func TestMarshalInput_UpdateBucketMetadataInventoryTableConfiguration(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *UpdateBucketMetadataInventoryTableConfigurationRequest
	var input *OperationInput
	var err error

	request = &UpdateBucketMetadataInventoryTableConfigurationRequest{}
	input = &OperationInput{
		OpName: "UpdateBucketMetadataInventoryTableConfiguration",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"metadataInventoryTable": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &UpdateBucketMetadataInventoryTableConfigurationRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "UpdateBucketMetadataInventoryTableConfiguration",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"metadataInventoryTable": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, InventoryTableConfiguration.")

	request = &UpdateBucketMetadataInventoryTableConfigurationRequest{
		Bucket: Ptr("oss-demo"),
		InventoryTableConfiguration: &InventoryTableConfiguration{
			ConfigurationState: Ptr("DISABLED"),
		},
	}
	input = &OperationInput{
		OpName: "UpdateBucketMetadataInventoryTableConfiguration",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"metadataInventoryTable": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.Nil(t, err)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t,
		"<InventoryTableConfiguration><ConfigurationState>DISABLED</ConfigurationState></InventoryTableConfiguration>",
		string(body))
}

func TestUnmarshalOutput_UpdateBucketMetadataInventoryTableConfiguration(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	output := &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &UpdateBucketMetadataInventoryTableConfigurationResult{}
	err := c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "OK", result.Status)
}

func TestMarshalInput_UpdateBucketMetadataJournalTableConfiguration(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *UpdateBucketMetadataJournalTableConfigurationRequest
	var input *OperationInput
	var err error

	request = &UpdateBucketMetadataJournalTableConfigurationRequest{}
	input = &OperationInput{
		OpName: "UpdateBucketMetadataJournalTableConfiguration",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"metadataJournalTable": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &UpdateBucketMetadataJournalTableConfigurationRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "UpdateBucketMetadataJournalTableConfiguration",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"metadataJournalTable": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, JournalTableConfiguration.")

	request = &UpdateBucketMetadataJournalTableConfigurationRequest{
		Bucket: Ptr("oss-demo"),
		JournalTableConfiguration: &JournalTableConfiguration{
			RecordExpiration: &RecordExpiration{
				Expiration: Ptr("DISABLED"),
			},
		},
	}
	input = &OperationInput{
		OpName: "UpdateBucketMetadataJournalTableConfiguration",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"metadataJournalTable": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.Nil(t, err)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t,
		"<JournalTableConfiguration><RecordExpiration><Expiration>DISABLED</Expiration></RecordExpiration></JournalTableConfiguration>",
		string(body))
}

func TestUnmarshalOutput_UpdateBucketMetadataJournalTableConfiguration(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	output := &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &UpdateBucketMetadataJournalTableConfigurationResult{}
	err := c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "OK", result.Status)
}

func TestMarshalInput_DeleteBucketMetadataTableConfiguration(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DeleteBucketMetadataTableConfigurationRequest
	var input *OperationInput
	var err error

	request = &DeleteBucketMetadataTableConfigurationRequest{}
	input = &OperationInput{
		OpName: "DeleteBucketMetadataTableConfiguration",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"metadataConfiguration": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &DeleteBucketMetadataTableConfigurationRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "DeleteBucketMetadataTableConfiguration",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"metadataConfiguration": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, MarshalUpdateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_DeleteBucketMetadataTableConfiguration(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error

	output = &OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &DeleteBucketMetadataTableConfigurationResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, 204, result.StatusCode)
	assert.Equal(t, "No Content", result.Status)

	body := `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>MetadataConfigurationNotFound</Code>
  <Message>The metadata table configuration does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
</Error>`
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "Not Found",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DeleteBucketMetadataTableConfigurationResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, 404, result.StatusCode)
}
