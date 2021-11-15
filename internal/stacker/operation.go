package stacker

import (
	"fmt"
	"sort"

	"github.com/cloneable/stacker/internal/git"
)

type operation struct {
	git  git.Git
	repo git.Repository

	err error
}

func op(git git.Git) *operation {
	return &operation{
		git: git,
	}
}

func (o *operation) Failf(s string, args ...interface{}) error {
	o.err = fmt.Errorf(s, args...)
	return o.err
}

func (o *operation) Err() error {
	if o != nil && o.err != nil {
		o.git.DumpOutput()
		return o.err
	}
	return nil
}

func (o *operation) snapshot() {
	if o == nil || o.err != nil {
		return
	}
	repo, err := git.SnapshotRepository(o.git)
	if err != nil {
		o.err = fmt.Errorf("snapshot: %w", err)
		return
	}
	o.repo = repo
}

func (o *operation) head() git.BranchName {
	if o == nil || o.err != nil {
		return ""
	}
	head, ok := o.repo.Head()
	if !ok {
		o.err = fmt.Errorf("HEAD is unset")
		return ""
	}
	return head
}

func (o *operation) hasRef(name git.RefName) bool {
	if o == nil || o.err != nil {
		return false
	}
	_, found := o.repo.LookupRef(name)
	return found
}

func (o *operation) ref(name git.RefName) git.Ref {
	if o == nil || o.err != nil {
		return git.Ref{}
	}
	ref, ok := o.repo.LookupRef(name)
	if !ok {
		o.err = fmt.Errorf("ref %s not found", name)
		return git.Ref{}
	}
	return ref
}

func (o *operation) branch(name string) git.BranchName {
	if o == nil || o.err != nil {
		return ""
	}
	branchName, ok := o.repo.LookupBranch(name)
	if !ok {
		o.err = fmt.Errorf("ref %s not found", name)
		return ""
	}
	return branchName
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
		"--no-deref",
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
		"--no-deref",
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

func (o *operation) push(name git.BranchName, remote git.RemoteName, expect git.ObjectName) {
	if o == nil || o.err != nil {
		return
	}
	_, err := o.git.Exec(
		"push",
		"--set-upstream",
		fmt.Sprintf("--force-with-lease=%s:%s", name, expect),
		remote.String(),
		fmt.Sprintf("%s:%s", name, name),
	)
	if err != nil {
		o.err = fmt.Errorf("push: %w", err)
		return
	}
}

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

func (o *operation) fetchAllPrune() {
	if o == nil || o.err != nil {
		return
	}
	_, err := o.git.Exec(
		"fetch",
		"--all",
		"--prune",
	)
	if err != nil {
		o.err = fmt.Errorf("fetchAllPrune: %w", err)
		return
	}
}

func (o *operation) trackedBranches() []git.BranchName {
	if o == nil || o.err != nil {
		return nil
	}
	branchMap := make(map[git.BranchName]bool)
	for _, ref := range o.repo.AllRefs() {
		if ref.Stacker() {
			branchMap[ref.Name().BranchName()] = true
		}
	}
	branches := make([]git.BranchName, 0, len(branchMap))
	for branch := range branchMap {
		branches = append(branches, branch)
	}
	sort.Slice(branches, func(i, j int) bool { return branches[i] < branches[j] })
	return branches
}
