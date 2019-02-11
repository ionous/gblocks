package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
)

type XmlElement struct {
	*js.Object
	TagName    string      `js:"tagName"`
	Attributes *Attributes `js:"attributes"`
}

// an NamedNodeMap
type Attributes struct {
	*js.Object
}

type Attribute struct {
	*js.Object
	Name  string     `js:name`
	Value *js.Object `js:value`
}

// https://developer.mozilla.org/en-US/docs/Web/API/HTMLCollection
type HtmlCollection struct {
	*js.Object
}

var _xmlSource *XmlDoc

type XmlDoc struct {
	*js.Object
}

func (x *XmlDoc) createElement(tag string) *XmlElement {
	obj := x.Call("createElement", tag)
	return &XmlElement{Object: obj}
}

func xmlSource() *XmlDoc {
	if _xmlSource == nil {
		obj := js.Global.Get("document").Get("implementation").Call("createDocument", "", "", nil)
		_xmlSource = &XmlDoc{Object: obj}
	}
	return _xmlSource
}

func NewXmlElement(name string, attrs ...Attrs) *XmlElement {
	dom := xmlSource().createElement(name)
	for _, attrs := range attrs {
		for k, v := range attrs {
			dom.SetAttribute(k, v)
		}
	}
	return dom
}

func (m *XmlElement) AppendChild(child *XmlElement) *XmlElement {
	m.Call("appendChild", child)
	return child
}

func (m *XmlElement) GetAttribute(k string) (ret *Attribute) {
	if obj := m.Call("getAttribute", k); obj.Bool() {
		ret = &Attribute{Object: obj}
	}
	return
}

func (m *XmlElement) FirstElementChild() (ret *XmlElement) {
	if obj := m.Get("firstElementChild"); obj.Bool() {
		ret = &XmlElement{Object: obj}
	}
	return
}

func (m *XmlElement) Children() *HtmlCollection {
	obj := m.Get("children")
	return &HtmlCollection{Object: obj}
}

func (m *XmlElement) OuterHTML() string {
	return m.Get("outerHTML").String()
}

func (m *XmlElement) SetAttribute(k string, v interface{}) {
	m.Call("setAttribute", k, v)
}

func (m *XmlElement) SetInnerHTML(text string) {
	m.Set("innerHTML", text)
}

//HTMLCollection
func (m *XmlElement) GetElementsByTagName(tagName string) (ret *HtmlCollection) {
	if obj := m.Call("getElementsByTagName", tagName); obj != nil && obj.Bool() {
		ret = &HtmlCollection{Object: obj}
	}
	return
}

func (a *HtmlCollection) Num() int {
	return a.Get("length").Int()
}

func (a *HtmlCollection) Index(i int) (ret *XmlElement) {
	obj := a.Call("item", i)
	return &XmlElement{Object: obj}
}

func (a *Attributes) Num() int {
	return a.Get("length").Int()
}

func (a *Attributes) Index(i int) (ret *Attribute) {
	obj := a.Call("item", i)
	return &Attribute{Object: obj}
}

func (a *Attributes) Int(i int) (ret int) {
	if attr := a.Call("item", i); attr.Bool() {
		ret = attr.Get("value").Int()
	}
	return
}

func (a *Attributes) String(i int) (ret string) {
	if attr := a.Call("item", i); attr.Bool() {
		ret = attr.Get("value").String()
	}
	return
}
