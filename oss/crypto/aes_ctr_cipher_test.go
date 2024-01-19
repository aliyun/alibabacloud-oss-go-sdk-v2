package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContentEncryptCipherError(t *testing.T) {
	masterRsaCipher, _ := CreateMasterRsa(matDesc, rsaPublicKey, rsaPrivateKey)
	contentProvider := CreateAesCtrCipher(masterRsaCipher)
	cc, err := contentProvider.ContentCipher()
	assert.Nil(t, err)

	var cipherData CipherData
	cipherData.RandomKeyIv(31, 15)

	_, err = cc.Clone(cipherData)
	assert.NotNil(t, err)
}

func TestCreateCipherDataError(t *testing.T) {
	// crypto bucket
	masterRsaCipher, _ := CreateMasterRsa(matDesc, "", "")
	contentProvider := CreateAesCtrCipher(masterRsaCipher)

	v := contentProvider.(aesCtrCipherBuilder)
	_, err := v.createCipherData()
	assert.NotNil(t, err)
}

func TestContentCipherCDError(t *testing.T) {
	var cd CipherData

	// crypto bucket
	masterRsaCipher, _ := CreateMasterRsa(matDesc, "", "")
	contentProvider := CreateAesCtrCipher(masterRsaCipher)

	v := contentProvider.(aesCtrCipherBuilder)
	_, err := v.contentCipherCD(cd)
	assert.NotNil(t, err)

	_, err = contentProvider.ContentCipherEnv(Envelope{})
	assert.NotNil(t, err)
}
