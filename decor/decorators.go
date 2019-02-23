package decor

import (
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	r "reflect"
)

type Decorators struct {
	// decor to function
	registry map[string]Fn
	// Container+Input to function
	// this allows us to find decorations while processing paths w.o reflection
	shortcut map[block.Type]Fn
}

func (d *Decorators) Register(name string, fn Fn) {
	if d.registry == nil {
		d.registry = make(map[string]Fn)
	}
	if _, exists := d.registry[name]; exists {
		panic("decoration already registered " + name)
	}
	d.registry[name] = fn
}

func (d *Decorators) Find(container block.Type, input block.Item) (ret Fn, okay bool) {
	decorName := block.SpecialType(container.String(), input.String())
	if fn, ok := d.shortcut[decorName]; ok {
		ret, okay = fn, true
	}
	return
}

// ex. `decor:"functionName"`
// FIX: implement + NextInput
func (d *Decorators) RegisterField(container block.Type, field r.StructField) (err error) {
	if decor, ok := field.Tag.Lookup("decor"); ok {
		if fn, ok := d.registry[decor]; !ok {
			err = errutil.New("unknown decoration", decor, "in", container.StructName(), field.Name)
		} else {
			in := block.ItemFromField(field)
			name := block.SpecialType(container.String(), in.String())
			d.shortcut[name] = fn
		}
	}
	return
}
