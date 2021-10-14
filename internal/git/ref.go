package git

type RefName string

func (rn RefName) String() string { return string(rn) }

type Commit string

func (c Commit) String() string { return string(c) }

type Ref struct {
	Name   RefName
	Commit string
}
