package dom

import "encoding/xml"

// BlockInput holds a connection from one statement block to another.
// Usually rendered as:
// Usually rendered as:
// |   Input BlockInput `xml:",any,omitempty"`
// |   <next><block/></next>
type BlockInput struct {
	*Block
}

// custom serialization to toggle b/t block and shadow
func (k BlockInput) MarshalXML(enc *xml.Encoder, start xml.StartElement) (err error) {
	if n := k.Block; n != nil { // omit empty content
		if e := EncodeBlock(enc, n); e != nil {
			err = e
		}
	}
	return
}

func (k *BlockInput) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) (err error) {
	if n, e := DecodeBlock(dec, &start); e != nil {
		err = e
	} else {
		k.Block = n
	}
	return
}
