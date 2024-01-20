package oss

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperationMetadata(t *testing.T) {
	m := OperationMetadata{}
	m.Get("key")
	assert.Nil(t, m.Get("key"))

	m.Set("key", "value")
	assert.Equal(t, "value", m.Get("key"))

	m.Add("key", "value2")
	assert.Equal(t, "value", m.Get("key"))
	vs := m.Values("key")
	assert.Len(t, vs, 2)
	assert.Equal(t, "value", vs[0])
	assert.Equal(t, "value2", vs[1])

	m.Set("key", "value3")
	assert.Equal(t, "value3", m.Get("key"))
	assert.Len(t, m.Values("key"), 1)

	assert.True(t, m.Has("key"))
	assert.False(t, m.Has("key1"))

	m.Add("key1", "value1-1")
	assert.True(t, m.Has("key"))
	assert.True(t, m.Has("key1"))

	m = OperationMetadata{}
	m.Add("key", "value")
	m.Add("key1", "value1")

	cm := m.Clone()

	assert.True(t, cm.Has("key"))
	assert.True(t, cm.Has("key1"))

	cm.Set("key2", "value2")

	assert.True(t, cm.Has("key"))
	assert.True(t, cm.Has("key1"))
	assert.True(t, cm.Has("key2"))

	assert.True(t, m.Has("key"))
	assert.True(t, m.Has("key1"))
	assert.False(t, m.Has("key2"))

	ms := m
	ms.Set("key2", "value2")
	assert.True(t, m.Has("key"))
	assert.True(t, m.Has("key1"))
	assert.True(t, m.Has("key2"))
}
