package stacker

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloneable/stacker/internal/git"
)

const RefNamespace = "stacker"

var errUnimplemented = errors.New("unimplemented")

type Stacker struct {
	git git.Git
}

func New(repoPath string) *Stacker {
	return &Stacker{
		git: &git.Runner{
			Env:     nil,
			WorkDir: repoPath,
		},
	}
}

func (s *Stacker) Init(ctx context.Context, force bool) error {
	// TODO: read refs, branches, remotes
	// TODO: validate stacker refs against branches
	// TODO: determine list of needed refs
	// TODO: print and create list of created refs
	return errUnimplemented
}

func (s *Stacker) Show(ctx context.Context) error {
	// TODO: list all stacker tracked branches with a status
	return errUnimplemented
}

func (s *Stacker) Clean(ctx context.Context, force bool, branches ...string) error {
	// TODO: for each branch
	// TODO: ... check if fully merged
	// TODO: ... check if remote ref == local branch
	// TODO: ... delete stacker refs
	// TODO: ... or print warning
	return errUnimplemented
}

func (s *Stacker) Create(ctx context.Context, name string) error {
	cur, err := s.git.CurrentBranch()
	if err != nil {
		return fmt.Errorf("CurrentBranch: %w", err)
	}
	branch, err := s.git.ParseBranchName(name)
	if err != nil {
		return fmt.Errorf("ParseBranchName: %w", err)
	}
	if err := s.git.CreateBranch(branch, cur); err != nil {
		return fmt.Errorf("CreateBranch: %w", err)
	}
	if err := s.git.SwitchBranch(branch); err != nil {
		return fmt.Errorf("SwitchBranch: %w", err)
	}
	// TODO: determine base branch and its origin
	// TODO: create new branch off of base branch
	// TODO: add remote tracking
	// TODO: set stacker refs: base symref, start commit
	return nil
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
