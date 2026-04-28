package dataprocess

import (
	"encoding/xml"
	"io"
)

type MapEntry map[string]string

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func (m *MapEntry) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = MapEntry{}
	for {
		var e xmlMapEntry
		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}
