package oss

import (
	"bytes"
	"encoding/hex"
	"hash/crc32"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func hexStrToByte(hexString string) string {
	byteData, _ := hex.DecodeString(hexString)
	return string(byteData)
}

func TestMarshalInput_CreateSelectObjectMeta(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *CreateSelectObjectMetaRequest
	var input *OperationInput
	var err error

	request = &CreateSelectObjectMetaRequest{}
	input = &OperationInput{
		OpName: "CreateSelectObjectMeta",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
	}
	err = c.marshalInput(request, input, marshalMetaRequest, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &CreateSelectObjectMetaRequest{
		Bucket: Ptr("oss-demo"),
	}
	err = c.marshalInput(request, input, marshalMetaRequest, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &CreateSelectObjectMetaRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
	}
	err = c.marshalInput(request, input, marshalMetaRequest, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &CreateSelectObjectMetaRequest{
		Bucket:      Ptr("oss-demo"),
		Key:         Ptr("oss-key"),
		MetaRequest: &SelectRequest{},
	}
	err = c.marshalInput(request, input, marshalMetaRequest, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid field, MetaRequest.")

	request = &CreateSelectObjectMetaRequest{
		Bucket:      Ptr("oss-demo"),
		Key:         Ptr("oss-demo"),
		MetaRequest: &CsvMetaRequest{},
	}
	err = c.marshalInput(request, input, marshalMetaRequest, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-process"], "csv/meta")
	data, _ := io.ReadAll(input.Body)
	assert.Equal(t, "<CsvMetaRequest></CsvMetaRequest>", string(data))

	request = &CreateSelectObjectMetaRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
		MetaRequest: &CsvMetaRequest{
			OverwriteIfExists: Ptr(true),
		},
	}
	err = c.marshalInput(request, input, marshalMetaRequest, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-process"], "csv/meta")
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, "<CsvMetaRequest><OverwriteIfExists>true</OverwriteIfExists></CsvMetaRequest>", string(data))

	request = &CreateSelectObjectMetaRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
		MetaRequest: &CsvMetaRequest{
			OverwriteIfExists: Ptr(true),
			InputSerialization: &InputSerialization{
				CSV: &InputSerializationCSV{
					RecordDelimiter: Ptr("\n"),
					FieldDelimiter:  Ptr(","),
					QuoteCharacter:  Ptr("\""),
				},
			},
		},
	}

	err = c.marshalInput(request, input, marshalMetaRequest, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-process"], "csv/meta")
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, "<CsvMetaRequest><InputSerialization><CSV><RecordDelimiter>Cg==</RecordDelimiter><FieldDelimiter>LA==</FieldDelimiter><QuoteCharacter>Ig==</QuoteCharacter></CSV></InputSerialization><OverwriteIfExists>true</OverwriteIfExists></CsvMetaRequest>", string(data))

	request = &CreateSelectObjectMetaRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
		MetaRequest: &CsvMetaRequest{
			OverwriteIfExists: Ptr(true),
			InputSerialization: &InputSerialization{
				CSV: &InputSerializationCSV{
					RecordDelimiter: Ptr("\n"),
					FieldDelimiter:  Ptr(","),
					QuoteCharacter:  Ptr("\""),
				},
				CompressionType: Ptr("None"),
			},
		},
	}
	err = c.marshalInput(request, input, marshalMetaRequest, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-process"], "csv/meta")
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, "<CsvMetaRequest><InputSerialization><CSV><RecordDelimiter>Cg==</RecordDelimiter><FieldDelimiter>LA==</FieldDelimiter><QuoteCharacter>Ig==</QuoteCharacter></CSV><CompressionType>None</CompressionType></InputSerialization><OverwriteIfExists>true</OverwriteIfExists></CsvMetaRequest>", string(data))

	request = &CreateSelectObjectMetaRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
		MetaRequest: &JsonMetaRequest{
			InputSerialization: &InputSerialization{
				JSON: &InputSerializationJSON{
					JSONType: Ptr("LINES"),
				},
			},
		},
	}
	err = c.marshalInput(request, input, marshalMetaRequest, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-process"], "json/meta")
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, "<JsonMetaRequest><InputSerialization><JSON><Type>LINES</Type></JSON></InputSerialization></JsonMetaRequest>", string(data))
}

func TestUnmarshalOutput_CreateSelectObjectMeta(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	hexString := "01800006000000250000000000000000000000000000000000000000000000c8000000010000000000000130000000182e46a93f70"
	body := hexStrToByte(hexString)
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":     {"5C06A3B67B8B5A3DA422****"},
			"Date":                 {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"Content-Type":         {"application/vnd.ms-excel"},
			"ETag":                 {"\"D41D8CD98F00B204E9800998ECF8****\""},
			"x-oss-hash-crc64ecma": {"316181249502703****"},
			"Content-MD5":          {"1B2M2Y8AsgTpgAmY7PhC****"},
		},
		Body: io.NopCloser(strings.NewReader(string(body))),
	}
	result := &CreateSelectObjectMetaResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyCreateSelectObjectMeta)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C06A3B67B8B5A3DA422****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	assert.Equal(t, result.TotalScanned, int64(0))
	assert.Equal(t, result.MetaStatus, 200)
	assert.Equal(t, result.SplitsCount, int32(1))
	assert.Equal(t, result.RowsCount, int64(304))
	assert.Equal(t, result.ColumnsCount, int32(24))
	assert.Equal(t, result.ErrorMsg, ".")

	hexString = "01800007000000210000000000000000000278f700000000000278f7000000c80000000100000000000000642e6e6a03f9"

	body = hexStrToByte(hexString)
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":     {"5C06A3B67B8B5A3DA422****"},
			"Date":                 {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"Content-Type":         {"application/json"},
			"ETag":                 {"\"D41D8CD98F00B204E9800998ECF8****\""},
			"x-oss-hash-crc64ecma": {"316181249502703****"},
			"Content-MD5":          {"1B2M2Y8AsgTpgAmY7PhC****"},
		},
		Body: io.NopCloser(strings.NewReader(body)),
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyCreateSelectObjectMeta)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C06A3B67B8B5A3DA422****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	assert.Equal(t, result.TotalScanned, int64(162039))
	assert.Equal(t, result.MetaStatus, 200)
	assert.Equal(t, result.SplitsCount, int32(1))
	assert.Equal(t, result.RowsCount, int64(100))
	assert.Equal(t, result.ErrorMsg, ".")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidInputRecordDelimiter</Code>
  <Message>Invalid RecordDelimiter parameter:CnNhZGFzZA==</Message>
  <RequestId>65698AC54B39ED97D780****</RequestId>
  <HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
  <EC>0016-00000829</EC>
</Error>`
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidInputRecordDelimiter",
		Headers: http.Header{
			"X-Oss-Request-Id": {"65698AC54B39ED97D780****"},
			"Content-Type":     {"application/xml"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyCreateSelectObjectMeta)
	assert.NotNil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidInputRecordDelimiter")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "65698AC54B39ED97D780****")
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
			"X-Oss-Request-Id": {"568D5566F2D0F89F5C0E****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyCreateSelectObjectMeta)
	assert.NotNil(t, err)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "568D5566F2D0F89F5C0E****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_SelectObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *SelectObjectRequest
	var input *OperationInput
	var err error

	request = &SelectObjectRequest{}
	input = &OperationInput{
		OpName: "SelectObject",
		Method: "POST",
		Bucket: request.Bucket,
		Key:    request.Key,
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
	}
	err = c.marshalInput(request, input, marshalSelectObjectRequest, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &SelectObjectRequest{
		Bucket: Ptr("oss-demo"),
	}
	err = c.marshalInput(request, input, marshalSelectObjectRequest, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &SelectObjectRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
	}
	err = c.marshalInput(request, input, marshalSelectObjectRequest, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &SelectObjectRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
	}
	err = c.marshalInput(request, input, marshalSelectObjectRequest, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &SelectObjectRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
		SelectRequest: &SelectRequest{
			Expression: Ptr("select Year, StateAbbr, CityName, Short_Question_Text from ossobject where Measure like '%blood pressure%Years'"),
			InputSerializationSelect: InputSerializationSelect{
				CsvBodyInput: &CSVSelectInput{
					FileHeaderInfo: Ptr("Use"),
				},
			},
		},
	}
	err = c.marshalInput(request, input, marshalSelectObjectRequest, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-process"], "csv/select")
	data, _ := io.ReadAll(input.Body)
	assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IFllYXIsIFN0YXRlQWJiciwgQ2l0eU5hbWUsIFNob3J0X1F1ZXN0aW9uX1RleHQgZnJvbSBvc3NvYmplY3Qgd2hlcmUgTWVhc3VyZSBsaWtlICclYmxvb2QgcHJlc3N1cmUlWWVhcnMn</Expression><InputSerialization><CSV><FileHeaderInfo>Use</FileHeaderInfo></CSV></InputSerialization><OutputSerialization></OutputSerialization></SelectRequest>", string(data))

	request = &SelectObjectRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
		SelectRequest: &SelectRequest{
			Expression: Ptr("select * from ossobject'"),
			InputSerializationSelect: InputSerializationSelect{
				CsvBodyInput: &CSVSelectInput{
					FileHeaderInfo:   Ptr("None"),
					CommentCharacter: Ptr("#"),
					RecordDelimiter:  Ptr("\n"),
					FieldDelimiter:   Ptr(","),
					QuoteCharacter:   Ptr("\""),
					Range:            Ptr(""),
				},
			},
		},
	}
	err = c.marshalInput(request, input, marshalSelectObjectRequest, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-process"], "csv/select")
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0ICogZnJvbSBvc3NvYmplY3Qn</Expression><InputSerialization><CSV><FileHeaderInfo>None</FileHeaderInfo><RecordDelimiter>Cg==</RecordDelimiter><FieldDelimiter>LA==</FieldDelimiter><QuoteCharacter>Ig==</QuoteCharacter><CommentCharacter>Iw==</CommentCharacter><Range></Range></CSV></InputSerialization><OutputSerialization></OutputSerialization></SelectRequest>", string(data))

	request = &SelectObjectRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
		SelectRequest: &SelectRequest{
			Expression: Ptr("select Year,StateAbbr, CityName, Short_Question_Text from ossobject"),
			InputSerializationSelect: InputSerializationSelect{
				CsvBodyInput: &CSVSelectInput{
					FileHeaderInfo: Ptr("Use"),
					Range:          Ptr("0-2"),
				},
			},
		},
	}
	err = c.marshalInput(request, input, marshalSelectObjectRequest, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-process"], "csv/select")
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IFllYXIsU3RhdGVBYmJyLCBDaXR5TmFtZSwgU2hvcnRfUXVlc3Rpb25fVGV4dCBmcm9tIG9zc29iamVjdA==</Expression><InputSerialization><CSV><FileHeaderInfo>Use</FileHeaderInfo><Range>line-range=0-2</Range></CSV></InputSerialization><OutputSerialization></OutputSerialization></SelectRequest>", string(data))

	request = &SelectObjectRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
		SelectRequest: &SelectRequest{
			Expression: Ptr("select avg(cast(year as int)), max(cast(year as int)), min(cast(year as int)) from ossobject where year = 2015"),
			InputSerializationSelect: InputSerializationSelect{
				CsvBodyInput: &CSVSelectInput{
					FileHeaderInfo: Ptr("Use"),
				},
			},
		},
	}
	err = c.marshalInput(request, input, marshalSelectObjectRequest, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-process"], "csv/select")
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IGF2ZyhjYXN0KHllYXIgYXMgaW50KSksIG1heChjYXN0KHllYXIgYXMgaW50KSksIG1pbihjYXN0KHllYXIgYXMgaW50KSkgZnJvbSBvc3NvYmplY3Qgd2hlcmUgeWVhciA9IDIwMTU=</Expression><InputSerialization><CSV><FileHeaderInfo>Use</FileHeaderInfo></CSV></InputSerialization><OutputSerialization></OutputSerialization></SelectRequest>", string(data))

	request = &SelectObjectRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
		SelectRequest: &SelectRequest{
			Expression: Ptr("select avg(cast(data_value as double)), max(cast(data_value as double)), sum(cast(data_value as double)) from ossobject"),
			InputSerializationSelect: InputSerializationSelect{
				CsvBodyInput: &CSVSelectInput{
					FileHeaderInfo: Ptr("Use"),
				},
			},
		},
	}
	err = c.marshalInput(request, input, marshalSelectObjectRequest, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-process"], "csv/select")
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IGF2ZyhjYXN0KGRhdGFfdmFsdWUgYXMgZG91YmxlKSksIG1heChjYXN0KGRhdGFfdmFsdWUgYXMgZG91YmxlKSksIHN1bShjYXN0KGRhdGFfdmFsdWUgYXMgZG91YmxlKSkgZnJvbSBvc3NvYmplY3Q=</Expression><InputSerialization><CSV><FileHeaderInfo>Use</FileHeaderInfo></CSV></InputSerialization><OutputSerialization></OutputSerialization></SelectRequest>", string(data))

	request = &SelectObjectRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
		SelectRequest: &SelectRequest{
			Expression: Ptr("select _1, _2 from ossobject"),
			OutputSerializationSelect: OutputSerializationSelect{
				CsvBodyOutput: &CSVSelectOutput{
					RecordDelimiter: Ptr("\r\n"),
					FieldDelimiter:  Ptr("|"),
				},
			},
		},
	}
	err = c.marshalInput(request, input, marshalSelectObjectRequest, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-process"], "csv/select")
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IF8xLCBfMiBmcm9tIG9zc29iamVjdA==</Expression><InputSerialization></InputSerialization><OutputSerialization><CSV><RecordDelimiter>&#xD;&#xA;</RecordDelimiter><FieldDelimiter>|</FieldDelimiter></CSV></OutputSerialization></SelectRequest>", string(data))

	request = &SelectObjectRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
		SelectRequest: &SelectRequest{
			Expression: Ptr("select * from ossobject"),
			OutputSerializationSelect: OutputSerializationSelect{
				EnablePayloadCrc: Ptr(true),
			},
		},
	}
	err = c.marshalInput(request, input, marshalSelectObjectRequest, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-process"], "csv/select")
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0ICogZnJvbSBvc3NvYmplY3Q=</Expression><InputSerialization></InputSerialization><OutputSerialization><EnablePayloadCrc>true</EnablePayloadCrc></OutputSerialization></SelectRequest>", string(data))

	request = &SelectObjectRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
		SelectRequest: &SelectRequest{
			Expression: Ptr("select _1, _2 from ossobject"),
			SelectOptions: &SelectOptions{
				SkipPartialDataRecord: Ptr(true),
			},
		},
	}
	err = c.marshalInput(request, input, marshalSelectObjectRequest, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-process"], "csv/select")
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IF8xLCBfMiBmcm9tIG9zc29iamVjdA==</Expression><InputSerialization></InputSerialization><OutputSerialization></OutputSerialization><Options><SkipPartialDataRecord>true</SkipPartialDataRecord></Options></SelectRequest>", string(data))

	request = &SelectObjectRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
		SelectRequest: &SelectRequest{
			Expression: Ptr("select _1,from ossobject"),
			OutputSerializationSelect: OutputSerializationSelect{
				OutputRawData: Ptr(true),
			},
		},
	}

	err = c.marshalInput(request, input, marshalSelectObjectRequest, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-process"], "csv/select")
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IF8xLGZyb20gb3Nzb2JqZWN0</Expression><InputSerialization></InputSerialization><OutputSerialization><OutputRawData>true</OutputRawData></OutputSerialization></SelectRequest>", string(data))

	request = &SelectObjectRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
		SelectRequest: &SelectRequest{
			Expression: Ptr("select _1,from ossobject"),
			OutputSerializationSelect: OutputSerializationSelect{
				OutputRawData:  Ptr(true),
				KeepAllColumns: Ptr(true),
			},
		},
	}
	err = c.marshalInput(request, input, marshalSelectObjectRequest, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-process"], "csv/select")
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IF8xLGZyb20gb3Nzb2JqZWN0</Expression><InputSerialization></InputSerialization><OutputSerialization><OutputRawData>true</OutputRawData><KeepAllColumns>true</KeepAllColumns></OutputSerialization></SelectRequest>", string(data))

	request = &SelectObjectRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
		SelectRequest: &SelectRequest{
			Expression: Ptr("select _1,from ossobject"),
			OutputSerializationSelect: OutputSerializationSelect{
				OutputHeader: Ptr(true),
			},
			InputSerializationSelect: InputSerializationSelect{
				CsvBodyInput: &CSVSelectInput{
					FileHeaderInfo: Ptr("Use"),
				},
			},
		},
	}
	err = c.marshalInput(request, input, marshalSelectObjectRequest, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-process"], "csv/select")
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IF8xLGZyb20gb3Nzb2JqZWN0</Expression><InputSerialization><CSV><FileHeaderInfo>Use</FileHeaderInfo></CSV></InputSerialization><OutputSerialization><OutputHeader>true</OutputHeader></OutputSerialization></SelectRequest>", string(data))

	request = &SelectObjectRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
		SelectRequest: &SelectRequest{
			Expression: Ptr("select name from ossobject"),
			OutputSerializationSelect: OutputSerializationSelect{
				OutputHeader:     Ptr(true),
				EnablePayloadCrc: Ptr(true),
			},
			InputSerializationSelect: InputSerializationSelect{
				CsvBodyInput: &CSVSelectInput{
					FileHeaderInfo: Ptr("Use"),
				},
			},
		},
	}
	err = c.marshalInput(request, input, marshalSelectObjectRequest, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-process"], "csv/select")
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IG5hbWUgZnJvbSBvc3NvYmplY3Q=</Expression><InputSerialization><CSV><FileHeaderInfo>Use</FileHeaderInfo></CSV></InputSerialization><OutputSerialization><EnablePayloadCrc>true</EnablePayloadCrc><OutputHeader>true</OutputHeader></OutputSerialization></SelectRequest>", string(data))

	request = &SelectObjectRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
		SelectRequest: &SelectRequest{
			Expression: Ptr("select person.firstname as aaa as firstname, person.lastname, extra from ossobject'"),
			InputSerializationSelect: InputSerializationSelect{
				CompressionType: Ptr("GZIP"),
				CsvBodyInput: &CSVSelectInput{
					FileHeaderInfo:             Ptr("IGNORE"),
					RecordDelimiter:            Ptr("\n"),
					FieldDelimiter:             Ptr(","),
					QuoteCharacter:             Ptr("\""),
					CommentCharacter:           Ptr("#"),
					SplitRange:                 Ptr("10-12"),
					AllowQuotedRecordDelimiter: Ptr(true),
				},
			},
			OutputSerializationSelect: OutputSerializationSelect{
				CsvBodyOutput: &CSVSelectOutput{
					RecordDelimiter: Ptr("\n"),
					FieldDelimiter:  Ptr(","),
				},
				KeepAllColumns:   Ptr(false),
				OutputRawData:    Ptr(true),
				EnablePayloadCrc: Ptr(true),
				OutputHeader:     Ptr(false),
			},
			SelectOptions: &SelectOptions{
				SkipPartialDataRecord:    Ptr(false),
				MaxSkippedRecordsAllowed: Ptr(2),
			},
		},
	}
	err = c.marshalInput(request, input, marshalSelectObjectRequest, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-process"], "csv/select")
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IHBlcnNvbi5maXJzdG5hbWUgYXMgYWFhIGFzIGZpcnN0bmFtZSwgcGVyc29uLmxhc3RuYW1lLCBleHRyYSBmcm9tIG9zc29iamVjdCc=</Expression><InputSerialization><CSV><FileHeaderInfo>IGNORE</FileHeaderInfo><RecordDelimiter>Cg==</RecordDelimiter><FieldDelimiter>LA==</FieldDelimiter><QuoteCharacter>Ig==</QuoteCharacter><CommentCharacter>Iw==</CommentCharacter><Range>split-range=10-12</Range><AllowQuotedRecordDelimiter>true</AllowQuotedRecordDelimiter></CSV><CompressionType>GZIP</CompressionType></InputSerialization><OutputSerialization><CSV><RecordDelimiter>&#xA;</RecordDelimiter><FieldDelimiter>,</FieldDelimiter></CSV><OutputRawData>true</OutputRawData><KeepAllColumns>false</KeepAllColumns><EnablePayloadCrc>true</EnablePayloadCrc><OutputHeader>false</OutputHeader></OutputSerialization><Options><SkipPartialDataRecord>false</SkipPartialDataRecord><MaxSkippedRecordsAllowed>2</MaxSkippedRecordsAllowed></Options></SelectRequest>", string(data))

	request = &SelectObjectRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
		SelectRequest: &SelectRequest{
			Expression: Ptr("select * from ossobject.objects[*] where party = 'Democrat'"),
			InputSerializationSelect: InputSerializationSelect{
				JsonBodyInput: &JSONSelectInput{
					JSONType: Ptr("DOCUMENT"),
				},
			},
			OutputSerializationSelect: OutputSerializationSelect{
				JsonBodyOutput: &JSONSelectOutput{
					RecordDelimiter: Ptr(","),
				},
			},
		},
	}
	err = c.marshalInput(request, input, marshalSelectObjectRequest, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-process"], "json/select")
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0ICogZnJvbSBvc3NvYmplY3Qub2JqZWN0c1sqXSB3aGVyZSBwYXJ0eSA9ICdEZW1vY3JhdCc=</Expression><InputSerialization><JSON><Type>DOCUMENT</Type></JSON></InputSerialization><OutputSerialization><JSON><RecordDelimiter>LA==</RecordDelimiter></JSON></OutputSerialization></SelectRequest>", string(data))

	request = &SelectObjectRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
		SelectRequest: &SelectRequest{
			Expression: Ptr("select * from ossobject where party = 'Democrat'"),
			InputSerializationSelect: InputSerializationSelect{
				JsonBodyInput: &JSONSelectInput{
					JSONType: Ptr("LINES"),
				},
			},
			OutputSerializationSelect: OutputSerializationSelect{
				JsonBodyOutput: &JSONSelectOutput{
					RecordDelimiter: Ptr(","),
				},
			},
		},
	}
	err = c.marshalInput(request, input, marshalSelectObjectRequest, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-process"], "json/select")
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0ICogZnJvbSBvc3NvYmplY3Qgd2hlcmUgcGFydHkgPSAnRGVtb2NyYXQn</Expression><InputSerialization><JSON><Type>LINES</Type></JSON></InputSerialization><OutputSerialization><JSON><RecordDelimiter>LA==</RecordDelimiter></JSON></OutputSerialization></SelectRequest>", string(data))

	request = &SelectObjectRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
		SelectRequest: &SelectRequest{
			Expression: Ptr("select person.firstname as aaa as firstname, person.lastname, extra from ossobject'"),
			InputSerializationSelect: InputSerializationSelect{
				JsonBodyInput: &JSONSelectInput{
					JSONType: Ptr("LINES"),
					Range:    Ptr("0-1"),
				},
			},
			OutputSerializationSelect: OutputSerializationSelect{
				JsonBodyOutput: &JSONSelectOutput{
					RecordDelimiter: Ptr(","),
				},
			},
		},
	}
	err = c.marshalInput(request, input, marshalSelectObjectRequest, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-process"], "json/select")
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IHBlcnNvbi5maXJzdG5hbWUgYXMgYWFhIGFzIGZpcnN0bmFtZSwgcGVyc29uLmxhc3RuYW1lLCBleHRyYSBmcm9tIG9zc29iamVjdCc=</Expression><InputSerialization><JSON><Type>LINES</Type><Range>line-range=0-1</Range></JSON></InputSerialization><OutputSerialization><JSON><RecordDelimiter>LA==</RecordDelimiter></JSON></OutputSerialization></SelectRequest>", string(data))

	request = &SelectObjectRequest{
		Bucket: Ptr("oss-demo"),
		Key:    Ptr("oss-key"),
		SelectRequest: &SelectRequest{
			Expression: Ptr("select person.firstname as aaa as firstname, person.lastname, extra from ossobject'"),
			InputSerializationSelect: InputSerializationSelect{
				CompressionType: Ptr("GZIP"),
				JsonBodyInput: &JSONSelectInput{
					JSONType:                Ptr("LINES"),
					Range:                   Ptr("0-1"),
					ParseJSONNumberAsString: Ptr(true),
					SplitRange:              Ptr("10-12"),
				},
			},
			OutputSerializationSelect: OutputSerializationSelect{
				JsonBodyOutput: &JSONSelectOutput{
					RecordDelimiter: Ptr(","),
				},
				OutputRawData:    Ptr(true),
				EnablePayloadCrc: Ptr(true),
			},
			SelectOptions: &SelectOptions{
				SkipPartialDataRecord:    Ptr(false),
				MaxSkippedRecordsAllowed: Ptr(2),
			},
		},
	}
	err = c.marshalInput(request, input, marshalSelectObjectRequest, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["x-oss-process"], "json/select")
	data, _ = io.ReadAll(input.Body)
	assert.Equal(t, "<SelectRequest><Expression>c2VsZWN0IHBlcnNvbi5maXJzdG5hbWUgYXMgYWFhIGFzIGZpcnN0bmFtZSwgcGVyc29uLmxhc3RuYW1lLCBleHRyYSBmcm9tIG9zc29iamVjdCc=</Expression><InputSerialization><JSON><Type>LINES</Type><Range>split-range=10-12</Range><ParseJsonNumberAsString>true</ParseJsonNumberAsString></JSON><CompressionType>GZIP</CompressionType></InputSerialization><OutputSerialization><JSON><RecordDelimiter>LA==</RecordDelimiter></JSON><OutputRawData>true</OutputRawData><EnablePayloadCrc>true</EnablePayloadCrc></OutputSerialization><Options><SkipPartialDataRecord>false</SkipPartialDataRecord><MaxSkippedRecordsAllowed>2</MaxSkippedRecordsAllowed></Options></SelectRequest>", string(data))
}

func TestUnmarshalOutput_SelectObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	hexString := "01800001000016cc000000000000000000012dfb323031352c414c2c4269726d696e6768616d2c3231323233370a323031352c414c2c4269726d696e6768616d2c3231323233370a323031352c414c2c4269726d696e6768616d2c333034320a323031352c414c2c4269726d696e6768616d2c323733350a323031352c414c2c4269726d696e6768616d2c333333380a323031352c414c2c4269726d696e6768616d2c323836340a323031352c414c2c4269726d696e6768616d2c323537370a323031352c414c2c4269726d696e6768616d2c333835390a323031352c414c2c4269726d696e6768616d2c353335340a323031352c414c2c4269726d696e6768616d2c323837360a323031352c414c2c4269726d696e6768616d2c323138310a323031352c414c2c4269726d696e6768616d2c333138390a323031352c414c2c4269726d696e6768616d2c333339300a323031352c414c2c4269726d696e6768616d2c313839340a323031352c414c2c4269726d696e6768616d2c333838350a323031352c414c2c4269726d696e6768616d2c333138360a323031352c414c2c4269726d696e6768616d2c323633300a323031352c414c2c4269726d696e6768616d2c323933360a323031352c414c2c4269726d696e6768616d2c323935320a323031352c414c2c4269726d696e6768616d2c333235370a323031352c414c2c4269726d696e6768616d2c333632390a323031352c414c2c4269726d696e6768616d2c333939320a323031352c414c2c4269726d696e6768616d2c323036340a323031352c414c2c4269726d696e6768616d2c333737390a323031352c414c2c4269726d696e6768616d2c323230330a323031352c414c2c4269726d696e6768616d2c333633370a323031352c414c2c4269726d696e6768616d2c3933310a323031352c414c2c4269726d696e6768616d2c3934370a323031352c414c2c4269726d696e6768616d2c323437370a323031352c414c2c4269726d696e6768616d2c323738300a323031352c414c2c4269726d696e6768616d2c343638330a323031352c414c2c4269726d696e6768616d2c353036330a323031352c414c2c4269726d696e6768616d2c353430390a323031352c414c2c4269726d696e6768616d2c343139390a323031352c414c2c4269726d696e6768616d2c313738330a323031352c414c2c4269726d696e6768616d2c333737320a323031352c414c2c4269726d696e6768616d2c323334310a323031352c414c2c4269726d696e6768616d2c353030330a323031352c414c2c4269726d696e6768616d2c333438300a323031352c414c2c4269726d696e6768616d2c323934340a323031352c414c2c4269726d696e6768616d2c313836310a323031352c414c2c4269726d696e6768616d2c313136370a323031352c414c2c4269726d696e6768616d2c333134360a323031352c414c2c4269726d696e6768616d2c333438320a323031352c414c2c4269726d696e6768616d2c313530370a323031352c414c2c4269726d696e6768616d2c323538370a323031352c414c2c4269726d696e6768616d2c323838310a323031352c414c2c4269726d696e6768616d2c333734300a323031352c414c2c4269726d696e6768616d2c333436330a323031352c414c2c4269726d696e6768616d2c313832340a323031352c414c2c4269726d696e6768616d2c343336370a323031352c414c2c4269726d696e6768616d2c323337320a323031352c414c2c4269726d696e6768616d2c333431330a323031352c414c2c4269726d696e6768616d2c343231360a323031352c414c2c4269726d696e6768616d2c343933330a323031352c414c2c4269726d696e6768616d2c353033390a323031352c414c2c4269726d696e6768616d2c313937350a323031352c414c2c4269726d696e6768616d2c313632310a323031352c414c2c4269726d696e6768616d2c323532340a323031352c414c2c4269726d696e6768616d2c343631320a323031352c414c2c4269726d696e6768616d2c3131340a323031352c414c2c4269726d696e6768616d2c37340a323031352c414c2c4269726d696e6768616d2c313532380a323031352c414c2c4269726d696e6768616d2c3136380a323031352c414c2c4269726d696e6768616d2c3137320a323031352c414c2c4269726d696e6768616d2c3531340a323031352c414c2c4269726d696e6768616d2c38360a323031352c414c2c4269726d696e6768616d2c313638380a323031352c414c2c4269726d696e6768616d2c34320a323031352c414c2c4269726d696e6768616d2c390a323031352c414c2c4269726d696e6768616d2c3831350a323031352c414c2c4269726d696e6768616d2c313036320a323031352c414c2c4269726d696e6768616d2c313338350a323031352c414c2c4269726d696e6768616d2c3932380a323031352c414c2c4269726d696e6768616d2c313135370a323031352c414c2c4269726d696e6768616d2c360a323031352c414c2c4269726d696e6768616d2c313931350a323031352c414c2c4269726d696e6768616d2c3330340a323031352c414c2c4269726d696e6768616d2c34340a323031352c414c2c4269726d696e6768616d2c32330a323031352c414c2c4269726d696e6768616d2c3134340a323031352c414c2c4269726d696e6768616d2c3430330a323031352c414c2c4269726d696e6768616d2c313036360a323031352c414c2c4269726d696e6768616d2c3431380a323031352c414c2c4269726d696e6768616d2c3431300a323031352c414c2c4269726d696e6768616d2c3337310a323031352c414c2c4269726d696e6768616d2c34340a323031352c414c2c4269726d696e6768616d2c3439380a323031352c414c2c4269726d696e6768616d2c3131330a323031352c414c2c4269726d696e6768616d2c313236310a323031352c414c2c4269726d696e6768616d2c390a323031352c414c2c4269726d696e6768616d2c313531340a323031352c414c2c4269726d696e6768616d2c343432340a323031352c414c2c4269726d696e6768616d2c313738320a323031352c414c2c4269726d696e6768616d2c3935320a323031352c414c2c4269726d696e6768616d2c323737380a323031352c414c2c4269726d696e6768616d2c3339370a323031352c414c2c4269726d696e6768616d2c3634340a323031352c414c2c4269726d696e6768616d2c31360a323031352c414c2c4269726d696e6768616d2c3936380a323031352c414c2c4269726d696e6768616d2c3231323233370a323031352c414c2c4269726d696e6768616d2c3231323233370a323031352c414c2c4269726d696e6768616d2c333034320a323031352c414c2c4269726d696e6768616d2c323733350a323031352c414c2c4269726d696e6768616d2c333333380a323031352c414c2c4269726d696e6768616d2c323836340a323031352c414c2c4269726d696e6768616d2c323537370a323031352c414c2c4269726d696e6768616d2c333835390a323031352c414c2c4269726d696e6768616d2c353335340a323031352c414c2c4269726d696e6768616d2c323837360a323031352c414c2c4269726d696e6768616d2c323138310a323031352c414c2c4269726d696e6768616d2c333138390a323031352c414c2c4269726d696e6768616d2c333339300a323031352c414c2c4269726d696e6768616d2c313839340a323031352c414c2c4269726d696e6768616d2c333838350a323031352c414c2c4269726d696e6768616d2c333138360a323031352c414c2c4269726d696e6768616d2c323633300a323031352c414c2c4269726d696e6768616d2c323933360a323031352c414c2c4269726d696e6768616d2c323935320a323031352c414c2c4269726d696e6768616d2c333235370a323031352c414c2c4269726d696e6768616d2c333632390a323031352c414c2c4269726d696e6768616d2c333939320a323031352c414c2c4269726d696e6768616d2c323036340a323031352c414c2c4269726d696e6768616d2c333737390a323031352c414c2c4269726d696e6768616d2c323230330a323031352c414c2c4269726d696e6768616d2c333633370a323031352c414c2c4269726d696e6768616d2c3933310a323031352c414c2c4269726d696e6768616d2c3934370a323031352c414c2c4269726d696e6768616d2c323437370a323031352c414c2c4269726d696e6768616d2c323738300a323031352c414c2c4269726d696e6768616d2c343638330a323031352c414c2c4269726d696e6768616d2c353036330a323031352c414c2c4269726d696e6768616d2c353430390a323031352c414c2c4269726d696e6768616d2c343139390a323031352c414c2c4269726d696e6768616d2c313738330a323031352c414c2c4269726d696e6768616d2c333737320a323031352c414c2c4269726d696e6768616d2c323334310a323031352c414c2c4269726d696e6768616d2c353030330a323031352c414c2c4269726d696e6768616d2c333438300a323031352c414c2c4269726d696e6768616d2c323934340a323031352c414c2c4269726d696e6768616d2c313836310a323031352c414c2c4269726d696e6768616d2c313136370a323031352c414c2c4269726d696e6768616d2c333134360a323031352c414c2c4269726d696e6768616d2c333438320a323031352c414c2c4269726d696e6768616d2c313530370a323031352c414c2c4269726d696e6768616d2c323538370a323031352c414c2c4269726d696e6768616d2c323838310a323031352c414c2c4269726d696e6768616d2c333734300a323031352c414c2c4269726d696e6768616d2c333436330a323031352c414c2c4269726d696e6768616d2c313832340a323031352c414c2c4269726d696e6768616d2c343336370a323031352c414c2c4269726d696e6768616d2c323337320a323031352c414c2c4269726d696e6768616d2c333431330a323031352c414c2c4269726d696e6768616d2c343231360a323031352c414c2c4269726d696e6768616d2c343933330a323031352c414c2c4269726d696e6768616d2c353033390a323031352c414c2c4269726d696e6768616d2c313937350a323031352c414c2c4269726d696e6768616d2c313632310a323031352c414c2c4269726d696e6768616d2c323532340a323031352c414c2c4269726d696e6768616d2c343631320a323031352c414c2c4269726d696e6768616d2c3131340a323031352c414c2c4269726d696e6768616d2c37340a323031352c414c2c4269726d696e6768616d2c313532380a323031352c414c2c4269726d696e6768616d2c3136380a323031352c414c2c4269726d696e6768616d2c3137320a323031352c414c2c4269726d696e6768616d2c3531340a323031352c414c2c4269726d696e6768616d2c38360a323031352c414c2c4269726d696e6768616d2c313638380a323031352c414c2c4269726d696e6768616d2c34320a323031352c414c2c4269726d696e6768616d2c390a323031352c414c2c4269726d696e6768616d2c3831350a323031352c414c2c4269726d696e6768616d2c313036320a323031352c414c2c4269726d696e6768616d2c313338350a323031352c414c2c4269726d696e6768616d2c3932380a323031352c414c2c4269726d696e6768616d2c313135370a323031352c414c2c4269726d696e6768616d2c360a323031352c414c2c4269726d696e6768616d2c313931350a323031352c414c2c4269726d696e6768616d2c3330340a323031352c414c2c4269726d696e6768616d2c34340a323031352c414c2c4269726d696e6768616d2c32330a323031352c414c2c4269726d696e6768616d2c3134340a323031352c414c2c4269726d696e6768616d2c3430330a323031352c414c2c4269726d696e6768616d2c313036360a323031352c414c2c4269726d696e6768616d2c3431380a323031352c414c2c4269726d696e6768616d2c3431300a323031352c414c2c4269726d696e6768616d2c3337310a323031352c414c2c4269726d696e6768616d2c34340a323031352c414c2c4269726d696e6768616d2c3439380a323031352c414c2c4269726d696e6768616d2c3131330a323031352c414c2c4269726d696e6768616d2c313236310a323031352c414c2c4269726d696e6768616d2c390a323031352c414c2c4269726d696e6768616d2c313531340a323031352c414c2c4269726d696e6768616d2c343432340a323031352c414c2c4269726d696e6768616d2c313738320a323031352c414c2c4269726d696e6768616d2c3935320a323031352c414c2c4269726d696e6768616d2c323737380a323031352c414c2c4269726d696e6768616d2c3339370a323031352c414c2c4269726d696e6768616d2c3634340a323031352c414c2c4269726d696e6768616d2c31360a323031352c414c2c4269726d696e6768616d2c3936380a323031352c414c2c4269726d696e6768616d2c3231323233370a323031352c414c2c4269726d696e6768616d2c3231323233370a323031352c414c2c4269726d696e6768616d2c333034320a323031352c414c2c4269726d696e6768616d2c323733350a323031352c414c2c4269726d696e6768616d2c333333380a323031352c414c2c4269726d696e6768616d2c323836340a323031352c414c2c4269726d696e6768616d2c323537370a323031352c414c2c4269726d696e6768616d2c333835390a323031352c414c2c4269726d696e6768616d2c353335340a323031352c414c2c4269726d696e6768616d2c323837360a323031352c414c2c4269726d696e6768616d2c323138310a323031352c414c2c4269726d696e6768616d2c333138390a323031352c414c2c4269726d696e6768616d2c333339300a323031352c414c2c4269726d696e6768616d2c313839340a323031352c414c2c4269726d696e6768616d2c333838350a323031352c414c2c4269726d696e6768616d2c333138360a323031352c414c2c4269726d696e6768616d2c323633300a323031352c414c2c4269726d696e6768616d2c323933360a323031352c414c2c4269726d696e6768616d2c323935320a323031352c414c2c4269726d696e6768616d2c333235370a323031352c414c2c4269726d696e6768616d2c333632390a323031352c414c2c4269726d696e6768616d2c333939320a323031352c414c2c4269726d696e6768616d2c323036340a323031352c414c2c4269726d696e6768616d2c333737390a323031352c414c2c4269726d696e6768616d2c323230330a323031352c414c2c4269726d696e6768616d2c333633370a323031352c414c2c4269726d696e6768616d2c3933310a323031352c414c2c4269726d696e6768616d2c3934370a323031352c414c2c4269726d696e6768616d2c323935320a323031352c414c2c4269726d696e6768616d2c323437370a323031352c414c2c4269726d696e6768616d2c323738300a323031352c414c2c4269726d696e6768616d2c343638330a323031352c414c2c4269726d696e6768616d2c353036330a323031352c414c2c4269726d696e6768616d2c353430390a323031352c414c2c4269726d696e6768616d2c343139390a323031352c414c2c4269726d696e6768616d2c313738330a323031352c414c2c4269726d696e6768616d2c333737320a323031352c414c2c4269726d696e6768616d2c323334310a323031352c414c2c4269726d696e6768616d2c353030330a323031352c414c2c4269726d696e6768616d2c333438300a323031352c414c2c4269726d696e6768616d2c323934340a323031352c414c2c4269726d696e6768616d2c313836310a323031352c414c2c4269726d696e6768616d2c313136370a323031352c414c2c4269726d696e6768616d2c333134360a323031352c414c2c4269726d696e6768616d2c333438320a323031352c414c2c4269726d696e6768616d2c313530370a000000000180000100000020000000000000000000012efc323031352c414c2c4269726d696e6768616d2c323538370a000000000180000500000014000000000000000000012efc0000000000012efc000000c800000000"
	body := hexStrToByte(hexString)
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":     {"5C06A3B67B8B5A3DA422****"},
			"Date":                 {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"Content-Type":         {"application/vnd.ms-excel"},
			"ETag":                 {"\"D41D8CD98F00B204E9800998ECF8****\""},
			"x-oss-hash-crc64ecma": {"316181249502703****"},
			"Content-MD5":          {"1B2M2Y8AsgTpgAmY7PhC****"},
		},
		Body: io.NopCloser(strings.NewReader(string(body))),
	}
	result := &SelectObjectResult{}
	readerWrapper := &ReaderWrapper{
		Body:                output.Body,
		WriterForCheckCrc32: crc32.NewIEEE(),
	}
	readerWrapper.OutputRawData = strings.ToUpper(output.Headers.Get("x-oss-select-output-raw")) == "TRUE"
	result.Body = readerWrapper
	assert.Nil(t, err)
	err = c.unmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C06A3B67B8B5A3DA422****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	dataByte, err := io.ReadAll(result.Body)
	assert.Equal(t, string(dataByte[:25]), "2015,AL,Birmingham,212237")

	hexString = "018000010000014d000000000000000000000145596561722c5374617465416262722c5374617465446573632c436974794e616d652c47656f677261706869634c6576656c2c44617461536f757263652c43617465676f72792c556e6971756549442c4d6561737572652c446174615f56616c75655f556e69742c4461746156616c75655479706549442c446174615f56616c75655f547970652c446174615f56616c75652c4c6f775f436f6e666964656e63655f4c696d69742c486967685f436f6e666964656e63655f4c696d69742c446174615f56616c75655f466f6f746e6f74655f53796d626f6c2c446174615f56616c75655f466f6f746e6f74652c506f70756c6174696f6e436f756e742c47656f4c6f636174696f6e2c43617465676f727949442c4d65617375726549642c43697479464950532c5472616374464950532c53686f72745f5175657374696f6e5f546578740d0a0000000001800005000000140000000000000000000001450000000000000145000000c800000000"
	body = hexStrToByte(hexString)
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":     {"5C06A3B67B8B5A3DA422****"},
			"Date":                 {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"Content-Type":         {"application/vnd.ms-excel"},
			"ETag":                 {"\"D41D8CD98F00B204E9800998ECF8****\""},
			"x-oss-hash-crc64ecma": {"316181249502703****"},
			"Content-MD5":          {"1B2M2Y8AsgTpgAmY7PhC****"},
		},
		Body: io.NopCloser(strings.NewReader(string(body))),
	}
	readerWrapper = &ReaderWrapper{
		Body:                output.Body,
		WriterForCheckCrc32: crc32.NewIEEE(),
	}
	readerWrapper.OutputRawData = strings.ToUpper(output.Headers.Get("x-oss-select-output-raw")) == "TRUE"
	result.Body = readerWrapper
	assert.Nil(t, err)
	err = c.unmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C06A3B67B8B5A3DA422****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	dataByte, err = io.ReadAll(result.Body)
	assert.Equal(t, string(dataByte[:25]), "Year,StateAbbr,StateDesc,")

	hexString = "018000010000003e0000000000000000000002e1323031352c55532c2c4865616c746820496e737572616e63650d0a323031352c55532c2c4865616c746820496e737572616e63650d0a0000000001800005000000140000000000000000000002e100000000000002e1000000c800000000"
	body = hexStrToByte(hexString)
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":     {"5C06A3B67B8B5A3DA422****"},
			"Date":                 {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"Content-Type":         {"application/vnd.ms-excel"},
			"ETag":                 {"\"D41D8CD98F00B204E9800998ECF8****\""},
			"x-oss-hash-crc64ecma": {"316181249502703****"},
			"Content-MD5":          {"1B2M2Y8AsgTpgAmY7PhC****"},
		},
		Body: io.NopCloser(strings.NewReader(string(body))),
	}
	readerWrapper = &ReaderWrapper{
		Body:                output.Body,
		WriterForCheckCrc32: crc32.NewIEEE(),
	}
	readerWrapper.OutputRawData = strings.ToUpper(output.Headers.Get("x-oss-select-output-raw")) == "TRUE"
	result.Body = readerWrapper
	assert.Nil(t, err)
	err = c.unmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C06A3B67B8B5A3DA422****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	dataByte, err = io.ReadAll(result.Body)
	assert.Equal(t, string(dataByte[:25]), "2015,US,,Health Insurance")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidRange</Code>
  <Message>The split range is invalid. The non-negative value should be from 0 to total_splits - 1. Its semantics is same as the value in Http Range header.Actual split range is 90-109.</Message>
  <RequestId>65698FC48B404DFCB581CC8F</RequestId>
  <HostId>demo-walker-6961.oss-cn-hangzhou.aliyuncs.com</HostId>
  <EC>0016-00000817</EC>
  <RecommendDoc>https://api.aliyun.com/troubleshoot?q=0016-00000817</RecommendDoc>
</Error>`
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidRange",
		Headers: http.Header{
			"X-Oss-Request-Id": {"65698AC54B39ED97D780****"},
			"Content-Type":     {"application/xml"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	readerWrapper = &ReaderWrapper{
		Body:                output.Body,
		WriterForCheckCrc32: crc32.NewIEEE(),
	}
	readerWrapper.OutputRawData = strings.ToUpper(output.Headers.Get("x-oss-select-output-raw")) == "TRUE"
	result.Body = readerWrapper
	assert.Nil(t, err)
	err = c.unmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidRange")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "65698AC54B39ED97D780****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}
