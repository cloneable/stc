package git

type Git interface {
	ParseBranchName(name string) (BranchName, error)
	CurrentBranch() (BranchName, error)
	CreateBranch(newBranch, baseBranch BranchName) error
	SwitchBranch(branch BranchName) error
	CreateSymref(name, target RefName, reason string) error
	CreateRef(name RefName, commit Commit) error
	GetRef(refName RefName) (Ref, error)
}

type BranchName string
