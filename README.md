# stc

(work in progress; very experimental)

Easy git rebasing of stacked feature branches

## Audience

stc is recommended for a workflow where
*  individual contributors work in the same repository,
*  commits are made to dev branches that are forked from the main or a topic
   branch,
*  each dev branch is completely owned and used by the one contributor and
*  contributors are expected to provide a clean commit history (without merges)
   for review.

stc helps with managing entire *stacks* of these dev branches, i.e. when
they are branched off of and depend on each other, allowing for a more rapid
development. While one branch is undergoing review another branch can be stacked
onto it and worked in.

## Workflow

1. Create a new branch.

   `stc start my-new-feature`

2. Do some work.

   `git add ...`. `git commit ...`. Repeat.

3. Meanwhile the base branch has merges from other contributors. Fetch
   everything.

   `stc sync`

4. Rebase the branch, so it forks off at the head of the base branch.

   `stc rebase`

5. Publish the branch by pushing it the first time.

   `stc push`

6. Oh, there's a bad commit.

   `git rebase -i HEAD~5`

7. Now the local branch and the remote branch divert. Push again to (forcefully)
   set remote branch to what local branch points at.

   `stc push`

## Installation

Requires Rust toolchain to be installed.

```sh
cargo install stc
```

```sh
cargo install --git https://github.com/cloneable/stc
```

## Usage

```
stc 0.1.0

[WIP] Easy stacking of dev branches in git repositories.

USAGE:
    stc <SUBCOMMAND>

OPTIONS:
    -h, --help       Print help information
    -V, --version    Print version information

SUBCOMMANDS:
    clean     Cleans any stc related refs and settings from repo.
    fix       Adds, updates and deletes tracking refs if needed.
    help      Print this message or the help of the given subcommand(s)
    init      Initializes the repo and tries to set stc refs for any non-default branches.
    push      Sets remote branch head to what local branch head points to.
    rebase    Rebases current branch on top of its base branch.
    start     Starts a new branch off of current branch.
    sync      Fetches all branches and tags and prunes deleted ones.
```

## Under the Hood

### Tracking Refs

Stc uses custom refs to track branches:

*  `refs/stc/base/<branchname>`

   A symref to the parent branch in the stack. Created when a branch is created.
   Only updated when a branch is moved within the stack. Deleted when either
   branch is deleted.

*  `refs/stc/start/<branchname>`

   The commit where the branch starts. Created when a branch is created. Updated
   after a rebase. Deleted when the branch is deleted.

*  `refs/stc/remote/<branchname>`

   The commit where the remote branch head is expected to be. Created when a
   branch is pushed for the first time. Updated after a push. Deleted when the
   branch is deleted.

Stc does not update/delete any refs outside `refs/stc/`.

### Config

Stc puts a few settings into `$REPO/.git/config`:

*  `log.excludeDecoration = refs/stc/`

   To hide the trackings refs in `git log` output.

*  `transfer.hideRefs = refs/stc/`

   Probably over-paranoid. Just so there's no possibility these refs leave the
   local repo.

*  `stc.*`

   Any stc specific settings.

Stc will only touch the repo's config. It's safe to make these settings in
`--global` or `--system`.

### Git Commands

Stc mainly uses two commands to manage branches:

*  `git rebase --committer-date-is-author-date --onto refs/stc/base/<branch> refs/stc/start/<branch> <branch>`

   Rebases a branch onto its base branch starting at start marker. If there are
   conflict, rebasing will stop. Fix conflicts and use `--continue` or
   `--abort`. After a successful rebase the start marker ref will be moved.

*  `git push --set-upstream --force-with-lease=<branch>:<expected-commit> <remote> <branch>:<branch>`

   Pushes a branch to remote, potentially replacing the commit chain, so the
   remote branch looks exactly like the local branch. `<expected-commit>` is
   what `refs/stc/remote/<branch>` points to.

In addition, stc uses a few more commands to help keeping track of things,
like `for-each-ref`, `update-ref`, `symbolic-ref`, `check-ref-format`.
