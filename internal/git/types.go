package git

import "path"

type Tag struct {
	name string
}

func (t Tag) RefName() RefName {
	return RefName(path.Join("refs/tags", t.name))
}

type Head struct {
	name string
}

func (h Head) RefName() RefName {
	return RefName(path.Join("refs/heads", h.name))
}

type Remote string

type Branch struct {
	head   Head
	remote Remote
}

type RemoteBranch struct {
}

type Label struct {
	ref Ref
}

type BranchName string
