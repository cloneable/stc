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

func (s *Stacker) Init(ctx context.Context) error {
	op := op(s.git)
	op.configAdd("transfer.hideRefs", git.StackerRefPrefix)
	op.configAdd("log.excludeDecoration", git.StackerRefPrefix)

	// TODO: read refs, branches, remotes
	// TODO: validate stacker refs against branches
	// TODO: determine list of needed refs
	// TODO: print and create list of created refs
	return op.Err()
}

func (s *Stacker) Clean(ctx context.Context) error {
	op := op(s.git)
	op.configUnsetPattern("transfer.hideRefs", git.StackerRefPrefix)
	op.configUnsetPattern("log.excludeDecoration", git.StackerRefPrefix)
	// for _, ref := range op.listStackerRefs() {
	// 	op.deleteRef(ref.Name(), ref.ObjectName())
	// }
	// TODO: for each branch
	// TODO: ... check if fully merged
	// TODO: ... check if remote ref == local branch
	// TODO: ... delete stacker refs
	// TODO: ... or print warning
	return op.Err()
}

func (s *Stacker) Start(ctx context.Context, name string) error {
	op := op(s.git)
	op.snapshot()
	baseB := op.head()
	newName := op.parseBranchName(name)
	op.createBranch(newName, baseB)
	op.switchBranch(newName)
	op.createSymref(newName.StackerBaseRefName(), baseB.RefName(), "stacker: base branch marker")
	baseRef := op.ref(baseB.RefName())
	op.createRef(newName.StackerStartRefName(), baseRef.ObjectName())
	return op.Err()
}

func (s *Stacker) Push(ctx context.Context) error {
	op := op(s.git)

	var expectedCommit git.ObjectName
	{
		op.snapshot()
		curB := op.head()
		symRef := op.ref(curB.StackerBaseRefName())
		baseRef := op.ref(symRef.SymRefTarget())
		if op.hasRef(curB.StackerRemoteRefName()) {
			expectedCommit = op.ref(curB.StackerRemoteRefName()).ObjectName()
		} else {
			expectedCommit = git.NonExistantObject
		}
		op.push(curB, baseRef.Remote(), expectedCommit)
	}
	{
		op.snapshot()
		curB := op.head()
		curRef := op.ref(curB.RefName())
		op.updateRef(curB.StackerRemoteRefName(), curRef.ObjectName(), expectedCommit)
	}

	// TODO: for each branch
	// TODO: ... determine state (already pushed?)
	// TODO: ... determine upstream
	// TODO: ... push branch to remote
	return op.Err()
}

func (s *Stacker) Rebase(ctx context.Context) error {
	op := op(s.git)

	// TODO: branches

	op.snapshot()
	branch := op.head()
	baseRef := op.ref(branch.StackerBaseRefName())
	startRef := op.ref(branch.StackerStartRefName())
	op.rebaseOnto(branch)
	op.updateRef(branch.StackerStartRefName(), baseRef.ObjectName(), startRef.ObjectName())

	// TODO: if len(branch) == 0 use current head as branch (head must be branch head)
	// TODO: for each branch
	// TODO: ... determine list of all stacked branches
	// TODO: ... add unselected branches to list
	// TODO: topologically sort selected branches
	// TODO: for each selected branch
	// TODO: ... call git rebase --onto
	// TODO: ... update stacker ref
	return op.Err()
}

func (s *Stacker) Sync(ctx context.Context) error {
	op := op(s.git)
	op.fetchAllPrune()

	// TODO: fast-forward base branches that are not stacker branches.
	// TODO: push all or selected ahead/rebased stacker branches.

	return op.Err()
}

func (s *Stacker) Fix(ctx context.Context, branches ...string) error {
	op := op(s.git)

	// TODO: this is hacky. refactor.
	if len(branches) == 2 {
		branch := op.parseBranchName(branches[0])
		baseBranch := op.parseBranchName(branches[1])
		if baseSymrefName := branch.StackerBaseRefName(); op.hasRef(baseSymrefName) {
			baseSymref := op.ref(baseSymrefName)
			if baseSymref.SymRefTarget() != baseBranch.RefName() {
				return op.Failf("base branch already defined: %v", baseSymref.SymRefTarget())
			}
		} else {
			op.createSymref(branch.StackerBaseRefName(), baseBranch.RefName(), "stacker: set base branch")
		}
		if startRefName := branch.StackerStartRefName(); op.hasRef(startRefName) {
			// TODO: check if base or ancestor of base
		} else {
			forkpoint := op.forkpoint(baseBranch.RefName(), branch.RefName())
			op.createRef(branch.StackerStartRefName(), forkpoint)
		}
		return nil
	} else if len(branches) != 0 {
		return op.Failf("invalid arguments: %v", branches)
	}

	op.snapshot()
	for _, branch := range op.trackedBranches() {
		if !op.hasRef(branch.RefName()) {
			if r := branch.StackerBaseRefName(); op.hasRef(r) {
				op.deleteSymref(r)
			}
			if r := branch.StackerStartRefName(); op.hasRef(r) {
				ref := op.ref(r)
				op.deleteRef(r, ref.ObjectName())
			}
			if r := branch.StackerRemoteRefName(); op.hasRef(r) {
				ref := op.ref(r)
				op.deleteRef(r, ref.ObjectName())
			}
		}
	}

	op.snapshot()
	for _, branch := range op.trackedBranches() {
		// for each existing branch that's somehow still being tracked:
		baseSymrefName := branch.StackerBaseRefName()
		startRefName := branch.StackerStartRefName()
		if op.hasRef(baseSymrefName) {
			// there's a base symref
			if !op.hasRef(startRefName) {
				// but no start ref
				baseSymref := op.ref(branch.StackerBaseRefName())
				if !op.hasRef(baseSymref.SymRefTarget()) {
					// TODO: base branch doesn't exist (anymore)
					continue
				}
				// figure out forkpoint from what the base symref points to and the branch
				// TODO: forkpoint can fail
				forkpoint := op.forkpoint(baseSymref.SymRefTarget(), branch.RefName())
				// write the commit as start ref
				op.createRef(branch.StackerStartRefName(), forkpoint)
			}
		} else {
			// there's no base symref
			if op.hasRef(startRefName) {
				// but there's a start ref
				// TODO: check for branch at that commit? consult reflog?
			}
		}
	}

	// TODO: no /base/, but /start/ -> look for branch head at /start/, set /base/
	// TODO: no /start/, but /base/ -> use git merge-base to find fork point
	// TODO: no /start/ nor /base/ -> do nothing, offer explicit way to track
	// TODO: no /remote/, but remote branch exists? -> set ref, if ancestor, if not -> error
	// TODO: no remote branch, but /remote/ -> delete ref (check origin?)

	return op.Err()
}
