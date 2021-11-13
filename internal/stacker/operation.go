package stacker

import (
	"bufio"
	"fmt"
	"strings"

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

	err error
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

// func (o *operation) readBranch(name string) *branch {
// 	if o == nil || o.err != nil  {
// 		return nil
// 	}
// 	bname, err := o.git.ParseBranchName(name)
// 	if err != nil {
// 		o.err = fmt.Errorf("ParseBranchName: %w", err)
// 		return nil
// 	}
// 	return bname
// }

func (o *operation) createBranch(name git.BranchName, base *branch) *branch {
	if o == nil || o.err != nil {
		return nil
	}
	_, err := o.git.Exec(
		"branch",
		string(name),
		string(base.branch),
	)
	if err != nil {
		o.err = fmt.Errorf("createBranch: %w", err)
		return nil
	}
	return &branch{
		branch: name,
	}
}

func (o *operation) currentBranch() *branch {
	if o == nil || o.err != nil {
		return nil
	}
	res, err := o.git.Exec(
		"branch",
		"--show-current",
	)
	if err != nil {
		o.err = fmt.Errorf("currentBranch: %w", err)
		return nil
	}
	return &branch{
		branch: git.BranchName(strings.TrimSpace(res.Stdout.String())),
	}
}

func (o *operation) switchBranch(b *branch) {
	if o == nil || o.err != nil {
		return
	}
	_, err := o.git.Exec(
		"switch",
		"--no-guess",
		string(b.branch),
	)
	if err != nil {
		o.err = fmt.Errorf("switchBranch: %w", err)
		return
	}
}

func (o *operation) createSymref(branch, refBranch *branch, reason string) {
	if o == nil || o.err != nil {
		return
	}
	_, err := o.git.Exec(
		"symbolic-ref",
		"-m",
		reason,
		string(branch.baseRefName()),
		string(refBranch.refName()),
	)
	if err != nil {
		o.err = fmt.Errorf("createSymref: %w", err)
		return
	}
}

func (o *operation) createRef(branch *branch, commit git.ObjectName) {
	if o == nil || o.err != nil {
		return
	}
	_, err := o.git.Exec(
		"update-ref",
		"--create-reflog",
		string(branch.startRefName()),
		string(commit),
		strings.Repeat("0", 40),
	)
	if err != nil {
		o.err = fmt.Errorf("createRef: %w", err)
		return
	}
}

func (o *operation) getRef(branch *branch) git.Ref {
	if o == nil || o.err != nil {
		return git.Ref{}
	}
	res, err := o.git.Exec(
		"show-ref",
		"--verify",
		branch.refName().String(),
	)
	if err != nil {
		o.err = fmt.Errorf("getRef: %w", err)
		return git.Ref{}
	}
	ref, err := git.ParseRef(strings.TrimSpace(res.Stdout.String()))
	if err != nil {
		o.err = fmt.Errorf("ParseRef: %w", err)
		return git.Ref{}
	}
	return ref
}

func (o *operation) updateRef(refName git.RefName, newCommit, oldCommit git.ObjectName) {
	if o == nil || o.err != nil {
		return
	}
	_, err := o.git.Exec(
		"update-ref",
		"--create-reflog",
		string(refName),
		string(newCommit),
		string(oldCommit),
	)
	if err != nil {
		o.err = fmt.Errorf("updateRef: %w", err)
		return
	}
}

func (o *operation) deleteRef(refName git.RefName, oldCommit git.ObjectName) {
	if o == nil || o.err != nil {
		return
	}
	_, err := o.git.Exec(
		"update-ref",
		"--no-deref",
		"-d",
		string(refName),
		string(oldCommit),
	)
	if err != nil {
		o.err = fmt.Errorf("deleteRef: %w", err)
		return
	}
}

func (o *operation) rebaseOnto(b *branch) {
	if o == nil || o.err != nil {
		return
	}
	_, err := o.git.Exec(
		"rebase",
		"--committer-date-is-author-date",
		"--onto",
		b.baseRefName().String(),
		b.startRefName().String(),
		b.refName().String(),
	)
	if err != nil {
		o.err = fmt.Errorf("rebaseOnto: %w", err)
		return
	}
}

func (o *operation) pushForce(b *branch) {
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

// func (o *operation) pushUpstream(b *branch) {
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

func (o *operation) listStackerRefs() []git.Ref {
	if o == nil || o.err != nil {
		return nil
	}
	res, err := o.git.Exec(
		"for-each-ref",
		"--format=%(objectname) %(refname)",
		stackerRefPrefix,
	)
	if err != nil {
		o.err = fmt.Errorf("rebaseOnto: %w", err)
		return nil
	}
	var refs []git.Ref
	scan := bufio.NewScanner(&res.Stdout)
	for scan.Scan() {
		ref, err := git.ParseRef(scan.Text())
		if err != nil {
			o.err = fmt.Errorf("ParseRef: %w", err)
			return nil
		}
		refs = append(refs, ref)
	}
	return refs
}

func (o *operation) listRefs() []git.Ref {
	if o == nil || o.err != nil {
		return nil
	}
	res, err := o.git.Exec(
		"show-ref",
	)
	if err != nil {
		o.err = fmt.Errorf("rebaseOnto: %w", err)
		return nil
	}
	var refs []git.Ref
	scan := bufio.NewScanner(&res.Stdout)
	for scan.Scan() {
		ref, err := git.ParseRef(scan.Text())
		if err != nil {
			o.err = fmt.Errorf("ParseRef: %w", err)
			return nil
		}
		refs = append(refs, ref)
	}
	return refs
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
