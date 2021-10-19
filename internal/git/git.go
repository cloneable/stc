package git

type Git interface {
	Bare() (bool, error)
	RepoRoot() (string, error)
	ValidBranchName(name string) (bool, error)
	CurrentBranch() (string, error)
	CreateBranch(branch, baseBranch string) error
	SwitchBranch(branch string) error
}
