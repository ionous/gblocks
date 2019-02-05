package gblocks

type EnumPair [2]string // display, uniquifier

type RegisteredEnum struct {
	pairs []EnumPair
}

type RegisteredEnums map[TypeName]*RegisteredEnum
