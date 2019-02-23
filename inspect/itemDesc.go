package inspect

import (
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/option"
)

// ItemDesc - fill the out dict with a description of this item suitable for block.
// the format doesnt appear to be documented anywhere; derived through examples.
func (it *Item) ItemDesc(scope block.Item, out block.Dict, pairs []EnumPair) (err error) {
	itemPath := scope.Push(it.Name)
	switch cls := it.Class; {
	case len(pairs) > 0 && cls != Int:
		err = errutil.New("enum expected int", itemPath, it)

	case cls == Statements:
		out.Insert(option.Name, itemPath.String())
		out.Insert(option.Type, block.StatementInput)

	case cls == Group:
		out.Insert(option.Name, itemPath.String())
		out.Insert(option.Type, block.DummyInput)
		out.Insert(option.Group, block.TypeFromStruct(it.Type))

	// input containing another block
	case cls == Input:
		out.Insert(option.Name, itemPath.String())
		out.Insert(option.Type, block.ValueInput)

	// a field of some sort ( ex. angle, checkbox, colour, date, dropdown, image, label, number, text, variable )
	case cls == Bool:
		out.Insert(option.Name, itemPath.String())
		out.Insert(option.Type, block.CheckboxField)

	case cls == Int:
		out.Insert(option.Name, itemPath.String())
		if len(pairs) > 0 {
			out.Insert(option.Type, block.DropdownField)
			out.Insert(option.Choices, pairs)
		} else {
			out.Insert(option.Type, block.NumberField)
			out.Insert(option.Precision, 1)
		}

	case cls == Uint:
		out.Insert(option.Name, itemPath.String())
		out.Insert(option.Type, block.NumberField)
		out.Insert(option.Precision, 1)
		out.Insert(option.Min, 0)

	case cls == Float:
		out.Insert(option.Name, itemPath.String())
		out.Insert(option.Type, block.NumberField)

	case cls == Label:
		out.Insert(option.Name, itemPath.String())
		out.Insert(option.Type, block.LabelField)
		out.Insert(option.Text, it.Name.Friendly())

	case cls == Text:
		out.Insert(option.Name, itemPath.String())
		out.Insert(option.Type, block.TextField) // "field_input"
		out.Insert(option.Text, it.Name.Friendly())

	case cls == Date:
		out.Insert(option.Name, itemPath.String())
		out.Insert(option.Type, block.DateField)

	default:
		err = errutil.New("invalid type", itemPath, it)
		break
	}
	return
}
