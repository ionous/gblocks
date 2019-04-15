package mock

import "github.com/ionous/gblocks/block"

type InputList struct {
	Inputs []block.Input
}

func (k *InputList) NumInputs() int {
	return len(k.Inputs)
}

func (k *InputList) Input(i int) block.Input {
	return k.Inputs[i]
}

func (k *InputList) SetInput(i int, in block.Input) {
	k.Inputs[i] = in
}

func (k *InputList) RemoveInput(name string) {
	if _, i := k.InputByName(name); i >= 0 {
		a := k.Inputs
		k.Inputs = append(a[:i], a[i+1:]...)
	}
}

func (k *InputList) InputByName(name string) (retInput block.Input, retIndex int) {
	var found bool
	for i, in := range k.Inputs {
		if in.InputName() == name {
			retInput, retIndex = in, i
			found = true
			break
		}
	}
	if !found {
		retIndex = -1
	}
	return
}

func (k *InputList) Interpolate(msg string, args []block.Dict) {
	for _, a := range args {
		n, nok := a["name"].(string)
		t, tok := a["type"].(string)
		if nok && tok {
			in := &MockInput{Name: n, Type: t}
			k.Inputs = append(k.Inputs, in)
		}
	}
}
