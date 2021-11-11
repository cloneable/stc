package stacker

import (
	"context"
	"errors"

	"github.com/cloneable/stacker/internal/git"
)

var errUnimplemented = errors.New("unimplemented")

type Stacker struct {
	git git.Git
}

func New(repoPath string) *Stacker {
	return &Stacker{
		git: &git.Runner{
			Env:           nil,
			WorkDir:       repoPath,
			PrintCommands: true,
		},
	}
}

func (s *Stacker) Init(ctx context.Context, force bool) error {
	op := op(s.git)
	op.configAdd("transfer.hideRefs", stackerRefPrefix)
	op.configAdd("log.excludeDecoration", stackerRefPrefix)

	// TODO: read refs, branches, remotes
	// TODO: validate stacker refs against branches
	// TODO: determine list of needed refs
	// TODO: print and create list of created refs
	return op.Err()
}

func (s *Stacker) Clean(ctx context.Context, force bool, branches ...string) error {
	op := op(s.git)
	op.configUnsetPattern("transfer.hideRefs", stackerRefPrefix)
	op.configUnsetPattern("log.excludeDecoration", stackerRefPrefix)
	for _, ref := range op.listStackerRefs() {
		op.deleteRef(ref.Name, ref.Commit)
	}
	// TODO: for each branch
	// TODO: ... check if fully merged
	// TODO: ... check if remote ref == local branch
	// TODO: ... delete stacker refs
	// TODO: ... or print warning
	return op.Err()
}

func (s *Stacker) Show(ctx context.Context) error {
	// TODO: list all stacker tracked branches with a status
	return errUnimplemented
}

func (s *Stacker) Start(ctx context.Context, name string) error {
	op := op(s.git)
	baseB := op.currentBranch()
	newName := op.parseBranchName(name)
	newB := op.createBranch(newName, baseB)
	op.switchBranch(newB)
	op.createSymref(newB, baseB, "stacker: mark base branch")
	ref := op.getRef(baseB)
	op.createRef(newB, ref.Commit)
	return op.Err()
}

func (s *Stacker) Publish(ctx context.Context, branches ...string) error {
	// TODO: for each branch
	// TODO: ... determine state (already pushed?)
	// TODO: ... determine upstream
	// TODO: ... push branch to remote
	return errUnimplemented
}

func (s *Stacker) Delete(ctx context.Context, branch string) error {
	return errUnimplemented
}

func (s *Stacker) Rebase(ctx context.Context, branches ...string) error {
	// op := op(s.git)

	// var bs []*branch
	// if len(branches) == 0 {
	// 	bs = append(bs, op.currentBranch())
	// } else {
	// 	for _, name := range branches {
	// 		bs = append(bs, op.parseBranchName(name))
	// 	}
	// }

	// baseB := op.currentBranch()
	// newName := op.parseBranchName(name)
	// newB := op.createBranch(newName, baseB)
	// op.switchBranch(newB)
	// op.createSymref(newB, baseB, "stacker: mark base branch")
	// ref := op.getRef(baseB)
	// op.createRef(newB, ref.Commit)
	// return op.Err()

	// TODO: if len(branch) == 0 use current head as branch (head must be branch head)
	// TODO: for each branch
	// TODO: ... determine list of all stacked branches
	// TODO: ... add unselected branches to list
	// TODO: topologically sort selected branches
	// TODO: for each selected branch
	// TODO: ... call git rebase --onto
	// TODO: ... update stacker ref
	return errUnimplemented
}

func (s *Stacker) Sync(ctx context.Context, branches ...string) error {
	// TODO: if len(branch) == 0 use current head as branch (head must be branch head)
	// TODO: for each branch
	// TODO: ... determine list of all stacked branches
	// TODO: ... add unselected branches to list
	return errUnimplemented
}
