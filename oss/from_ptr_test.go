package oss

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFromPtr(t *testing.T) {
	assert.Equal(t, false, ToBool(nil))
	assert.Equal(t, 0, ToInt(nil))
	assert.Equal(t, "", ToString(nil))
	assert.Equal(t, true, ToTime(nil).IsZero())
	assert.Equal(t, time.Duration(0), ToDuration(nil))

	assert.Equal(t, true, ToBool(Ptr(true)))
	assert.Equal(t, 10, ToInt(Ptr(10)))
	assert.Equal(t, "123", ToString(Ptr("123")))
	now := time.Now()
	assert.Equal(t, now, ToTime(Ptr(now)))
	assert.Equal(t, 15*time.Second, ToDuration(Ptr(15*time.Second)))
}
