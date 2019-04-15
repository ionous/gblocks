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
	Categories Categories `xml:"category,omitempty"`
	Blocks     Blocks     `xml:"block,omitempty"`
}

type Categories []*Category
type Blocks []*Block

type Category struct {
	Name       string     `xml:"name,attr"`
	Expanded   bool       `xml:"expanded,attr,omitempty"`
	Colour     string     `xml:"colour,attr,omitempty"`
	Categories Categories `xml:"category,omitempty"`
	Blocks     BlockList  `xml:"block,omitempty"`
}

func (t *Toolbox) IndentHTML(indent string) (ret string) {
	container := struct {
		XMLName xml.Name
		*Toolbox
	}{names.xml, t}
	var bytes []byte
	var e error
	if len(indent) > 0 {
		bytes, e = xml.MarshalIndent(container, "", indent)
	} else {
		bytes, e = xml.Marshal(container)
	}
	if e != nil {
		panic(e)
	} else {
		ret = string(bytes)
	}
	return
}

// return <xml></xml>
func (t *Toolbox) OuterHTML() string {
	return t.IndentHTML("")
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
