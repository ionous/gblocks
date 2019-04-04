package dom

import "encoding/xml"

// ShapeInput holds a connection from one statement block to another.
// Usually rendered as:
// Usually rendered as:
// |   Input ShapeInput `xml:",any,omitempty"`
// |   <next><block/></next>
type ShapeInput struct {
	Shape
}

// custom serialization to toggle b/t block and shadow
func (k ShapeInput) MarshalXML(enc *xml.Encoder, start xml.StartElement) (err error) {
	if n := k.Shape; n != nil { // omit empty content
		if e := encodeShape(enc, n); e != nil {
			err = e
		}
	}
	return
}

func (k *ShapeInput) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) (err error) {
	if n, e := decodeShape(dec, &start); e != nil {
		err = e
	} else {
		k.Shape = n
	}
	return
}
