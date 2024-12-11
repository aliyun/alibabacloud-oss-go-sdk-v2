package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_CreateBucketDataRedundancyTransition(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *CreateBucketDataRedundancyTransitionRequest
	var input *OperationInput
	var err error

	request = &CreateBucketDataRedundancyTransitionRequest{}
	input = &OperationInput{
		OpName: "CreateBucketDataRedundancyTransition",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"redundancyTransition": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"redundancyTransition"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &CreateBucketDataRedundancyTransitionRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "CreateBucketDataRedundancyTransition",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"redundancyTransition": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"redundancyTransition"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Contains(t, err.Error(), "missing required field, TargetRedundancyType.")

	request = &CreateBucketDataRedundancyTransitionRequest{
		Bucket:               Ptr("oss-demo"),
		TargetRedundancyType: Ptr("ZRS"),
	}
	input = &OperationInput{
		OpName: "CreateBucketDataRedundancyTransition",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"redundancyTransition": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"redundancyTransition"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-target-redundancy-type"], "ZRS")
}

func TestUnmarshalOutput_CreateBucketDataRedundancyTransition(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
			<BucketDataRedundancyTransition>
			  <TaskId>4be5beb0f74f490186311b268bf6****</TaskId>
			</BucketDataRedundancyTransition>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &CreateBucketDataRedundancyTransitionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.BucketDataRedundancyTransition.TaskId, "4be5beb0f74f490186311b268bf6****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &CreateBucketDataRedundancyTransitionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &CreateBucketDataRedundancyTransitionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &CreateBucketDataRedundancyTransitionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetBucketDataRedundancyTransition(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetBucketDataRedundancyTransitionRequest
	var input *OperationInput
	var err error

	request = &GetBucketDataRedundancyTransitionRequest{}
	input = &OperationInput{
		OpName: "GetBucketDataRedundancyTransition",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"redundancyTransition": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"redundancyTransition"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetBucketDataRedundancyTransitionRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetBucketDataRedundancyTransition",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"redundancyTransition": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"redundancyTransition"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, RedundancyTransitionTaskid.")

	request = &GetBucketDataRedundancyTransitionRequest{
		Bucket:                     Ptr("oss-demo"),
		RedundancyTransitionTaskid: Ptr("123"),
	}
	input = &OperationInput{
		OpName: "GetBucketDataRedundancyTransition",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"redundancyTransition": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"redundancyTransition"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-redundancy-transition-taskid"], "123")
}

func TestUnmarshalOutput_GetBucketDataRedundancyTransition(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	putBody := `<?xml version="1.0" encoding="UTF-8"?>
<BucketDataRedundancyTransition>
  <Bucket>examplebucket</Bucket>
  <TaskId>4be5beb0f74f490186311b268bf6****</TaskId>
  <Status>Queueing</Status>
  <CreateTime>2023-11-17T09:11:58.000Z</CreateTime>
</BucketDataRedundancyTransition>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	assert.Nil(t, err)
	result := &GetBucketDataRedundancyTransitionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.BucketDataRedundancyTransition.Bucket, "examplebucket")
	assert.Equal(t, *result.BucketDataRedundancyTransition.TaskId, "4be5beb0f74f490186311b268bf6****")
	assert.Equal(t, *result.BucketDataRedundancyTransition.Status, "Queueing")
	assert.Equal(t, *result.BucketDataRedundancyTransition.CreateTime, "2023-11-17T09:11:58.000Z")

	putBody = `<?xml version="1.0" encoding="UTF-8"?>
<BucketDataRedundancyTransition>
  <Bucket>examplebucket</Bucket>
  <TaskId>909c6c818dd041d1a44e0fdc66aa****</TaskId>
  <Status>Processing</Status>
  <CreateTime>2023-11-17T09:14:39.000Z</CreateTime>
  <StartTime>2023-11-17T09:14:39.000Z</StartTime>
  <ProcessPercentage>0</ProcessPercentage>
  <EstimatedRemainingTime>100</EstimatedRemainingTime>
</BucketDataRedundancyTransition>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	assert.Nil(t, err)
	result = &GetBucketDataRedundancyTransitionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.BucketDataRedundancyTransition.Bucket, "examplebucket")
	assert.Equal(t, *result.BucketDataRedundancyTransition.TaskId, "909c6c818dd041d1a44e0fdc66aa****")
	assert.Equal(t, *result.BucketDataRedundancyTransition.Status, "Processing")
	assert.Equal(t, *result.BucketDataRedundancyTransition.CreateTime, "2023-11-17T09:14:39.000Z")
	assert.Equal(t, *result.BucketDataRedundancyTransition.StartTime, "2023-11-17T09:14:39.000Z")
	assert.Equal(t, *result.BucketDataRedundancyTransition.ProcessPercentage, int32(0))
	assert.Equal(t, *result.BucketDataRedundancyTransition.EstimatedRemainingTime, int64(100))

	putBody = `<?xml version="1.0" encoding="UTF-8"?>
<BucketDataRedundancyTransition>
  <Bucket>examplebucket</Bucket>
  <TaskId>909c6c818dd041d1a44e0fdc66aa****</TaskId>
  <Status>Finished</Status>
  <CreateTime>2023-11-17T09:14:39.000Z</CreateTime>
  <StartTime>2023-11-17T09:14:39.000Z</StartTime>
  <ProcessPercentage>100</ProcessPercentage>
  <EstimatedRemainingTime>0</EstimatedRemainingTime>
  <EndTime>2023-11-18T09:14:39.000Z</EndTime>
</BucketDataRedundancyTransition>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	assert.Nil(t, err)
	result = &GetBucketDataRedundancyTransitionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.BucketDataRedundancyTransition.Bucket, "examplebucket")
	assert.Equal(t, *result.BucketDataRedundancyTransition.TaskId, "909c6c818dd041d1a44e0fdc66aa****")
	assert.Equal(t, *result.BucketDataRedundancyTransition.Status, "Finished")
	assert.Equal(t, *result.BucketDataRedundancyTransition.CreateTime, "2023-11-17T09:14:39.000Z")
	assert.Equal(t, *result.BucketDataRedundancyTransition.StartTime, "2023-11-17T09:14:39.000Z")
	assert.Equal(t, *result.BucketDataRedundancyTransition.ProcessPercentage, int32(100))
	assert.Equal(t, *result.BucketDataRedundancyTransition.EstimatedRemainingTime, int64(0))
	assert.Equal(t, *result.BucketDataRedundancyTransition.EndTime, "2023-11-18T09:14:39.000Z")

	putBody = `<?xml version="1.0" encoding="UTF-8"?>
		<Error>
		<Code>NoSuchBucket</Code>
		<Message>The specified bucket does not exist.</Message>
		<RequestId>66C2FF09FDF07830343C72EC</RequestId>
		<HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
		<BucketName>bucket</BucketName>
		<EC>0015-00000101</EC>
	</Error>`
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	assert.Nil(t, err)
	result = &GetBucketDataRedundancyTransitionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	putBody = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	assert.Nil(t, err)
	result = &GetBucketDataRedundancyTransitionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_DeleteBucketDataRedundancyTransition(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DeleteBucketDataRedundancyTransitionRequest
	var input *OperationInput
	var err error

	request = &DeleteBucketDataRedundancyTransitionRequest{}
	input = &OperationInput{
		OpName: "DeleteBucketDataRedundancyTransition",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"redundancyTransition": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"redundancyTransition"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &DeleteBucketDataRedundancyTransitionRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "DeleteBucketDataRedundancyTransition",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"redundancyTransition": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"redundancyTransition"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Contains(t, err.Error(), "missing required field, RedundancyTransitionTaskid.")

	request = &DeleteBucketDataRedundancyTransitionRequest{
		Bucket:                     Ptr("oss-demo"),
		RedundancyTransitionTaskid: Ptr("123"),
	}
	input = &OperationInput{
		OpName: "DeleteBucketDataRedundancyTransition",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"redundancyTransition": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"redundancyTransition"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-redundancy-transition-taskid"], "123")
}

func TestUnmarshalOutput_DeleteBucketDataRedundancyTransition(t *testing.T) {
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
	result := &DeleteBucketDataRedundancyTransitionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DeleteBucketDataRedundancyTransitionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DeleteBucketDataRedundancyTransitionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body := `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DeleteBucketDataRedundancyTransitionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_ListBucketDataRedundancyTransition(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *ListBucketDataRedundancyTransitionRequest
	var input *OperationInput
	var err error

	request = &ListBucketDataRedundancyTransitionRequest{}
	input = &OperationInput{
		OpName: "ListBucketDataRedundancyTransition",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"redundancyTransition": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"redundancyTransition"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &ListBucketDataRedundancyTransitionRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetBucketDataRedundancyTransition",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"redundancyTransition": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"redundancyTransition"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_ListBucketDataRedundancyTransition(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	putBody := `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketDataRedundancyTransition>
<BucketDataRedundancyTransition>
  <Bucket>examplebucket</Bucket>
  <TaskId>4be5beb0f74f490186311b268bf6****</TaskId>
  <Status>Queueing</Status>
  <CreateTime>2023-11-17T09:11:58.000Z</CreateTime>
</BucketDataRedundancyTransition>
</ListBucketDataRedundancyTransition>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	assert.Nil(t, err)
	result := &ListBucketDataRedundancyTransitionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].Bucket, "examplebucket")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].TaskId, "4be5beb0f74f490186311b268bf6****")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].Status, "Queueing")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].CreateTime, "2023-11-17T09:11:58.000Z")

	putBody = `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketDataRedundancyTransition>
<BucketDataRedundancyTransition>
  <Bucket>examplebucket</Bucket>
  <TaskId>909c6c818dd041d1a44e0fdc66aa****</TaskId>
  <Status>Processing</Status>
  <CreateTime>2023-11-17T09:14:39.000Z</CreateTime>
  <StartTime>2023-11-17T09:14:39.000Z</StartTime>
  <ProcessPercentage>0</ProcessPercentage>
  <EstimatedRemainingTime>100</EstimatedRemainingTime>
</BucketDataRedundancyTransition>
</ListBucketDataRedundancyTransition>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	assert.Nil(t, err)
	result = &ListBucketDataRedundancyTransitionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].Bucket, "examplebucket")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].TaskId, "909c6c818dd041d1a44e0fdc66aa****")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].Status, "Processing")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].CreateTime, "2023-11-17T09:14:39.000Z")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].StartTime, "2023-11-17T09:14:39.000Z")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].ProcessPercentage, int32(0))
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].EstimatedRemainingTime, int64(100))

	putBody = `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketDataRedundancyTransition>
<BucketDataRedundancyTransition>
  <Bucket>examplebucket</Bucket>
  <TaskId>909c6c818dd041d1a44e0fdc66aa****</TaskId>
  <Status>Finished</Status>
  <CreateTime>2023-11-17T09:14:39.000Z</CreateTime>
  <StartTime>2023-11-17T09:14:39.000Z</StartTime>
  <ProcessPercentage>100</ProcessPercentage>
  <EstimatedRemainingTime>0</EstimatedRemainingTime>
  <EndTime>2023-11-18T09:14:39.000Z</EndTime>
</BucketDataRedundancyTransition>
</ListBucketDataRedundancyTransition>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	assert.Nil(t, err)
	result = &ListBucketDataRedundancyTransitionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].Bucket, "examplebucket")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].TaskId, "909c6c818dd041d1a44e0fdc66aa****")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].Status, "Finished")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].CreateTime, "2023-11-17T09:14:39.000Z")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].StartTime, "2023-11-17T09:14:39.000Z")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].ProcessPercentage, int32(100))
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].EstimatedRemainingTime, int64(0))
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].EndTime, "2023-11-18T09:14:39.000Z")

	putBody = `<?xml version="1.0" encoding="UTF-8"?>
		<Error>
		<Code>NoSuchBucket</Code>
		<Message>The specified bucket does not exist.</Message>
		<RequestId>66C2FF09FDF07830343C72EC</RequestId>
		<HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
		<BucketName>bucket</BucketName>
		<EC>0015-00000101</EC>
	</Error>`
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	assert.Nil(t, err)
	result = &ListBucketDataRedundancyTransitionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	putBody = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	assert.Nil(t, err)
	result = &ListBucketDataRedundancyTransitionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_ListUserDataRedundancyTransition(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *ListUserDataRedundancyTransitionRequest
	var input *OperationInput
	var err error

	request = &ListUserDataRedundancyTransitionRequest{}
	input = &OperationInput{
		OpName: "ListUserDataRedundancyTransition",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"redundancyTransition": "",
		},
	}
	input.OpMetadata.Set(signer.SubResource, []string{"redundancyTransition"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)

	request = &ListUserDataRedundancyTransitionRequest{
		ContinuationToken: Ptr("123"),
		MaxKeys:           int32(10),
	}
	input = &OperationInput{
		OpName: "ListUserDataRedundancyTransition",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"redundancyTransition": "",
		},
	}
	input.OpMetadata.Set(signer.SubResource, []string{"redundancyTransition"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["continuation-token"], "123")
	assert.Equal(t, input.Parameters["max-keys"], "10")
}

func TestUnmarshalOutput_ListUserDataRedundancyTransition(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	putBody := `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketDataRedundancyTransition>
<BucketDataRedundancyTransition>
  <Bucket>examplebucket</Bucket>
  <TaskId>4be5beb0f74f490186311b268bf6****</TaskId>
  <Status>Queueing</Status>
  <CreateTime>2023-11-17T09:11:58.000Z</CreateTime>
</BucketDataRedundancyTransition>
<BucketDataRedundancyTransition>
  <Bucket>examplebucket1</Bucket>
  <TaskId>909c6c818dd041d1a44e0fdc66aa****</TaskId>
  <Status>Processing</Status>
  <CreateTime>2023-11-17T09:14:39.000Z</CreateTime>
  <StartTime>2023-11-17T09:14:39.000Z</StartTime>
  <ProcessPercentage>0</ProcessPercentage>
  <EstimatedRemainingTime>100</EstimatedRemainingTime>
</BucketDataRedundancyTransition>
<BucketDataRedundancyTransition>
  <Bucket>examplebucket2</Bucket>
  <TaskId>909c6c818dd041d1a44e0fdc66aa****</TaskId>
  <Status>Finished</Status>
  <CreateTime>2023-11-17T09:14:39.000Z</CreateTime>
  <StartTime>2023-11-17T09:14:39.000Z</StartTime>
  <ProcessPercentage>100</ProcessPercentage>
  <EstimatedRemainingTime>0</EstimatedRemainingTime>
  <EndTime>2023-11-18T09:14:39.000Z</EndTime>
</BucketDataRedundancyTransition>
<IsTruncated>false</IsTruncated>
<NextContinuationToken></NextContinuationToken>
</ListBucketDataRedundancyTransition>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	assert.Nil(t, err)
	result := &ListUserDataRedundancyTransitionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, "", *result.ListBucketDataRedundancyTransition.NextContinuationToken)
	assert.False(t, *result.ListBucketDataRedundancyTransition.IsTruncated)
	assert.Equal(t, 3, len(result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions))
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].Bucket, "examplebucket")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].TaskId, "4be5beb0f74f490186311b268bf6****")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].Status, "Queueing")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].CreateTime, "2023-11-17T09:11:58.000Z")

	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[1].Bucket, "examplebucket1")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[1].TaskId, "909c6c818dd041d1a44e0fdc66aa****")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[1].Status, "Processing")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[1].CreateTime, "2023-11-17T09:14:39.000Z")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[1].StartTime, "2023-11-17T09:14:39.000Z")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[1].ProcessPercentage, int32(0))
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[1].EstimatedRemainingTime, int64(100))

	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].Bucket, "examplebucket2")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].TaskId, "909c6c818dd041d1a44e0fdc66aa****")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].Status, "Finished")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].CreateTime, "2023-11-17T09:14:39.000Z")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].StartTime, "2023-11-17T09:14:39.000Z")
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].ProcessPercentage, int32(100))
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].EstimatedRemainingTime, int64(0))
	assert.Equal(t, *result.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].EndTime, "2023-11-18T09:14:39.000Z")

	putBody = `<?xml version="1.0" encoding="UTF-8"?>
		<Error>
		<Code>NoSuchBucket</Code>
		<Message>The specified bucket does not exist.</Message>
		<RequestId>66C2FF09FDF07830343C72EC</RequestId>
		<HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
		<BucketName>bucket</BucketName>
		<EC>0015-00000101</EC>
	</Error>`
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	assert.Nil(t, err)
	result = &ListUserDataRedundancyTransitionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	putBody = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(putBody))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	assert.Nil(t, err)
	result = &ListUserDataRedundancyTransitionResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}
