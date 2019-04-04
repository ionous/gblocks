package dom

import "encoding/xml"

type Mutation struct {
	Input string `xml:"name,attr"`
	Atoms Atoms  `xml:"atom,omitempty"`
}

// employs a custom marshaller to write a flat array of strings to dom elements with attributes
// <atom type="name">
type Atoms struct {
	Types []string
}

type atomAttr struct {
	Type string `xml:"type,attr"`
}

func (a Atoms) MarshalXML(enc *xml.Encoder, start xml.StartElement) (err error) {
	for _, t := range a.Types {
		if e := enc.EncodeElement(atomAttr{t}, start); e != nil {
			err = e
			break
		}
	}
	return
}

func (a *Atoms) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	// FIX: speed up by decoding directly?
	// d.Token() Token holds either: StartElement, EndElement, CharData, Comment, ProcInst, or Directive.
	var attrs []atomAttr
	if e := d.DecodeElement(&attrs, &start); err != nil {
		err = e
	} else {
		for _, attr := range attrs {
			a.Types = append(a.Types, attr.Type)
		}
	}
	return
}
