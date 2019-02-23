package dom

import (
	"github.com/gopherjs/gopherjs/js"
)

type Attrs map[string]string

type Element struct {
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

func (x *XmlDoc) createElement(tag string) *Element {
	obj := x.Call("createElement", tag)
	return &Element{Object: obj}
}

func xmlSource() *XmlDoc {
	if _xmlSource == nil {
		obj := js.Global.Get("document").Get("implementation").Call("createDocument", "", "", nil)
		_xmlSource = &XmlDoc{Object: obj}
	}
	return _xmlSource
}

func NewElement(name string, attrs ...Attrs) *Element {
	dom := xmlSource().createElement(name)
	for _, attrs := range attrs {
		for k, v := range attrs {
			dom.SetAttribute(k, v)
		}
	}
	return dom
}

func (m *Element) AppendChild(child *Element) *Element {
	m.Call("appendChild", child)
	return child
}

func (m *Element) GetAttribute(k string) (ret *Attribute) {
	if obj := m.Call("getAttribute", k); obj.Bool() {
		ret = &Attribute{Object: obj}
	}
	return
}

func (m *Element) FirstElementChild() (ret *Element) {
	if obj := m.Get("firstElementChild"); obj.Bool() {
		ret = &Element{Object: obj}
	}
	return
}

func (m *Element) Children() *HtmlCollection {
	obj := m.Get("children")
	return &HtmlCollection{Object: obj}
}

func (m *Element) OuterHTML() string {
	return m.Get("outerHTML").String()
}

func (m *Element) SetAttribute(k string, v interface{}) {
	m.Call("setAttribute", k, v)
}

func (m *Element) SetInnerHTML(text string) {
	m.Set("innerHTML", text)
}

//HTMLCollection
func (m *Element) GetElementsByTagName(tagName string) (ret *HtmlCollection) {
	if obj := m.Call("getElementsByTagName", tagName); obj != nil && obj.Bool() {
		ret = &HtmlCollection{Object: obj}
	}
	return
}

func (a *HtmlCollection) Num() int {
	return a.Get("length").Int()
}

func (a *HtmlCollection) Index(i int) (ret *Element) {
	obj := a.Call("item", i)
	return &Element{Object: obj}
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
