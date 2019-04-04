package dom

import "encoding/xml"

type ShapeList struct {
	Shapes []Shape // block or shadow
}

func (l *ShapeList) Append(it Shape) {
	l.Shapes = append(l.Shapes, it)
}

func (l ShapeList) MarshalXML(enc *xml.Encoder, _ xml.StartElement) (err error) {
	for _, n := range l.Shapes {
		if e := encodeShape(enc, n); e != nil {
			err = e
			break
		}
	}
	return
}

// called multiple times for each tag matched by ShapeList field
func (l *ShapeList) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) (err error) {
	if shape, e := decodeShape(dec, &start); e != nil {
		err = e
	} else {
		l.Shapes = append(l.Shapes, shape)
	}
	return
}
