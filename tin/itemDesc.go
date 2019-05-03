package tin

import (
	r "reflect"
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/mutant"
	"github.com/ionous/gblocks/option"
	"github.com/ionous/gblocks/pascal"
)

type context struct {
	mutant.Atomizer
	mutables Mutables
}

func (c *context) addMutationInput(inputName, mutationName string, out *mutant.InMutations) (ret *Mutable, err error) {
	if out == nil {
		// typically a mutation in a mutation
		err = errutil.New("invalid context for mutation", mutationName)
	} else if m, ok := c.mutables.FindMutable(mutationName); !ok {
		err = errutil.New("couldnt find mutation named", mutationName)
	} else {
		if out.Mutations == nil {
			out.Mutations = make(map[string]mutant.InMutation)
		}
		out.Inputs = append(out.Inputs, inputName)
		out.Mutations[inputName] = m
		ret = m
	}
	return
}

// create arg descriptions for the passed pointer type
func (c *context) buildItems(scope string, ptrType r.Type, out *mutant.InMutations) (ret block.Args, err error) {
	var args block.Args
	structType := ptrType.Elem()
	for i, cnt := 0, structType.NumField(); i < cnt; i++ {
		if field := structType.Field(i); len(field.PkgPath) == 0 {
			if field.Name != block.NextStatement {
				name := pascal.ToCaps(field.Name)
				if len(scope) > 0 {
					// ex. a$ muiBlockId $ FIELD
					name = block.Scope(scope, name)
				}
				if desc, e := c.itemDesc(name, &field, out); e != nil {
					err = errutil.Append(err, e)
				} else if len(desc) > 0 {
					args.AddArg(desc)
				}
			}
		}
	}
	if err == nil {
		ret = args
	}
	return
}

// desc for now -- could return an item element like quark
func (c *context) itemDesc(name string, field *r.StructField, outMutations *mutant.InMutations) (ret block.Dict, err error) {
	outDesc := make(block.Dict)
	tags := parseTags(string(field.Tag))
	switch cls := Classify(field.Type); cls {
	case Option:
		// skip for now

	// a field of some sort ( ex. angle, checkbox, colour, date, dropdown, image, label, number, text, variable )
	case Bool:
		block.Merge(outDesc, tags, option.Name, name)
		block.Merge(outDesc, tags, option.Type, block.CheckboxField)

	case Int:
		block.Merge(outDesc, tags, option.Name, name)
		if choices := c.GetPairs(field.Type.Name()); choices != nil {
			block.Merge(outDesc, tags, option.Type, block.DropdownField)
			block.Merge(outDesc, tags, option.Choices, choices)
		} else {
			block.Merge(outDesc, tags, option.Type, block.NumberField)
			block.Merge(outDesc, tags, option.Precision, 1)
		}

	case Uint:
		block.Merge(outDesc, tags, option.Name, name)
		block.Merge(outDesc, tags, option.Type, block.NumberField)
		block.Merge(outDesc, tags, option.Precision, 1)
		block.Merge(outDesc, tags, option.Min, 0)

	case Float:
		block.Merge(outDesc, tags, option.Name, name)
		block.Merge(outDesc, tags, option.Type, block.NumberField)

	case Text:
		block.Merge(outDesc, tags, option.Name, name)
		if _, ok := outDesc["readOnly"]; ok {
			block.Merge(outDesc, tags, option.Type, block.LabelField)
			block.Merge(outDesc, tags, option.Text, pascal.ToSpaces(field.Name))
		} else {
			block.Merge(outDesc, tags, option.Type, block.TextField) // "field_input"
			block.Merge(outDesc, tags, option.Text, pascal.ToSpaces(field.Name))
		}

	case Date:
		block.Merge(outDesc, tags, option.Name, name)
		block.Merge(outDesc, tags, option.Type, block.DateField)

	// input pointing to another block
	case InputClass:
		var limits block.Limits
		inputType := inputOption(tags)
		switch inputType {
		case block.ValueInput:
			limits = c.GetTermsByType(field.Type)

		case block.StatementInput:
			limits = c.GetStatementsByType(field.Type)

		case block.DummyInput:
			targetType := pascal.ToUnderscore(field.Type.Elem().Name())
			if _, e := c.addMutationInput(name, targetType, outMutations); e != nil {
				err = e
			} else {
				// note: we dont expand the fixed mutation b/c we'd need to know block id.
				// instead we use atoms, injecting when we generate a toolbox; load from xml.
			}
		}
		if err == nil {
			block.Merge(outDesc, tags, option.Name, name)
			block.Merge(outDesc, tags, option.Type, inputType)
			if !limits.IsUnlimited() {
				block.Merge(outDesc, tags, option.Check, limits.Check())
			}
		}

	default:
		err = errutil.New("invalid type", name, field.Type.Name())
		break
	}
	if err == nil {
		ret = outDesc
	}
	return
}

// fix? if really needed to decouple this, could have an "input" factory
func inputOption(out block.Dict) string {
	opt, _ := out[option.Input].(string)
	lower := strings.ToLower(opt)
	switch lower {
	case "":
		lower = "value"
	case "mutation":
		lower = "dummy"
	}
	return "input_" + lower
}
