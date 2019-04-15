package dom

import (
	"encoding/xml"

	"github.com/ionous/errutil"
)

type Block struct {
	Id       string         `xml:"id,attr,omitempty"`
	Type     string         `xml:"type,attr,omitempty"`
	Shadow   bool           `xml:"-"`
	Mutation *BlockMutation `xml:"mutation,omitempty"`
	Next     BlockLink      `xml:"next,omitempty"` // values or fields
	Items    ItemList       `xml:",any,omitempty"` // values or fields
}

type BlockMutation struct {
	Inputs Mutations `xml:"pin"`
}
type Mutations []*Mutation

func (ms *BlockMutation) Append(in *Mutation) {
	ms.Inputs = append(ms.Inputs, in)
}

func (ms *BlockMutation) MarshalMutation() (ret string, err error) {
	out := struct {
		XMLName xml.Name  `xml:"mutation"`
		Inputs  Mutations `xml:"pin"`
	}{Inputs: ms.Inputs}

	if bytes, e := xml.Marshal(out); e != nil {
		err = e
	} else {
		ret = string(bytes)
	}
	return
}

func UnmarshalMutation(str string) (ret BlockMutation, err error) {
	var ms BlockMutation
	if e := xml.Unmarshal([]byte(str), &ms); e != nil {
		err = e
	} else {
		ret = ms
	}
	return
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
