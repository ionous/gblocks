package test

import "strconv"

type TestNames struct{ id int }

func (an *TestNames) GenerateUniqueName() string {
	x := strconv.Itoa(an.id)
	an.id++
	return "name" + x
}
