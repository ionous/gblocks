package dom

import "encoding/xml"

var names = struct {
	xml, block, shadow,
	value, statement, field xml.Name
}{
	xml.Name{Local: "xml"},
	xml.Name{Local: "block"},
	xml.Name{Local: "shadow"},
	xml.Name{Local: "value"},
	xml.Name{Local: "statement"},
	xml.Name{Local: "field"},
}

type Toolbox struct {
	Id         string     `xml:"id,attr,omitempty"`
	Style      string     `xml:"style,attr,omitempty"`
	Categories Categories `xml:"category,omitempty"`
	Blocks     Blocks     `xml:"block,omitempty"`
}

type Categories []*Category
type Blocks []*Block
type Items []Item
type Mutations []*Mutation

func (t *Toolbox) OuterHTML() string {
	output, err := xml.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(output)
}

type Category struct {
	Name       string     `xml:"name,attr"`
	Expanded   bool       `xml:"expanded,attr,omitempty"`
	Colour     string     `xml:"colour,attr,omitempty"`
	Categories Categories `xml:"category,omitempty"`
	Blocks     Blocks     `xml:"block,omitempty"`
}

func (t *Toolbox) MarshalToString() (ret string, err error) {
	if bytes, e := xml.Marshal(struct {
		XMLName xml.Name
		*Toolbox
	}{names.xml, t}); e != nil {
		err = e
	} else {
		ret = string(bytes)
	}
	return
}

func NewToolboxFromString(src string) (ret *Toolbox, err error) {
	t := new(Toolbox)
	if e := xml.Unmarshal([]byte(src), t); e != nil {
		err = e
	} else {
		ret = t
	}
	return
}
