package gblocks

import (
	r "reflect"
	"strconv"
)

type Attrs map[string]string

// Toolbox - create a dom element from the passed tag and attrs, and attach the passed content.
// content can include dom nodes or gblocks instance data.
// returns the parent node.
// see also: https://developers.google.com/blockly/guides/configure/web/toolbox
func Toolbox(tag string, attrs map[string]string, content ...interface{}) *DomElement {
	out := NewDomElement(tag, attrs)
	return NewTools(out, content...)
}

// NewTools - attach toolbox content to the passed parent.
// see also Toolbox.
func NewTools(parent *DomElement, content ...interface{}) *DomElement {
	for _, c := range content {
		if child, ok := c.(*DomElement); ok {
			parent.AppendChild(child)
		} else {
			v := r.ValueOf(c).Elem()
			parent.AppendChild(toolboxData(v))
		}
	}
	return parent
}

func unpack(v r.Value) (ret r.Value) {
	switch k := v.Kind(); k {
	case r.Ptr:
		ret = v.Elem()
	case r.Interface:
		ret = unpack(v.Elem())
	default:
		ret = v
	}
	return
}

func toolboxData(v r.Value) *DomElement {
	t := v.Type()
	n := toTypeName(t)
	el := NewDomElement("block", Attrs{"type": n.String()})
	var mutationEl *DomElement
	//
	for i, cnt := 0, t.NumField(); i < cnt; i++ {
		// skip unexpected symbols ( only unexported symbols have a pkg path )
		if f := t.Field(i); len(f.PkgPath) == 0 {
			switch f.Name {
			case PreviousField:
				// skip
			case NextField:
				// <next>, recursive
				if nv := v.FieldByIndex(f.Index); !nv.IsNil() {
					nextEl := el.AppendChild(NewDomElement("next"))
					nextEl.AppendChild(toolboxData(unpack(nv)))
				}

			default:
				name := pascalToCaps(f.Name)
				nv := v.FieldByIndex(f.Index)

				// see if the type implements the stringer, for instance an enum.
				type stringer interface{ String() string }
				if str, ok := nv.Interface().(stringer); ok {
					el.AppendChild(toolboxField(name, str.String()))
				} else {
					switch k := f.Type.Kind(); k {
					case r.Bool:
						field := toolboxField(name, strconv.FormatBool(nv.Bool()))
						el.AppendChild(field)

					case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
						field := toolboxField(name, strconv.FormatInt(nv.Int(), 10))
						el.AppendChild(field)

					case r.Uint, r.Uint8, r.Uint16, r.Uint32, r.Uint64:
						field := toolboxField(name, strconv.FormatUint(nv.Uint(), 10))
						el.AppendChild(field)

					case r.Float32, r.Float64:
						field := toolboxField(name, strconv.FormatFloat(nv.Float(), 'g', -1, 32))
						el.AppendChild(field)

					// input containing another block
					case r.Ptr:
					case r.Interface:
						if !nv.IsNil() {
							valEl := el.AppendChild(NewDomElement("value", Attrs{"name": name}))
							valEl.AppendChild(toolboxData(unpack(nv)))
						}

					case r.Slice:
						if !nv.IsNil() {
							if _, ok := f.Tag.Lookup(tag_mutation); ok {
								if mutationEl == nil {
									mutationEl = el.AppendChild(NewDomElement("mutation"))
								}
								appendMutation(name, nv, mutationEl)
							}
							top := el.AppendChild(NewDomElement("statement", Attrs{"name": name}))
							next := false
							for i, cnt := 0, nv.Len(); i < cnt; i++ {
								if next {
									top = top.AppendChild(NewDomElement("next"))
								}
								top = top.AppendChild(toolboxData(unpack(nv.Index(i))))
								next = true
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
	}
	return el
}

func toolboxField(name, val string) *DomElement {
	field := NewDomElement("field", Attrs{"name": name})
	field.SetInnerHTML(val)
	return field
}
