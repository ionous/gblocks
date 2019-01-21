package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
)

type DomElement struct {
	*js.Object
	TagName    string      `js:tagName`
	Attributes *Attributes `js:attributes`
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

func NewDomElement(name string, attrs ...Attrs) *DomElement {
	obj := js.Global.Get("document").Call("createElement", name)
	dom := &DomElement{Object: obj}
	for _, attrs := range attrs {
		for k, v := range attrs {
			dom.SetAttribute(k, v)
		}
	}
	return dom
}

func (m *DomElement) AppendChild(child *DomElement) *DomElement {
	m.Call("appendChild", child)
	return child
}

func (m *DomElement) GetAttribute(k string) (ret *Attribute) {
	if obj := m.Call("getAttribute", k); obj.Bool() {
		ret = &Attribute{Object: obj}
	}
	return
}

func (m *DomElement) Children() *HtmlCollection {
	obj := m.Get("children")
	return &HtmlCollection{Object: obj}
}

func (m *DomElement) OuterHTML() string {
	return m.Get("outerHTML").String()
}

func (m *DomElement) SetAttribute(k string, v interface{}) {
	m.Call("setAttribute", k, v)
}

func (m *DomElement) SetInnerHTML(text string) {
	m.Set("innerHTML", text)
}

func (a *HtmlCollection) Num() int {
	return a.Get("length").Int()
}

func (a *HtmlCollection) Index(i int) (ret *DomElement) {
	obj := a.Call("item", i)
	return &DomElement{Object: obj}
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
