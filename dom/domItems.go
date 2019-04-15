package dom

import (
	"encoding/xml"

	"github.com/ionous/errutil"
)

type Item interface{ Item() Item }

// ex: <value name=""><block type=""></block></value>
type Value struct {
	Name  string     `xml:"name,attr"`
	Input BlockInput `xml:",any,omitempty"`
}

func (it *Value) Item() Item { return it }

// ex: <statement name=""><block type=""></block></statement>
type Statement struct {
	Name  string     `xml:"name,attr"`
	Input BlockInput `xml:",any,omitempty"`
}

func (it *Statement) Item() Item { return it }

// ex. <field name="NUMBER">10</field>
type Field struct {
	Name    string `xml:"name,attr"`
	Content string `xml:",innerxml"`
}

func (it *Field) Item() Item { return it }

type ItemList struct {
	Items Items
}
type Items []Item

func (l *ItemList) Append(it Item) {
	l.Items = append(l.Items, it)
}

func (l ItemList) MarshalXML(enc *xml.Encoder, _ xml.StartElement) (err error) {
	for _, item := range l.Items {
		switch item := item.(type) {
		case *Value:
			err = enc.EncodeElement(item, xml.StartElement{Name: names.value})
		case *Statement:
			err = enc.EncodeElement(item, xml.StartElement{Name: names.statement})
		case *Field:
			err = enc.EncodeElement(item, xml.StartElement{Name: names.field})
		default:
			err = errutil.New("unknown type", item)
		}
	}
	return
}

// called multiple times for each tag matched by BlockList field
func (l *ItemList) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) (err error) {
	switch start.Name {
	case names.value:
		out := new(Value)
		if e := dec.DecodeElement(out, &start); e != nil {
			err = e
		} else {
			l.Items = append(l.Items, out)
		}
	case names.statement:
		out := new(Statement)
		if e := dec.DecodeElement(out, &start); e != nil {
			err = e
		} else {
			l.Items = append(l.Items, out)
		}
	case names.field:
		out := new(Field)
		if e := dec.DecodeElement(out, &start); e != nil {
			err = e
		} else {
			l.Items = append(l.Items, out)
		}
	default:
		err = errutil.New("unknown element", start.Name.Local)
	}
	return
}
