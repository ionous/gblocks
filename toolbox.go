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

// shared with toolbox creation
func appendMutation(name string, slice r.Value, out *XmlElement) {
	if cnt := slice.Len(); cnt > 0 {
		elements := make([]int, cnt)
		nameToIndex := make(map[r.Type]int)
		var typeNames []TypeName
		for i := 0; i < cnt; i++ {
			// note: we have to deref the element into its actual value before asking for its type.
			iface := slice.Index(i)
			ptr := iface.Elem()
			el := ptr.Elem()
			t := el.Type()
			typeIndex := len(nameToIndex)
			if existingIndex, ok := nameToIndex[t]; ok {
				typeIndex = existingIndex
			} else {
				nameToIndex[t] = typeIndex
				typeName := toTypeName(t)
				typeNames = append(typeNames, typeName)
			}
			elements[i] = typeIndex
		}
		el := out.AppendChild(NewXmlElement("data"))
		el.SetAttribute("name", name)
		el.SetAttribute("types", typeNames)
		el.SetAttribute("elements", elements)
	}
}

func toolboxData(v r.Value) *XmlElement {
	t := v.Type()
	n := toTypeName(t)
	el := NewXmlElement("block", Attrs{"type": n.String()})
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
							valEl := el.AppendChild(NewXmlElement("value", Attrs{"name": name}))
							valEl.AppendChild(toolboxData(unpack(nv)))
						}

					case r.Slice:
						if !nv.IsNil() {
							if _, ok := f.Tag.Lookup(tag_mutation); ok {
								if mutationEl == nil {
									mutationEl = el.AppendChild(NewXmlElement("mutation"))
								}
								appendMutation(name, nv, mutationEl)
							}
							top := el.AppendChild(NewXmlElement("statement", Attrs{"name": name}))
							next := false
							for i, cnt := 0, nv.Len(); i < cnt; i++ {
								if next {
									top = top.AppendChild(NewXmlElement("next"))
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

func toolboxField(name, val string) *XmlElement {
	field := NewXmlElement("field", Attrs{"name": name})
	field.SetInnerHTML(val)
	return field
}
