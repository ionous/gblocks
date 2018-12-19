package gblocks

import (
	r "reflect"
)

type BlockData struct {
	elem r.Value
}

func (ws *Workspace) BlockData(b *Block) (ret BlockData) {
	if elem := ws.dataPointerById(b.Id); !elem.IsValid() {
		panic("Unknown block " + b.Id)
	} else {
		ret.elem = elem.Elem()
	}
	return
}

func (d *BlockData) Elem() r.Value {
	return d.elem
}

func (d *BlockData) Ptr() r.Value {
	return d.elem
}

func (d *BlockData) Elements(in *Input) (ret r.Value, okay bool) {
	if m, ok := d.Mutation(in); ok {
		ret, okay = m.Elements(), true
	}
	return
}

func (d *BlockData) Mutation(in *Input) (ret Mutation, okay bool) {
	if d.elem.IsValid() && in != nil && in.mutations != 0 {
		fieldName := underscoreToPascal(in.Name)
		if field := d.elem.FieldByName(fieldName); field.IsValid() {
			if mutation, ok := field.Addr().Interface().(Mutation); ok {
				ret, okay = mutation, true
			}
		}
	}
	return
}
