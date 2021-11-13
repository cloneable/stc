package stacker

import (
	"fmt"

	"github.com/cloneable/stacker/internal/git"
)

type operation struct {
	git git.Git

	err error
}

func op(git git.Git) *operation {
	return &operation{
		git: git,
	}
}

func (o *operation) Err() error {
	if o != nil {
		o.git.DumpOutput()
		return o.err
	}
	return nil
}

func (o *operation) snapshot() git.Repository {
	if o == nil || o.err != nil {
		return git.Repository{}
	}
	repo, err := git.SnapshotRepository(o.git)
	if err != nil {
		o.err = fmt.Errorf("snapshot: %w", err)
		return git.Repository{}
	}
	return repo
}

func (o *operation) parseBranchName(name string) git.BranchName {
	if o == nil || o.err != nil {
		return ""
	}
	_, err := o.git.Exec(
		"check-ref-format",
		"--branch",
		name,
	)
	if err != nil {
		o.err = fmt.Errorf("parseBranchName: %w", err)
		return ""
	}
	return git.BranchName(name)
}

func (o *operation) createBranch(name, base git.BranchName) {
	if o == nil || o.err != nil {
		return
	}
	_, err := o.git.Exec(
		"branch",
		"--create-reflog",
		name.String(),
		base.String(),
	)
	if err != nil {
		o.err = fmt.Errorf("createBranch: %w", err)
		return
	}
}

func (o *operation) switchBranch(b git.BranchName) {
	if o == nil || o.err != nil {
		return
	}
	_, err := o.git.Exec(
		"switch",
		"--no-guess",
		b.String(),
	)
	if err != nil {
		o.err = fmt.Errorf("switchBranch: %w", err)
		return
	}
}

func (o *operation) createSymref(name, target git.RefName, reason string) {
	if o == nil || o.err != nil {
		return
	}
	_, err := o.git.Exec(
		"symbolic-ref",
		"-m",
		reason,
		name.String(),
		target.String(),
	)
	if err != nil {
		o.err = fmt.Errorf("createSymref: %w", err)
		return
	}
}

func (o *operation) createRef(name git.RefName, commit git.ObjectName) {
	if o == nil || o.err != nil {
		return
	}
	_, err := o.git.Exec(
		"update-ref",
		"--create-reflog",
		name.String(),
		commit.String(),
		git.NonExistantObject.String(),
	)
	if err != nil {
		o.err = fmt.Errorf("createRef: %w", err)
		return
	}
}

func (o *operation) updateRef(name git.RefName, newCommit, curCommit git.ObjectName) {
	if o == nil || o.err != nil {
		return
	}
	_, err := o.git.Exec(
		"update-ref",
		"--create-reflog",
		name.String(),
		newCommit.String(),
		curCommit.String(),
	)
	if err != nil {
		o.err = fmt.Errorf("updateRef: %w", err)
		return
	}
}

func (o *operation) deleteRef(name git.RefName, curCommit git.ObjectName) {
	if o == nil || o.err != nil {
		return
	}
	_, err := o.git.Exec(
		"update-ref",
		"--no-deref",
		"-d",
		name.String(),
		curCommit.String(),
	)
	if err != nil {
		o.err = fmt.Errorf("deleteRef: %w", err)
		return
	}
}

func (o *operation) rebaseOnto(name git.BranchName) {
	if o == nil || o.err != nil {
		return
	}
	_, err := o.git.Exec(
		"rebase",
		"--committer-date-is-author-date",
		"--onto",
		name.StackerBaseRefName().String(),
		name.StackerStartRefName().String(),
		name.String(),
	)
	if err != nil {
		o.err = fmt.Errorf("rebaseOnto: %w", err)
		return
	}
}

func (o *operation) pushForce(b git.BranchName) {
	if o == nil || o.err != nil {
		return
	}
	_, err := o.git.Exec(
		"push",
		"--force-with-lease",
		// TODO: explicit parameters
		// fmt.Sprintf("--force-with-lease=%s:%s", branchName, expectedCommit),
		// fmt.Sprintf("refs/heads/%s:refs/remotes/%s/%s", branchName, remoteName, branchName),
	)
	if err != nil {
		o.err = fmt.Errorf("pushForce: %w", err)
		return
	}
}

// func (o *operation) pushUpstream(b git.BranchName) {
// 	if o == nil || o.err != nil  {
// 		return
// 	}
// 	_, err := o.git.Exec(
// 		"push",
// 		"--set-upstream",
// 		remote,
// 		string(remoteBranch),
// 	)
// 	if err != nil {
// 		o.err = fmt.Errorf("pushUpstream: %w", err)
// 		return
// 	}
// }

func (o *operation) configSet(key, value string) {
	if o == nil || o.err != nil {
		return
	}
	_, err := o.git.Exec(
		"config",
		"--local",
		key,
		value,
	)
	if err != nil {
		o.err = fmt.Errorf("configSet: %w", err)
		return
	}
}

func (o *operation) configAdd(key, value string) {
	if o == nil || o.err != nil {
		return
	}
	_, err := o.git.Exec(
		"config",
		"--local",
		"--add",
		key,
		value,
	)
	if err != nil {
		o.err = fmt.Errorf("configAdd: %w", err)
		return
	}
}

func (o *operation) configUnsetPattern(key, pattern string) {
	if o == nil || o.err != nil {
		return
	}
	res, err := o.git.Exec(
		"config",
		"--local",
		"--fixed-value",
		"--unset-all",
		key,
		pattern,
	)
	// 5 means the nothing matched.
	if err != nil && res.ExitCode != 5 {
		o.err = fmt.Errorf("configUnsetPattern: %w", err)
		return
	}
}
