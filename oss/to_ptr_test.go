package oss

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPtr(t *testing.T) {
	b := true
	pb := Ptr(b)
	assert.NotNil(t, pb)
	assert.Equal(t, b, *pb)
}

func TestSliceOfPtrs(t *testing.T) {
	arr := SliceOfPtrs[int]()
	assert.Equal(t, len(arr), 0)
	arr = SliceOfPtrs(1, 2, 3, 4, 5)
	for i, v := range arr {
		assert.Equal(t, i+1, *v)
	}
}
