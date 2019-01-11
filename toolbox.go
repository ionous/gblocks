package gblocks

import (
	r "reflect"
	"strconv"
)

// https://developers.google.com/blockly/guides/configur e/web/toolbox
type Tools struct {
	reg *Registry
}

func NewTools(reg *Registry) *Tools {
	return &Tools{}
}

type Attrs map[string]string

func (t *Tools) Box(tag string, attrs map[string]string, content ...interface{}) *DomElement {
	out := NewDomElement(tag)
	for k, v := range attrs {
		out.SetAttribute(k, v)
	}
	for _, c := range content {
		if child, ok := c.(*DomElement); ok {
			out.AppendChild(child)
		} else {
			// must be a custom block type
		}
	}
	return out
}

func NewToolData(el *DomElement, src ...interface{}) *DomElement {
	for _, n := range src {
		v := r.ValueOf(n).Elem()
		el.AppendChild(toolboxData(v))
	}
	return el
}

func unptr(v r.Value) (ret r.Value) {
	switch k := v.Kind(); k {
	case r.Ptr:
		ret = v.Elem()
	default:
		ret = v
	}
	return
}

func toolboxData(v r.Value) *DomElement {
	t := v.Type()
	n := toTypeName(t)
	el := CreateDomElement("block", Attrs{"type": n.String()})
	var mutationEl *DomElement
	//
	for i, cnt := 0, t.NumField(); i < cnt; i++ {
		// skip unexpected symbols ( only unexported symbols have a pkg path )
		if f := t.Field(i); len(f.PkgPath) == 0 {
			switch f.Name {
			case PreviousStatementField:
				// skip
			case NextStatementField:
				// <next>, recursive
				if nv := v.FieldByIndex(f.Index); !nv.IsNil() {
					nextEl := el.AppendChild(NewDomElement("next"))
					nextEl.AppendChild(toolboxData(unptr(nv.Elem())))
				}

			default:
				name := pascalToCaps(f.Name)
				nv := v.FieldByIndex(f.Index)

				switch k := f.Type.Kind(); k {
				case r.Bool:
					field := toolboxField(name, strconv.FormatBool(nv.Bool()))
					el.AppendChild(field)

				case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
					field := toolboxField(name, strconv.FormatInt(nv.Int(), 10))
					el.AppendChild(field)

				case r.Uint, r.Uint8, r.Uint16, r.Uint32, r.Uint64:
					field := toolboxField(name, strconv.FormatUint(nv.Uint(), 19))
					el.AppendChild(field)

				case r.Float32, r.Float64:
					field := toolboxField(name, strconv.FormatFloat(nv.Float(), 'g', -1, 32))
					el.AppendChild(field)

				// input containing another block
				case r.Ptr:
				case r.Interface:
					if !nv.IsNil() {
						valEl := el.AppendChild(CreateDomElement("value", Attrs{"name": name}))
						valEl.AppendChild(toolboxData(unptr(nv.Elem())))
					}

				case r.Slice:
					if !nv.IsNil() {
						if _, ok := f.Tag.Lookup("mutation"); ok {
							if mutationEl == nil {
								mutationEl = el.AppendChild(NewDomElement("mutation"))
							}
							appendMutation(name, nv, mutationEl)
						} else {
							top := el.AppendChild(CreateDomElement("statement", Attrs{"name": name}))
							next := false
							for i, cnt := 0, nv.Len(); i < cnt; i++ {
								if next {
									top = top.AppendChild(NewDomElement("next"))
								}
								top = top.AppendChild(toolboxData(unptr(nv.Index(i).Elem())))
								next = true
							}
						}
					}

				default:
					if str := nv.String(); len(str) > 0 {
						field := toolboxField(name, str)
						el.AppendChild(field)
					}
				}
			}
		}
	}
	return el
}

func toolboxField(name, val string) *DomElement {
	field := CreateDomElement("field", Attrs{"name": name})
	field.SetInnerHTML(val)
	return field
}
