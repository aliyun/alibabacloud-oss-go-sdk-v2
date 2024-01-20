package oss

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeByExtension(t *testing.T) {
	filePath := "demo.html"
	typ := TypeByExtension(filePath)
	assert.Equal(t, "text/html", typ)

	filePath = "demo.htm"
	typ = TypeByExtension(filePath)
	assert.Equal(t, "text/html", typ)

	filePath = "demo.txt"
	typ = TypeByExtension(filePath)
	assert.Equal(t, "text/plain", typ)

	filePath = "demo"
	typ = TypeByExtension(filePath)
	assert.Equal(t, "", typ)
}
