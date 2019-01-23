package gblocks

import (
	r "reflect"
	"strconv"
	"strings"
)

type Attrs map[string]string

// Toolbox - create a dom element from the passed tag and attrs, and attach the passed content.
// content can include dom nodes or gblocks instance data.
// returns the parent node.
// see also: https://developers.google.com/blockly/guides/configure/web/toolbox
func Toolbox(tag string, attrs map[string]string, content ...interface{}) *XmlElement {
	out := NewXmlElement(tag, attrs)
	return NewTools(out, content...)
}

// NewTools - attach toolbox content to the passed parent.
// see also Toolbox.
func NewTools(parent *XmlElement, content ...interface{}) *XmlElement {
	for _, c := range content {
		if child, ok := c.(*XmlElement); ok {
			parent.AppendChild(child)
		} else {
			v := r.ValueOf(c).Elem()
			parent.AppendChild(toolboxBlock(v))
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

// shared with toolbox creation
// name is input name
func appendMutation(name string, slice r.Value, out *XmlElement) {
	if cnt := slice.Len(); cnt > 0 {
		atoms := out.AppendChild(NewXmlElement("atoms", Attrs{"name": name}))
		for i := 0; i < cnt; i++ {
			el := unpack(slice.Index(i))
			typeName := toTypeName(el.Type())
			atom := NewXmlElement("atom", Attrs{"type": typeName.String()})
			atoms.AppendChild(atom)
		}
	}
}

func toolboxBlock(v r.Value) *XmlElement {
	t := v.Type()
	n := toTypeName(t)
	el := NewXmlElement("block", Attrs{"type": n.String()})
	toolboxFields(el, "", v, t)
	return el
}

func toolboxFields(el *XmlElement, prefix string, v r.Value, t r.Type) {
	var mutationEl *XmlElement
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
					nextEl := el.AppendChild(NewXmlElement("next"))
					nextEl.AppendChild(toolboxBlock(unpack(nv)))
				}
			default:
				name := prefix + pascalToCaps(f.Name)
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
					case r.Ptr, r.Interface:
						// ATOM_INPUT 0 ptr
						// ATOM_FIELD 0 string
						// ATOM_INPUT 0 ptr
						if !nv.IsNil() {
							valEl := el.AppendChild(NewXmlElement("value", Attrs{"name": name}))
							valEl.AppendChild(toolboxBlock(unpack(nv)))
						}

					case r.Slice:
						if !nv.IsNil() {
							if _, ok := f.Tag.Lookup(tag_mutation); !ok {
								top := el.AppendChild(NewXmlElement("statement", Attrs{"name": name}))
								next := false
								for i, cnt := 0, nv.Len(); i < cnt; i++ {
									if next {
										top = top.AppendChild(NewXmlElement("next"))
									}
									top = top.AppendChild(toolboxBlock(unpack(nv.Index(i))))
									next = true
								}
							} else {
								if mutationEl == nil {
									mutationEl = el.AppendChild(NewXmlElement("mutation"))
								}
								appendMutation(name, nv, mutationEl)

								// we have to look at each atom to determine what to add
								for i, cnt := 0, nv.Len(); i < cnt; i++ {
									// prefix - ex. MUTANT/2/
									prefix := strings.Join([]string{name, strconv.Itoa(i), ""}, "/")
									elv := unpack(nv.Index(i))
									toolboxFields(el, prefix, elv, elv.Type())
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
	}
}

func toolboxField(name, val string) *XmlElement {
	field := NewXmlElement("field", Attrs{"name": name})
	field.SetInnerHTML(val)
	return field
}
