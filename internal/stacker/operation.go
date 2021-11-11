package stacker

import (
	"fmt"

	"github.com/cloneable/stacker/internal/git"
)

type branch struct {
	branch git.BranchName
}

const (
	stackerRefPrefix      = "refs/stacker/"
	stackerBaseRefPrefix  = stackerRefPrefix + "base/"
	stackerStartRefPrefix = stackerRefPrefix + "start/"
	branchRefPrefix       = "refs/heads/"
)

func (b *branch) baseRefName() git.RefName {
	return git.RefName(stackerBaseRefPrefix + b.branch)
}

func (b *branch) startRefName() git.RefName {
	return git.RefName(stackerStartRefPrefix + b.branch)
}

func (b *branch) refName() git.RefName {
	return git.RefName(branchRefPrefix + b.branch)
}

type operation struct {
	git git.Git

	failed bool
	err    error
}

func op(git git.Git) *operation {
	return &operation{
		git: git,
	}
}

func (o *operation) Err() error {
	if o != nil {
		return o.err
	}
	return nil
}

func (o *operation) parseBranchName(name string) git.BranchName {
	if o == nil || o.failed {
		return ""
	}
	bname, err := o.git.ParseBranchName(name)
	if err != nil {
		o.failed = true
		o.err = fmt.Errorf("ParseBranchName: %w", err)
		return ""
	}
	return bname
}

// func (o *operation) readBranch(name string) *branch {
// 	if o == nil || o.failed {
// 		return nil
// 	}
// 	bname, err := o.git.ParseBranchName(name)
// 	if err != nil {
// 		o.failed = true
// 		o.err = fmt.Errorf("ParseBranchName: %w", err)
// 		return nil
// 	}
// 	return bname
// }

func (o *operation) createBranch(name git.BranchName, base *branch) *branch {
	if o == nil || o.failed {
		return nil
	}
	err := o.git.CreateBranch(name, base.branch)
	if err != nil {
		o.failed = true
		o.err = fmt.Errorf("CreateBranch: %w", err)
		return nil
	}
	return &branch{
		branch: name,
	}
}

func (o *operation) currentBranch() *branch {
	if o == nil || o.failed {
		return nil
	}
	name, err := o.git.CurrentBranch()
	if err != nil {
		o.failed = true
		o.err = fmt.Errorf("CurrentBranch: %w", err)
		return nil
	}
	return &branch{
		branch: name,
	}
}

func (o *operation) switchBranch(b *branch) {
	if o == nil || o.failed {
		return
	}
	err := o.git.SwitchBranch(b.branch)
	if err != nil {
		o.failed = true
		o.err = fmt.Errorf("SwitchBranch: %w", err)
		return
	}
}

func (o *operation) createSymref(branch, refBranch *branch, reason string) {
	if o == nil || o.failed {
		return
	}
	err := o.git.CreateSymref(branch.baseRefName(), refBranch.refName(), reason)
	if err != nil {
		o.failed = true
		o.err = fmt.Errorf("CreateSymref: %w", err)
		return
	}
}

func (o *operation) createRef(branch *branch, commit git.Commit) {
	if o == nil || o.failed {
		return
	}
	err := o.git.CreateRef(branch.startRefName(), commit)
	if err != nil {
		o.failed = true
		o.err = fmt.Errorf("CreateRef: %w", err)
		return
	}
}

func (o *operation) getRef(branch *branch) git.Ref {
	if o == nil || o.failed {
		return git.Ref{}
	}
	ref, err := o.git.GetRef(branch.refName())
	if err != nil {
		o.failed = true
		o.err = fmt.Errorf("GetRef: %w", err)
		return git.Ref{}
	}
	return ref
}
