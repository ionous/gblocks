package tin

import (
	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/mutant"
	"github.com/ionous/gblocks/option"
	"github.com/ionous/gblocks/pascal"
)

type context struct {
	mutant.Atomizer
	mutations Mutations
}

// create arg descriptions for the passed pointer type
func (c *context) buildItems(scope string, ptrType r.Type, out *mutant.BlockMutations) (ret block.Args, err error) {
	var args block.Args
	structType := ptrType.Elem()
	// a field (ex. enum) followed by a mutation will vanish;
	// collapsing into the invisible dummy input used for tracking mutations.
	// we need to flush those fields into a separate visible dummy input.
	// ( or stop the mutation from hiding, but that needs more data in Mutation/s  )
	var standaloneFields int
	for i, cnt := 0, structType.NumField(); i < cnt; i++ {
		if field := structType.Field(i); len(field.PkgPath) == 0 {
			if field.Name != block.NextStatement {
				name := pascal.ToCaps(field.Name)
				if len(scope) > 0 {
					// ex. "a, FIELD"
					name = block.Scope(scope, name)
				}
				if desc, e := c.itemDesc(name, &field, out); e != nil {
					err = errutil.Append(err, e)
				} else if len(desc) > 0 {
					switch desc[option.Type] {
					case block.StatementInput, block.ValueInput:
						standaloneFields = 0
						break
					case block.DummyInput:
						if standaloneFields > 0 {
							// inject a dummy input before the new pending input to separate fields
							visibleDummy := block.Dict{option.Type: block.DummyInput}
							args.AddArg(visibleDummy)
						}
						break
					default:
						standaloneFields++
					}
					args.AddArg(desc)
					// atoms that are fields need to be encapsulated ( pushed down into ) a dummy input
					// otherwise they will get their own blank input, and RemoveAtoms will fail.
					// ( resulting in expansion of each mui input block into 4x workspace fields )
					if standaloneFields > 0 && len(scope) > 0 {
						visibleDummy := block.Dict{option.Name: name, option.Type: block.DummyInput}
						args.AddArg(visibleDummy)
						standaloneFields = 0
					}
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
func (c *context) itemDesc(name string, field *r.StructField, outMutations *mutant.BlockMutations) (ret block.Dict, err error) {
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
		if choices := c.GetPairs(field.Type.Name()); choices != nil {
			block.Merge(outDesc, tags, option.Type, block.DropdownField)
			block.Merge(outDesc, tags, option.Choices, choices)
		} else if _, ok := outDesc["readOnly"]; ok {
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
		inputType, _ := option.InputOption(tags)
		switch inputType {
		case block.ValueInput:
			limits = c.GetTermsByType(field.Type)

		case block.StatementInput:
			limits = c.GetStatementsByType(field.Type)

		case block.DummyInput:
			mutationName := pascal.ToUnderscore(field.Type.Elem().Name())
			if outMutations == nil {
				// typically a mutation in a mutation
				err = errutil.Fmt("invalid context for mutation %q", mutationName)
			} else if min, ok := c.mutations.GetMutationInfo(mutationName); !ok {
				err = errutil.Fmt("couldnt find mutation named  %q", mutationName)
			} else {
				outMutations.AddMutation(name, min)
				// note: we dont expand the fixed mutation b/c we'd need to know block id.
				// instead we use atoms, injecting when we generate a toolbox; load from xml.
			}
		}
		if err == nil {
			block.Merge(outDesc, tags, option.Name, name)
			block.Merge(outDesc, nil, option.Type, inputType) // pass nil to ignore original value
			if !limits.IsUnlimited() {
				block.Merge(outDesc, tags, option.Check, limits.Check())
			}
		}

	default:
		err = errutil.New("invalid type", field.Name, field.Type.Name(), cls.String())
		break
	}
	if err == nil {
		ret = outDesc
	}
	return
}
