package dom

import "encoding/xml"

type BlockList struct {
	Blocks []*Block // block or shadow
}

func (l *BlockList) Append(it *Block) {
	l.Blocks = append(l.Blocks, it)
}

func (l BlockList) MarshalXML(enc *xml.Encoder, _ xml.StartElement) (err error) {
	for _, n := range l.Blocks {
		if e := EncodeBlock(enc, n); e != nil {
			err = e
			break
		}
	}
	return
}

// called multiple times for each tag matched by a BlockList field
func (l *BlockList) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) (err error) {
	if shape, e := DecodeBlock(dec, &start); e != nil {
		err = e
	} else {
		l.Blocks = append(l.Blocks, shape)
	}
	return
}
