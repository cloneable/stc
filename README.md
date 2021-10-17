# stacker

Easy git rebasing of stacked feature branches

## Use

Stacker is good at rebasing multiple local branches that "sit" on each other:
stacked branches. These are usually branches owned by one person and send as
individual pull requests. If the base branch changes or any of the branches in
between the ones on top need to be rebased.

* `stacker init` checks and creates any stacker-related refs.
* `stacker show` prints a graph of the entire stack.
* `stacker clean [--force]` removes any stacker-related refs. `--force` must be
  used to when cleaning is called mid-rebase.
* `stacker purge` deletes all deletable (=fully merged) branches.
* `stacker create <branch>` creates new branch, marks it for remote tracking and
  switches to it.
* `stacker rebase [<branch>]` rebases current stack or stack starting at
  `<branch>`.
* `stacker push [<branch>]` pushes all outdated branches or all branches
  starting at `<branch>`.
* `stacker delete <branch>` safely deletes local branch and remote branch.

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

2.`git update-ref --create-reflog refs/stacker/start/<branch> <commit> 0000000000000000000000000000000000000000`

   Mark the base of a branch with a new ref.

3. `git update-ref --create-reflog refs/stacker/start/<branch> <new-commit> <old-commit>`

   Move the base marker after rebasing.

4. `git update-ref -d refs/stacker/start/<branch> <commit>`

   Delete the base marker.

5. `git show-ref --verify <ref>`

   Get commit of refs.

6. `git rebase --committer-date-is-author-date --onto <base-branch> <base-marker> <branch>`

   Rebasing a branch onto base branch starting at base marker. If there are
   conflict, rebasing will stop. Fix conflicts and use `--continue` or
   `--abort`.

7. `git push --force-with-lease=<branch>:<expected-commit> <local-branch>:<remote-branch>`

   Push a rebased branch to remote, replacing the commit chain.

TODO: use `git check-ref-format` to check names?