package vectors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
)

type VectorsClient struct {
	clientImpl *oss.Client
}

func NewVectorsClient(cfg *oss.Config, optFns ...func(*oss.Options)) *VectorsClient {
	newCfg := cfg.Copy()
	updateEndpoint(&newCfg)
	updateUserAgent(&newCfg)
	var signer = resolveSigner(&newCfg)
	vectorsOptFn := func(options *oss.Options) {
		options.Signer = signer
		options.EndpointProvider = &endpointProvider{
			accountId:    oss.ToString(newCfg.AccountId),
			endpoint:     options.Endpoint,
			endpointType: options.UrlStyle,
		}
	}
	allOptFns := append(optFns, vectorsOptFn)
	return &VectorsClient{
		clientImpl: oss.NewClient(&newCfg, allOptFns...),
	}
}

func updateEndpoint(cfg *oss.Config) {
	if len(oss.ToString(cfg.Endpoint)) > 0 {
		return
	}

	region := oss.ToString(cfg.Region)

	if !oss.IsValidRegion(region) {
		return
	}

	if oss.ToBool(cfg.UseInternalEndpoint) {
		cfg.Endpoint = oss.Ptr(fmt.Sprintf("oss-%s-internal.oss-vectors.aliyuncs.com", region))

	} else {
		cfg.Endpoint = oss.Ptr(fmt.Sprintf("oss-%s.oss-vectors.aliyuncs.com", region))
	}
}

func updateUserAgent(cfg *oss.Config) {
	userAgent := "vectors-client"

	if cfg.UserAgent != nil {
		userAgent = fmt.Sprintf("%s/%s", userAgent, oss.ToString(cfg.UserAgent))
	}
	cfg.UserAgent = oss.Ptr(userAgent)
}

func resolveSigner(cfg *oss.Config) signer.Signer {
	ver := oss.DefaultSignatureVersion
	if cfg.SignatureVersion != nil {
		ver = *cfg.SignatureVersion
	}

	switch ver {
	case oss.SignatureVersionV1:
		return &signer.SignerVectorsV1{
			AccountId: cfg.AccountId,
		}
	default:
		return &signer.SignerVectorsV4{
			AccountId: cfg.AccountId,
		}
	}
}

// fieldInfo holds details for the input/output of a single field.
type fieldInfo struct {
	idx   int
	flags int
}

const (
	fRequire int = 1 << iota

	fTypeUsermeta
	fTypeTime
	fTypeJson
)

func parseFiledFlags(tokens []string) int {
	var flags int = 0
	for _, token := range tokens {
		switch token {
		case "required":
			flags |= fRequire
		case "time":
			flags |= fTypeTime
		case "xml", "json":
			flags |= fTypeJson
		case "usermeta":
			flags |= fTypeUsermeta
		}
	}
	return flags
}

func validateInput(input *oss.OperationInput) error {
	if input == nil {
		return oss.NewErrParamNull("OperationInput")
	}

	if input.Bucket != nil && !oss.IsValidBucketName(input.Bucket) {
		return oss.NewErrParamInvalid("OperationInput.Bucket")
	}

	if !oss.IsValidMethod(input.Method) {
		return oss.NewErrParamInvalid("OperationInput.Method")
	}

	return nil
}

// oss api's style, but the body type of request is json.
func (c *VectorsClient) marshalInput(request any, input *oss.OperationInput, handlers ...func(any, *oss.OperationInput) error) error {
	// merge common fields
	if cm, ok := request.(oss.RequestCommonInterface); ok {
		h, p, b := cm.GetCommonFileds()
		// headers
		if len(h) > 0 {
			if input.Headers == nil {
				input.Headers = map[string]string{}
			}
			for k, v := range h {
				input.Headers[k] = v
			}
		}

		// parameters
		if len(p) > 0 {
			if input.Parameters == nil {
				input.Parameters = map[string]string{}
			}
			for k, v := range p {
				input.Parameters[k] = v
			}
		}

		// body
		input.Body = b
	}

	val := reflect.ValueOf(request)
	switch val.Kind() {
	case reflect.Pointer, reflect.Interface:
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct || input == nil {
		return nil
	}

	t := val.Type()
	for k := 0; k < t.NumField(); k++ {
		if tag, ok := t.Field(k).Tag.Lookup("input"); ok {
			// header|query|body,filed_name,[required,time,usermeta...]
			v := val.Field(k)
			var flags int = 0
			tokens := strings.Split(tag, ",")
			if len(tokens) < 2 {
				continue
			}

			// parse field flags
			if len(tokens) > 2 {
				flags = parseFiledFlags(tokens[2:])
			}
			// check required flag
			if oss.IsEmptyValue(v) {
				if flags&fRequire != 0 {
					return oss.NewErrParamRequired(t.Field(k).Name)
				}
				continue
			}

			switch tokens[0] {
			case "query":
				if input.Parameters == nil {
					input.Parameters = map[string]string{}
				}
				if v.Kind() == reflect.Pointer {
					v = v.Elem()
				}
				input.Parameters[tokens[1]] = fmt.Sprintf("%v", v.Interface())
			case "header":
				if input.Headers == nil {
					input.Headers = map[string]string{}
				}
				if v.Kind() == reflect.Pointer {
					v = v.Elem()
				}
				if flags&fTypeUsermeta != 0 {
					if m, ok := v.Interface().(map[string]string); ok {
						for k, v := range m {
							input.Headers[tokens[1]+k] = v
						}
					}
				} else {
					input.Headers[tokens[1]] = fmt.Sprintf("%v", v.Interface())
				}
			case "body":
				switch {
				case flags&fTypeJson != 0:
					var b bytes.Buffer
					var err error
					wrapper := map[string]interface{}{
						tokens[1]: v.Interface(),
					}
					encoder := json.NewEncoder(&b)
					encoder.SetEscapeHTML(false)
					err = encoder.Encode(wrapper)
					input.Body = bytes.NewReader(bytes.TrimRight(b.Bytes(), "\n"))
					if err != nil {
						return &oss.SerializationError{Err: err}
					}
				default:
					if r, ok := v.Interface().(io.Reader); ok {
						input.Body = r
					} else {
						return oss.NewErrParamTypeNotSupport(t.Field(k).Name)
					}
				}
			}
		}
	}

	if err := validateInput(input); err != nil {
		return err
	}

	for _, h := range handlers {
		if err := h(request, input); err != nil {
			return err
		}
	}

	return nil
}

// oss new api's style
func (c *VectorsClient) marshalInputJson(request any, input *oss.OperationInput, handlers ...func(any, *oss.OperationInput) error) error {
	// merge common fields
	if cm, ok := request.(oss.RequestCommonInterface); ok {
		h, p, b := cm.GetCommonFileds()
		// headers
		if len(h) > 0 {
			if input.Headers == nil {
				input.Headers = map[string]string{}
			}
			for k, v := range h {
				input.Headers[k] = v
			}
		}

		// parameters
		if len(p) > 0 {
			if input.Parameters == nil {
				input.Parameters = map[string]string{}
			}
			for k, v := range p {
				input.Parameters[k] = v
			}
		}

		// body
		input.Body = b
	}

	val := reflect.ValueOf(request)
	switch val.Kind() {
	case reflect.Pointer, reflect.Interface:
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct || input == nil {
		return nil
	}

	t := val.Type()
	bodyMap := map[string]any{}
	for k := 0; k < t.NumField(); k++ {
		if tag, ok := t.Field(k).Tag.Lookup("input"); ok {
			// header|query|body,filed_name,[required,time,usermeta...]
			v := val.Field(k)
			var flags int = 0
			tokens := strings.Split(tag, ",")
			if len(tokens) < 2 {
				continue
			}

			// parse field flags
			if len(tokens) > 2 {
				flags = parseFiledFlags(tokens[2:])
			}
			// check required flag
			if oss.IsEmptyValue(v) {
				if flags&fRequire != 0 {
					return oss.NewErrParamRequired(t.Field(k).Name)
				}
				continue
			}

			switch tokens[0] {
			case "query":
				if input.Parameters == nil {
					input.Parameters = map[string]string{}
				}
				if v.Kind() == reflect.Pointer {
					v = v.Elem()
				}
				input.Parameters[tokens[1]] = fmt.Sprintf("%v", v.Interface())
			case "header":
				if input.Headers == nil {
					input.Headers = map[string]string{}
				}
				if v.Kind() == reflect.Pointer {
					v = v.Elem()
				}
				if flags&fTypeUsermeta != 0 {
					if m, ok := v.Interface().(map[string]string); ok {
						for k, v := range m {
							input.Headers[tokens[1]+k] = v
						}
					}
				} else {
					input.Headers[tokens[1]] = fmt.Sprintf("%v", v.Interface())
				}
			case "body":
				switch {
				case flags&fTypeJson != 0:
					bodyMap[tokens[1]] = v.Interface()
				default:
					if r, ok := v.Interface().(io.Reader); ok {
						input.Body = r
					} else {
						return oss.NewErrParamTypeNotSupport(t.Field(k).Name)
					}
				}
			}
		}
	}

	if len(bodyMap) != 0 {
		b, err := json.Marshal(bodyMap)
		if err != nil {
			return &oss.SerializationError{Err: err}
		}
		input.Body = bytes.NewReader(b)
	}

	if err := validateInput(input); err != nil {
		return err
	}

	for _, h := range handlers {
		if err := h(request, input); err != nil {
			return err
		}
	}

	return nil
}

func (c *VectorsClient) unmarshalOutput(result any, output *oss.OperationOutput, handlers ...func(any, *oss.OperationOutput) error) error {
	// Common
	if cm, ok := result.(oss.ResultCommonInterface); ok {
		cm.CopyIn(output.Status, output.StatusCode, output.Headers, output.OpMetadata)
	}

	var err error
	for _, h := range handlers {
		if err = h(result, output); err != nil {
			break
		}
	}
	return err
}

func unmarshalBodyLikeXmlJson2(result any, output *oss.OperationOutput) error {
	var err error
	var body []byte
	if output.Body != nil {
		defer output.Body.Close()
		if body, err = io.ReadAll(output.Body); err != nil {
			return err
		}
	}

	var root map[string]json.RawMessage
	// extract body
	if len(body) > 0 {
		if err = json.Unmarshal(body, &root); err == nil {
			if len(root) == 1 {
				for _, v := range root {
					err = json.Unmarshal(v, result)
				}
			} else {
				// invalid json format
				err = fmt.Errorf("Not a valid json format")
			}
		}
	}

	if err != nil {
		err = &oss.DeserializationError{
			Err:      err,
			Snapshot: body,
		}
	}

	return err
}

func unmarshalBodyLikeXmlJson(result any, output *oss.OperationOutput) error {
	var err error
	var body []byte
	if output.Body != nil {
		defer output.Body.Close()
		if body, err = io.ReadAll(output.Body); err != nil {
			return err
		}
	}

	if len(body) == 0 {
		return nil
	}

	val := reflect.ValueOf(result)
	switch val.Kind() {
	case reflect.Pointer, reflect.Interface:
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct || output == nil {
		return nil
	}

	t := val.Type()
	idx := -1
	var filedName string
	for k := 0; k < t.NumField(); k++ {
		if tag, ok := t.Field(k).Tag.Lookup("output"); ok {
			tokens := strings.Split(tag, ",")
			if len(tokens) < 2 {
				continue
			}
			// header|query|body,filed_name,[required,time,usermeta...]
			switch tokens[0] {
			case "body":
				idx = k
				filedName = tokens[1]
				break
			}
		}
	}

	if idx >= 0 {
		var rawMessage []byte
		filedNames := strings.Split(filedName, ">")
		rawMessage = body
		for _, fname := range filedNames {
			var root map[string]json.RawMessage
			err = json.Unmarshal(rawMessage, &root)
			if err != nil {
				break
			}

			if len(root) != 1 {
				err = fmt.Errorf("Not a expected json format")
			}

			for k, v := range root {
				if k != fname {
					err = fmt.Errorf("Found key %s, but expect %s", k, fname)
					break
				}
				rawMessage = v
			}
		}

		if err == nil {
			dst := val.Field(idx)
			if dst.IsNil() {
				dst.Set(reflect.New(dst.Type().Elem()))
			}
			if len(rawMessage) != 0 {
				err = json.Unmarshal(rawMessage, dst.Interface())
			} else {
				// not found
				len := len(body)
				if len > 256 {
					len = 256
				}
				err = fmt.Errorf("Not found Root %s in response body. With part response body %s.", filedName, string(body[:len]))
			}
		}
	}

	if err != nil {
		err = &oss.DeserializationError{
			Err:      err,
			Snapshot: body,
		}
	}

	return err
}

func unmarshalBodyJsonStyle(result any, output *oss.OperationOutput) error {
	var err error
	var body []byte
	if output.Body != nil {
		defer output.Body.Close()
		if body, err = io.ReadAll(output.Body); err != nil {
			return err
		}
	}

	// extract body
	if len(body) > 0 {
		contentType := output.Headers.Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			err = json.Unmarshal(body, result)
		} else {
			err = fmt.Errorf("unsupport contentType:%s", contentType)
		}
		if err != nil {
			err = &oss.DeserializationError{
				Err:      err,
				Snapshot: body,
			}
		}
	}
	return err
}

func (c *VectorsClient) toClientError(err error, code string, output *oss.OperationOutput) error {
	if err == nil {
		return nil
	}

	return &oss.ClientError{
		Code: code,
		Message: fmt.Sprintf("execute %s fail, error code is %s, request id:%s",
			output.Input.OpName,
			code,
			output.Headers.Get(oss.HeaderOssRequestID),
		),
		Err: err}
}

const (
	contentTypeDefault = "application/octet-stream"
	contentTypeXML     = "application/xml"
	contentTypeJSON    = "application/json"
)
