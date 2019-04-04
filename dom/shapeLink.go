package dom

import "encoding/xml"

// ShapeLink holds a connection from one statement block to another.
// Usually rendered as:
// |   Next ShapeLink `xml:"next,omitempty"`
// |   <next><block/></next>
type ShapeLink struct {
	Shape
}

// custom serialization to toggle b/t block and shadow
func (k ShapeLink) MarshalXML(enc *xml.Encoder, start xml.StartElement) (err error) {
	// omit empty content
	if n := k.Shape; n != nil {
		// start contains the <next> opening tag
		if e := enc.EncodeToken(start); e != nil {
			err = e
		} else if e := encodeShape(enc, n); e != nil {
			err = e
		} else if e := enc.EncodeToken(xml.EndElement{start.Name}); e != nil {
			err = e
		}
	}
	return
}

func (k *ShapeLink) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) (err error) {
	// note: start is usually <next>
	return parse(dec, &parser{
		EndElement: func(curr *xml.EndElement) (okay bool, err error) {
			return // all done.
		},
		// <block> or <shadow>
		StartElement: func(curr *xml.StartElement) (okay bool, err error) {
			if n, e := decodeShape(dec, curr); e != nil {
				err = e
			} else {
				k.Shape, okay = n, true // keep parsing till the end element.
			}
			return
		},
	})
}
