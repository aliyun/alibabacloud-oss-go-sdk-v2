package crypto

import (
	"io"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func randStr(n int) string {
	b := make([]rune, n)
	randMarker := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[randMarker.Intn(len(letters))]
	}
	return string(b)
}

func randLowStr(n int) string {
	return strings.ToLower(randStr(n))
}

func TestAesCtr(t *testing.T) {
	var cipherData CipherData
	cipherData.RandomKeyIv(32, 16)
	cipher, _ := newAesCtr(cipherData)

	byteReader := strings.NewReader(randLowStr(100))
	enReader := cipher.Encrypt(byteReader)
	encrypter := &CryptoEncrypter{Body: byteReader, Encrypter: enReader}
	encrypter.Close()
	buff := make([]byte, 10)
	n, err := encrypter.Read(buff)
	assert.Equal(t, 0, n)
	assert.Equal(t, io.EOF, err)

	deReader := cipher.Encrypt(byteReader)
	Decrypter := &CryptoDecrypter{Body: byteReader, Decrypter: deReader}
	Decrypter.Close()
	buff = make([]byte, 10)
	n, err = Decrypter.Read(buff)
	assert.Equal(t, 0, n)
	assert.Equal(t, io.EOF, err)
}
