package dom

import "encoding/xml"

var BlockName xml.Name = xml.Name{Local: "block"}
var ShadowName xml.Name = xml.Name{Local: "shadow"}

type Toolbox struct {
	XMLName    xml.Name    `xml:"xml"`
	Id         string      `xml:"id,attr,omitempty"`
	Style      string      `xml:"style,attr,omitempty"`
	Categories []*Category `xml:"category,omitempty"`
	Blocks     []*Block    `xml:",omitempty"`
}

func (t *Toolbox) OuterHTML() string {
	output, err := xml.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(output)
}

type Category struct {
	Name       string      `xml:"name,attr"`
	Expanded   bool        `xml:"expanded,attr,omitempty"`
	Colour     string      `xml:"colour,attr,omitempty"`
	Categories []*Category `xml:"category,omitempty"`
	Blocks     []*Block    `xml:",omitempty"`
}

type Block struct {
	XMLName   xml.Name     // block or shadow
	Type      string       `xml:"type,attr"`
	Mutations *[]*Mutation `xml:"mutation>input,omitempty"`
	Items     []Item       `xml:",omitempty"` // values or fields
	Next      *Block       `xml:"next>XXX,omitempty"`
}

type Item interface{ Item() Item }

func (b *Block) AppendItem(it Item) {
	b.Items = append(b.Items, it)
}

func (b *Block) AppendMutation(m *Mutation) {
	(*b.Mutations) = append((*b.Mutations), m)
}

type Value struct {
	XMLName xml.Name `xml:"value"`
	Name    string   `xml:"name,attr"`
	Block   *Block   `xml:",omitempty"`
}

func (it *Value) Item() Item { return it }

type Statement struct {
	XMLName xml.Name `xml:"statement"`
	Name    string   `xml:"name,attr"`
	Block   *Block   `xml:",omitempty"`
}

func (it *Statement) Item() Item { return it }

type Field struct {
	XMLName xml.Name `xml:"field"`
	Name    string   `xml:"name,attr"`
	Content string   `xml:",innerxml"`
}

func (it *Field) Item() Item { return it }

type Mutation struct {
	Name          string  `xml:"name,attr"`
	MutableInputs []*Atom `xml:"atom,omitempty"`
}

func (m *Mutation) AppendAtom(name string) {
	m.MutableInputs = append(m.MutableInputs, &Atom{name})
}

type Atom struct {
	Type string `xml:"type,attr"`
}
