package oss

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCRC64(t *testing.T) {
	// hash first part
	d := "123456789"
	c := NewCRC64(0)
	c.Write([]byte(d))
	sum := c.Sum64()
	assert.Equal(t, uint64(0x995dc9bbdf1939fa), sum)

	// hash second part
	d1 := "This is a test of the emergency broadcast system."
	c = NewCRC64(0)
	c.Write([]byte(d1))
	sum1 := c.Sum64()
	assert.Equal(t, uint64(0x27db187fc15bbc72), sum1)

	//combine
	d2 := d + d1
	c = NewCRC64(0)
	c.Write([]byte(d2))
	sum2 := c.Sum64()

	csum := CRC64Combine(sum, sum1, uint64(len(d1)))
	assert.Equal(t, sum2, csum)

	// init from first part
	c = NewCRC64(sum)
	c.Write([]byte(d1))
	assert.Equal(t, sum2, c.Sum64())

	//reset
	c = NewCRC64(sum)
	c.Write([]byte(d))
	assert.NotEqual(t, sum2, c.Sum64())
	c.Reset()
	c.Write([]byte(d1))
	assert.Equal(t, sum2, c.Sum64())
}
