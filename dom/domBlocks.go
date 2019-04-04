package dom

import (
	"encoding/xml"

	"github.com/ionous/errutil"
)

type Block struct {
	Type      string     `xml:"type,attr"`
	Mutations *Mutations `xml:"mutation>input,omitempty"`
	Next      ShapeLink  `xml:"next,omitempty"` // values or fields
	Items     ItemList   `xml:",any,omitempty"` // values or fields
}

type Shadow Block

type Shape interface{ Shape() Shape }

func (b *Block) Shape() Shape  { return b }
func (s *Shadow) Shape() Shape { return s }

func (b *Block) AppendMutation(m *Mutation) {
	(*b.Mutations) = append((*b.Mutations), m)
}

func encodeShape(enc *xml.Encoder, shape Shape) (err error) {
	switch n := shape.(type) {
	case *Block:
		err = enc.EncodeElement(n, xml.StartElement{Name: names.block})
	case *Shadow:
		err = enc.EncodeElement(n, xml.StartElement{Name: names.shadow})
	default:
		err = errutil.New("unknown type", n)
	}
	return
}

func decodeShape(dec *xml.Decoder, start *xml.StartElement) (ret Shape, err error) {
	switch start.Name {
	case names.shadow:
		var out Shadow
		ret, err = &out, dec.DecodeElement(&out, start)
	case names.block:
		var out Block
		ret, err = &out, dec.DecodeElement(&out, start)
	default:
		err = errutil.New("unknown element", start.Name.Local)
	}
	return
}
