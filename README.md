# Stacker

Easy git rebasing of stacked feature branches

## Workflow

Stacker is recommended for a workflow where individual contributors work in the
same repository where commits are made in separate dev branches each owned by
one contributor and merged via pull request with code review into the main
branch or a topic branch. Stacker helps with managing entire *stacks* of these
dev branches, i.e. when they are branched off of and depend on each other,
allowing for a more rapid development.

## Installation

Requires Go toolchain 1.17 or later.

```sh
go install github.com/cloneable/stacker@latest
```

Make sure `stacker` is in your $PATH.

```sh
export PATH=$(go env GOPATH)/bin:$PATH
```

## Usage

(Planned Features)

* `stacker init [--force]` checks and creates any stacker-related refs. If
  `--force` is used invalid refs are removed or replaced too.
* `stacker clean [--force] [<branch>...]` removes any stacker-related refs.
  `--force` must be used to when cleaning is called mid-rebase.
* `stacker show` lists all stacker tracked branches with status as graph.
* `stacker start <branch>` starts new branch, marks it for remote tracking and
  switches to it.
* `stacker publish <branch>...` pushes one or more branches to remote and
  marks them for tracking.
* `stacker delete <branch>` safely deletes local branch and remote branch.
* `stacker rebase [<branch>...]` rebases current stack or stack starting at
  `<branch>`.
* `stacker sync [<branch>...]` fetches remote branches of stacker tracked
  branches and base branches, prunes deleted remote refs and pushes all outdated
  branches.

## Under the Hood

Stacker uses custom refs to track branches:

*  `refs/stacker/base/<branchname>`

   A symref to the parent branch in the stack. Created when a branch is created.
   Only updated when a branch is moved within the stack. Deleted when either
   branch is deleted.

*  `refs/stacker/start/<branchname>`

   The commit where the branch starts. Created when a branch is created. Updated
   after a rebase. Deleted when the branch is deleted.

Stacker does not update/delete any refs outside `refs/stacker/`.

## Used Git Commands

Following is the list of all used git commands. `git` is called with all
environment variables inherited from `stacker`. <br> Protected branch names:
master, main, release, production, staging. These names are not used with any
commands, with the exception of rebasing `--onto`.

1. `git for-each-ref --format='%(HEAD)%(refname) %(objecttype) %(objectname) %(upstream:remotename) %(upstream)'`

   Gather all refs at startup.

2. `git update-ref --create-reflog refs/stacker/start/<branch> <commit> 0000000000000000000000000000000000000000`

   Mark the base of a branch with a new ref.

3. `git update-ref --create-reflog refs/stacker/start/<branch> <new-commit> <old-commit>`

   Move the base marker after rebasing.

4. `git update-ref -d refs/stacker/start/<branch> <commit>`

   Delete the base marker.

5. `git show-ref --verify <ref>`

   Get commit of refs.

6. `git rebase --committer-date-is-author-date --onto refs/stacker/base/<branch> refs/stacker/start/<branch> <branch>`

   Rebasing a branch onto base branch starting at base marker. If there are
   conflict, rebasing will stop. Fix conflicts and use `--continue` or
   `--abort`. After a successful rebase the start marker ref will be moved.

7. `git push --force-with-lease=<branch>:<expected-commit> <local-branch>:<remote-branch>`

   Push a rebased branch to remote, replacing the commit chain.
