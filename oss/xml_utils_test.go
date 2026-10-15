package oss

import (
	"encoding/xml"
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

func TestXMLRoundTripDiff(t *testing.T) {
	t.Run("full match", func(t *testing.T) {
		raw := []byte(`<LifecycleConfiguration><Rule><ID>a</ID></Rule></LifecycleConfiguration>`)
		var cfg LifecycleConfiguration
		assert.NoError(t, xml.Unmarshal(raw, &cfg))
		missing, err := XMLRoundTripDiff(raw, &cfg)
		assert.NoError(t, err)
		assert.Nil(t, missing)
	})

	t.Run("unknown child element in Rule", func(t *testing.T) {
		raw := []byte(`<LifecycleConfiguration><Rule><ID>a</ID><NewExtension>foo</NewExtension></Rule></LifecycleConfiguration>`)
		var cfg LifecycleConfiguration
		assert.NoError(t, xml.Unmarshal(raw, &cfg))
		missing, err := XMLRoundTripDiff(raw, &cfg)
		assert.NoError(t, err)
		assert.Equal(t, []string{"LifecycleConfiguration/Rule/NewExtension"}, missing)
	})

	t.Run("attributes are ignored", func(t *testing.T) {
		raw := []byte(`<LifecycleConfiguration version="2"><Rule><ID>a</ID></Rule></LifecycleConfiguration>`)
		var cfg LifecycleConfiguration
		assert.NoError(t, xml.Unmarshal(raw, &cfg))
		missing, err := XMLRoundTripDiff(raw, &cfg)
		assert.NoError(t, err)
		assert.Nil(t, missing)
	})

	t.Run("deep nesting unknown field", func(t *testing.T) {
		raw := []byte(`<LifecycleConfiguration><Rule><ID>a</ID><Filter><Not><Prefix>pre</Prefix><NewTagExt>v</NewTagExt></Not></Filter></Rule></LifecycleConfiguration>`)
		var cfg LifecycleConfiguration
		assert.NoError(t, xml.Unmarshal(raw, &cfg))
		missing, err := XMLRoundTripDiff(raw, &cfg)
		assert.NoError(t, err)
		assert.Equal(t, []string{"LifecycleConfiguration/Rule/Filter/Not/NewTagExt"}, missing)
	})

	t.Run("multi sibling only second has unknown field", func(t *testing.T) {
		raw := []byte(`<LifecycleConfiguration><Rule><ID>a</ID></Rule><Rule><ID>b</ID><NewExtension>x</NewExtension></Rule></LifecycleConfiguration>`)
		var cfg LifecycleConfiguration
		assert.NoError(t, xml.Unmarshal(raw, &cfg))
		missing, err := XMLRoundTripDiff(raw, &cfg)
		assert.NoError(t, err)
		assert.Equal(t, []string{"LifecycleConfiguration/Rule/NewExtension"}, missing)
	})

	t.Run("whitespace and text nodes ignored", func(t *testing.T) {
		raw := []byte("<LifecycleConfiguration>\n  <Rule>\n    <ID>a</ID>\n  </Rule>\n</LifecycleConfiguration>")
		var cfg LifecycleConfiguration
		assert.NoError(t, xml.Unmarshal(raw, &cfg))
		missing, err := XMLRoundTripDiff(raw, &cfg)
		assert.NoError(t, err)
		assert.Nil(t, missing)
	})

	t.Run("empty input nil", func(t *testing.T) {
		missing, err := XMLRoundTripDiff(nil, &LifecycleConfiguration{})
		assert.NoError(t, err)
		assert.Nil(t, missing)
	})

	t.Run("empty input bytes", func(t *testing.T) {
		missing, err := XMLRoundTripDiff([]byte("   \n\t"), &LifecycleConfiguration{})
		assert.NoError(t, err)
		assert.Nil(t, missing)
	})

	t.Run("malformed XML returns error", func(t *testing.T) {
		raw := []byte(`<a><b></a>`)
		missing, err := XMLRoundTripDiff(raw, &LifecycleConfiguration{})
		assert.Error(t, err)
		assert.Nil(t, missing)
	})

	t.Run("round-trip element order differs", func(t *testing.T) {
		// struct order emits ID before Prefix; raw has them reversed.
		raw := []byte(`<LifecycleConfiguration><Rule><Prefix>p</Prefix><ID>a</ID></Rule></LifecycleConfiguration>`)
		var cfg LifecycleConfiguration
		assert.NoError(t, xml.Unmarshal(raw, &cfg))
		missing, err := XMLRoundTripDiff(raw, &cfg)
		assert.NoError(t, err)
		assert.Nil(t, missing)
	})

	t.Run("unknown field under repeated sibling does not falsely drag Tag", func(t *testing.T) {
		// Two <Not> siblings; the second one carries an unknown <NewExt>.
		// Regression guard: an earlier sort-based pairing reported Tag as
		// missing because the unknown field skewed sibling sort order.
		raw := []byte(`<LifecycleConfiguration><Rule><ID>r</ID><Filter>` +
			`<Not><Prefix>a</Prefix><Tag><Key>k</Key></Tag></Not>` +
			`<Not><Prefix>b</Prefix><NewExt>v</NewExt></Not>` +
			`</Filter></Rule></LifecycleConfiguration>`)
		var cfg LifecycleConfiguration
		assert.NoError(t, xml.Unmarshal(raw, &cfg))
		missing, err := XMLRoundTripDiff(raw, &cfg)
		assert.NoError(t, err)
		assert.Equal(t, []string{"LifecycleConfiguration/Rule/Filter/Not/NewExt"}, missing)
	})

	t.Run("multiple siblings mapped to single-pointer field", func(t *testing.T) {
		// LifecycleRuleNot.Tag is *Tag (single pointer); two <Tag> siblings
		// collapse into one on round-trip, so the diff must report Tag as
		// losing one instance.
		raw := []byte(`<LifecycleConfiguration><Rule><ID>r</ID><Filter>` +
			`<Not><Prefix>p</Prefix><Tag><Key>k1</Key></Tag><Tag><Key>k2</Key></Tag></Not>` +
			`</Filter></Rule></LifecycleConfiguration>`)
		var cfg LifecycleConfiguration
		assert.NoError(t, xml.Unmarshal(raw, &cfg))
		missing, err := XMLRoundTripDiff(raw, &cfg)
		assert.NoError(t, err)
		assert.Contains(t, missing, "LifecycleConfiguration/Rule/Filter/Not/Tag")
	})

	t.Run("rt emits non-omitempty element raw lacks is not reported", func(t *testing.T) {
		// Non-pointer field without omitempty: rt may emit elements raw lacks;
		// the one-directional diff must not flag these.
		raw := []byte(`<LifecycleConfiguration><Rule><ID>r</ID>` +
			`<Transition><Days>30</Days></Transition></Rule></LifecycleConfiguration>`)
		var cfg LifecycleConfiguration
		assert.NoError(t, xml.Unmarshal(raw, &cfg))
		missing, err := XMLRoundTripDiff(raw, &cfg)
		assert.NoError(t, err)
		assert.Nil(t, missing)
	})
}

func TestUnknownFieldsError_Error(t *testing.T) {
	ufe := &UnknownFieldsError{MissingPaths: []string{
		"LifecycleConfiguration/Rule/NewExtension",
		"LifecycleConfiguration/Rule/Filter/Not/NewFilterExt",
	}}
	msg := ufe.Error()
	assert.Contains(t, msg, "2")
	assert.Contains(t, msg, "LifecycleConfiguration/Rule/NewExtension")
	assert.Contains(t, msg, "LifecycleConfiguration/Rule/Filter/Not/NewFilterExt")
}
