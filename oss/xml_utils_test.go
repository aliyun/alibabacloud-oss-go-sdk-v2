package oss

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var s = `<?xml version="1.0" encoding="UTF-8"?>
<root version="1.6" writer="jack">
	<node>
		<id>12345678</id>
		<uid>1234</uid>
	</node>
	<node>
		<id>22345678</id>
		<uid>2234</uid>
	</node>
	<tag key="tag1" value="value1"/>
	<tag key="tag2" value="value2"/>
	<foo>bar</foo>
	<node/>
	<empty></empty>
</root>`

func TestDecode(t *testing.T) {
	// Decode XML document
	root := &XmlNode{}
	var err error
	var dec *XmlDecoderLite
	dec = NewXmlDecoderLite(strings.NewReader(s))
	err = dec.Decode(root)
	assert.Nil(t, err)
	value := root.GetMap()
	assert.NotNil(t, value["root"])
	items := value["root"]
	assert.Len(t, items, 6)
	nodes, ok := items.(map[string]any)
	assert.NotNil(t, nodes)
	assert.True(t, ok)

	assert.Equal(t, "1.6", nodes["+@version"].(string))
	assert.Equal(t, "jack", nodes["+@writer"].(string))

	assert.Len(t, nodes["node"], 3)
	assert.Len(t, ((nodes["node"].([]any))[0]).(map[string]any), 2)
	assert.Equal(t, "12345678", ((nodes["node"].([]any))[0]).(map[string]any)["id"].(string))
	assert.Equal(t, "1234", ((nodes["node"].([]any))[0]).(map[string]any)["uid"].(string))
	assert.Len(t, ((nodes["node"].([]any))[1]).(map[string]any), 2)
	assert.Equal(t, "22345678", ((nodes["node"].([]any))[1]).(map[string]any)["id"].(string))
	assert.Equal(t, "2234", ((nodes["node"].([]any))[1]).(map[string]any)["uid"].(string))
	assert.Nil(t, ((nodes["node"].([]any))[2]))

	assert.Len(t, nodes["tag"], 2)
	assert.Len(t, ((nodes["tag"].([]any))[0]).(map[string]any), 2)
	assert.Equal(t, "tag1", ((nodes["tag"].([]any))[0]).(map[string]any)["+@key"].(string))
	assert.Equal(t, "value1", ((nodes["tag"].([]any))[0]).(map[string]any)["+@value"].(string))
	assert.Len(t, ((nodes["tag"].([]any))[1]).(map[string]any), 2)
	assert.Equal(t, "tag2", ((nodes["tag"].([]any))[1]).(map[string]any)["+@key"].(string))
	assert.Equal(t, "value2", ((nodes["tag"].([]any))[1]).(map[string]any)["+@value"].(string))

	assert.Equal(t, "bar", nodes["foo"].(string))
	assert.Nil(t, nodes["empty"])
}

func TestTrimNonGraphic(t *testing.T) {
	table := []struct {
		in       string
		expected string
	}{
		{in: "foo", expected: "foo"},
		{in: " foo", expected: "foo"},
		{in: "foo ", expected: "foo"},
		{in: " foo ", expected: "foo"},
		{in: "   foo   ", expected: "foo"},
		{in: "foo bar", expected: "foo bar"},
		{in: "\n\tfoo\n\t", expected: "foo"},
		{in: "\n\tfoo\n\tbar\n\t", expected: "foo\n\tbar"},
		{in: "", expected: ""},
		{in: "\n", expected: ""},
		{in: "\n\v", expected: ""},
		{in: "ending with ä", expected: "ending with ä"},
		{in: "ä and ä", expected: "ä and ä"},
	}

	for _, scenario := range table {
		got := trimNonGraphic(scenario.in)
		assert.Equal(t, scenario.expected, got)
	}
}
