package mutant

type MutatedInput struct {
	InputName string
	// mutation that created us.
	Arch Mutation
	// list atoms created by the user
	Atoms []*AtomizedInput
}

type AtomizedInput struct {
	Id   string // mui block guid
	Type string // type that was created
}

func (a *AtomizedInput) String() string {
	return a.Id + ":" + a.Type
}
