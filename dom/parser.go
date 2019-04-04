package dom

import (
	"encoding/xml"

	"github.com/ionous/errutil"
)

// parser helps to simplify the reading and handling of tokens.
// there's an (optional) callback per token type.
// each callback can return false to end the parsing, or true to keep going.
// tokens w/o callbacks are eaten.
type parser struct {
	// <tag>
	StartElement func(*xml.StartElement) (keepGoing bool, err error)
	// </tag>
	EndElement func(*xml.EndElement) (keepGoing bool, err error)
	// <![CDATA[ ... ]]>
	CharData func(*xml.CharData) (keepGoing bool, err error)
	// <!-- -->
	Comment func(*xml.Comment) (keepGoing bool, err error)
	// <?target inst?>
	ProcInst func(*xml.ProcInst) (keepGoing bool, err error)
	// ex. <!DOCTYPE ...>
	Directive func(*xml.Directive) (keepGoing bool, err error)
}

// parsing note:
// for the type Toolbox { Shapes ShapeList `xml:"block,omitempty"` }
// if ShapeList implements UnmarshalXml() the start.Name.Local == block.
// the first Decoder.Token() will then be whatever follows that starting block;
// an immediate EndElement, for example, would indicate an empty element </>.
// a series of elements (</>,</>,</>) means UnmarshalXML() gets called multiple times.
func parse(dec *xml.Decoder, p *parser) (err error) {
	for keepGoing := true; err == nil && keepGoing; {
		if token, e := dec.Token(); e != nil {
			err = e
		} else {
			switch curr := token.(type) {
			case xml.StartElement:
				if p.StartElement != nil {
					keepGoing, err = p.StartElement(&curr)
				}
			case xml.EndElement:
				if p.EndElement != nil {
					keepGoing, err = p.EndElement(&curr)
				}
			case xml.CharData:
				if p.CharData != nil {
					keepGoing, err = p.CharData(&curr)
				}
			case xml.Comment:
				if p.Comment != nil {
					keepGoing, err = p.Comment(&curr)
				}
			case xml.ProcInst:
				if p.ProcInst != nil {
					keepGoing, err = p.ProcInst(&curr)
				}
			case xml.Directive:
				if p.Directive != nil {
					keepGoing, err = p.Directive(&curr)
				}
			default:
				err = errutil.New("unknown token", curr)
			}
		}
	}
	return
}
