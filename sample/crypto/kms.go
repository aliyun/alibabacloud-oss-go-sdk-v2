package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"io/ioutil"
	"log"
	"strings"

	kms "github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	kmssdk "github.com/aliyun/alibabacloud-dkms-transfer-go-sdk/sdk"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	osscrypto "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/crypto"
)

// CreateMasterAliKms3 Create master key interface implemented by ali kms 3.0
// matDesc will be converted to json string
func CreateMasterAliKms3(matDesc map[string]string, kmsID string, kmsClient *kmssdk.KmsTransferClient) (osscrypto.MasterCipher, error) {
	var masterCipher MasterAliKms3Cipher
	if kmsID == "" || kmsClient == nil {
		return masterCipher, fmt.Errorf("kmsID is empty or kmsClient is nil")
	}

	var jsonDesc string
	if len(matDesc) > 0 {
		b, err := json.Marshal(matDesc)
		if err != nil {
			return masterCipher, err
		}
		jsonDesc = string(b)
	}

	masterCipher.MatDesc = jsonDesc
	masterCipher.KmsID = kmsID
	masterCipher.KmsClient = kmsClient
	return masterCipher, nil
}

// MasterAliKms3Cipher ali kms master key interface
type MasterAliKms3Cipher struct {
	MatDesc   string
	KmsID     string
	KmsClient *kmssdk.KmsTransferClient
}

// GetWrapAlgorithm get master key wrap algorithm
func (mrc MasterAliKms3Cipher) GetWrapAlgorithm() string {
	return osscrypto.KmsAliCryptoWrap
}

// GetMatDesc get master key describe
func (mkms MasterAliKms3Cipher) GetMatDesc() string {
	return mkms.MatDesc
}

// Encrypt  encrypt data by ali kms
// Mainly used to encrypt object's symmetric secret key and iv
func (mkms MasterAliKms3Cipher) Encrypt(plainData []byte) ([]byte, error) {
	// kms Plaintext must be base64 encoded
	base64Plain := base64.StdEncoding.EncodeToString(plainData)
	request := kms.CreateEncryptRequest()
	request.RpcRequest.Scheme = "https"
	request.RpcRequest.Method = "POST"
	request.RpcRequest.AcceptFormat = "json"

	request.KeyId = mkms.KmsID
	request.Plaintext = base64Plain

	response, err := mkms.KmsClient.Encrypt(request)
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(response.CiphertextBlob)
}

// Decrypt decrypt data by ali kms
// Mainly used to decrypt object's symmetric secret key and iv
func (mkms MasterAliKms3Cipher) Decrypt(cryptoData []byte) ([]byte, error) {
	base64Crypto := base64.StdEncoding.EncodeToString(cryptoData)
	request := kms.CreateDecryptRequest()
	request.RpcRequest.Scheme = "https"
	request.RpcRequest.Method = "POST"
	request.RpcRequest.AcceptFormat = "json"
	request.CiphertextBlob = string(base64Crypto)
	response, err := mkms.KmsClient.Decrypt(request)
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(response.Plaintext)
}

var (
	region     string
	bucketName string
	objectName string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
	flag.StringVar(&objectName, "object", "", "The name of the object.")
}

func main() {
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := oss.NewClient(cfg)

	kmsRegion := "cn-hangzhou"
	kmsAccessKeyId := "access key id"
	kmsAccessKeySecret := "access key secret"
	kmsKeyId := "kms id"

	kmsClient, err := kmssdk.NewClientWithAccessKey(kmsRegion, kmsAccessKeyId, kmsAccessKeySecret, nil)
	if err != nil {
		log.Fatalf("failed to kms sdk client%v", err)
	}
	materialDesc := make(map[string]string)
	materialDesc["desc"] = "your kms encrypt key material describe information"

	masterKmsCipher, err := CreateMasterAliKms3(materialDesc, kmsKeyId, kmsClient)
	if err != nil {
		log.Fatalf("failed to create master AliKms3 %v", err)
	}

	eclient, err := oss.NewEncryptionClient(client, masterKmsCipher)
	request := &oss.PutObjectRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
		Body:   strings.NewReader("hi kms"),
	}
	result, err := eclient.PutObject(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put object with encryption client %v", err)
	}
	log.Printf("put object with encryption client result:%#v\n", result)

	getRequest := &oss.GetObjectRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
	}
	getResult, err := eclient.GetObject(context.TODO(), getRequest)
	if err != nil {
		log.Fatalf("failed to get object with encryption client %v", err)
	}
	defer getResult.Body.Close()

	data, err := ioutil.ReadAll(getResult.Body)
	if err != nil {
		log.Fatalf("failed to read all %v", err)
	}
	log.Printf("get object data:%s\n", data)
}
