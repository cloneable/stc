package git

type Git interface {
	Bare() (bool, error)
	RepoRoot() (string, error)
	ParseBranchName(name string) (BranchName, error)
	CurrentBranch() (BranchName, error)
	CreateBranch(newBranch, baseBranch BranchName) error
	SwitchBranch(branch BranchName) error
	PushUpstream(branch BranchName, remote string) error
}

type BranchName string
