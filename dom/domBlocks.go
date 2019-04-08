package dom

import (
	"encoding/xml"

	"github.com/ionous/errutil"
)

type Block struct {
	Id        string     `xml:"id,attr,omitempty"`
	Type      string     `xml:"type,attr"`
	Shadow    bool       `xml:"-"`
	Mutations *Mutations `xml:"mutation>input,omitempty"`
	Next      BlockLink  `xml:"next,omitempty"` // values or fields
	Items     ItemList   `xml:",any,omitempty"` // values or fields
}

func (b *Block) AppendMutation(m *Mutation) {
	(*b.Mutations) = append((*b.Mutations), m)
}

func EncodeBlock(enc *xml.Encoder, b *Block) (err error) {
	var n xml.Name
	if b.Shadow {
		n = names.shadow
	} else {
		n = names.block
	}
	return enc.EncodeElement(b, xml.StartElement{Name: n})
}

func DecodeBlock(dec *xml.Decoder, start *xml.StartElement) (ret *Block, err error) {
	var out Block
	if e := dec.DecodeElement(&out, start); e != nil {
		err = e
	} else {
		switch start.Name {
		case names.shadow:
			out.Shadow = true
			ret = &out
		case names.block:
			out.Shadow = false
			ret = &out
		default:
			err = errutil.New("unknown element", start.Name.Local)
		}
	}
	return
}
