use crate::git;
use ::std::{
    option::Option::{self, Some},
    result::Result::{self, Err, Ok},
    string::String,
    vec::Vec,
};

pub struct Stc<G: git::Git> {
    git: G,
}

impl<G: git::Git> Stc<G> {
    pub fn new(git: G) -> Self {
        Stc { git }
    }

    pub fn init(&self) -> Result<(), git::Status> {
        let g = &self.git;
        g.config_add("transfer.hideRefs", git::STC_REF_PREFIX)?;
        g.config_add("log.excludeDecoration", git::STC_REF_PREFIX)?;

        // TODO: read refs, branches, remotes
        // TODO: validate stc refs against branches
        // TODO: determine list of needed refs
        // TODO: print and create list of created refs

        Ok(())
    }

    pub fn clean(&self) -> Result<(), git::Status> {
        let g = &self.git;
        g.config_unset_pattern("transfer.hideRefs", git::STC_REF_PREFIX)?;
        g.config_unset_pattern("log.excludeDecoration", git::STC_REF_PREFIX)?;

        // TODO: for each branch
        // TODO: ... check if fully merged
        // TODO: ... check if remote ref == local branch
        // TODO: ... delete stc refs
        // TODO: ... or print warning

        Ok(())
    }

    pub fn start(&self, name: String) -> Result<(), git::Status> {
        let g = &self.git;
        let repo = g.snapshot()?;
        let base_branch = repo.head().ok_or_else(|| git::Status::with(1))?;
        let new_name = g.check_branchname(&name)?;
        g.create_branch(&new_name, base_branch)?;
        g.switch_branch(&new_name)?;
        g.create_symref(
            &new_name.stc_base_refname(),
            &base_branch.refname(),
            "stc: base branch marker",
        )?;
        let base_refname = base_branch.refname();
        let base_ref = repo
            .get_ref(&base_refname)
            .ok_or_else(|| git::Status::with(1))?;
        g.create_ref(&new_name.stc_start_refname(), &base_ref.objectname)?;

        Ok(())
    }

    pub fn push(&self) -> Result<(), git::Status> {
        let g = &self.git;

        let expected_commit: git::ObjectName;
        {
            let repo = g.snapshot()?;
            let cur_branch = repo.head().ok_or_else(|| git::Status::with(1))?;
            let stc_base_refname = cur_branch.stc_base_refname();
            let base_symref = repo
                .get_ref(&stc_base_refname)
                .ok_or_else(|| git::Status::with(1))?;
            let base_ref = repo
                .get_ref(&base_symref.symref_target)
                .ok_or_else(|| git::Status::with(1))?;
            if let Some(remote_ref) = repo.get_ref(&cur_branch.stc_remote_refname()) {
                expected_commit = remote_ref.objectname.owning_clone();
            } else {
                expected_commit = git::NON_EXISTANT_OBJECT;
            }
            g.push(cur_branch, &base_ref.remote, &expected_commit)?;
        }
        {
            let repo = g.snapshot()?;
            let cur_branch = repo.head().ok_or_else(|| git::Status::with(1))?;
            let cur_refname = cur_branch.refname();
            let cur_ref = repo
                .get_ref(&cur_refname)
                .ok_or_else(|| git::Status::with(1))?;
            g.update_ref(
                &cur_branch.stc_remote_refname(),
                &cur_ref.objectname,
                &expected_commit,
            )?;
        }

        Ok(())
    }

    pub fn rebase(&self) -> Result<(), git::Status> {
        let g = &self.git;

        let repo = g.snapshot()?;
        let branch = repo.head().ok_or_else(|| git::Status::with(1))?;
        let stc_base_refname = branch.stc_base_refname();
        let stc_start_refname = branch.stc_start_refname();
        let base_ref = repo
            .get_ref(&stc_base_refname)
            .ok_or_else(|| git::Status::with(1))?;
        let start_ref = repo
            .get_ref(&stc_start_refname)
            .ok_or_else(|| git::Status::with(1))?;
        g.rebase_onto(branch)?;
        g.update_ref(
            &branch.stc_start_refname(),
            &base_ref.objectname,
            &start_ref.objectname,
        )?;

        Ok(())
    }

    pub fn sync(&self) -> Result<(), git::Status> {
        let g = &self.git;

        g.fetch_all_prune()?;

        Ok(())
    }

    pub fn fix(&self, branch: Option<String>, base: Option<String>) -> Result<(), git::Status> {
        let g = &self.git;

        let repo = g.snapshot()?;
        // TODO: this is hacky. refactor.
        if let Some(branchname) = branch {
            if let Some(base_branchname) = base {
                let branch = g.check_branchname(&branchname)?;
                let base_branch = g.check_branchname(&base_branchname)?;
                if let Some(base_symref) = repo.get_ref(&branch.stc_base_refname()) {
                    if base_symref.symref_target != base_branch.refname() {
                        return Err(git::Status::new(
                            1,
                            Vec::<u8>::new(),
                            "base branch already defined".as_bytes().to_vec(),
                        ));
                    }
                } else {
                    g.create_symref(
                        &branch.stc_base_refname(),
                        &base_branch.refname(),
                        "stc: set base branch",
                    )?;
                }
                if let Some(_start_ref) = repo.get_ref(&branch.stc_start_refname()) {
                    // TODO: check if base or ancestor of base
                } else {
                    let forkpoint = g.forkpoint(&base_branch.refname(), &branch.refname())?;
                    g.create_ref(&branch.stc_start_refname(), &forkpoint)?;
                }
            } else {
                return Err(git::Status::new(
                    1,
                    Vec::<u8>::new(),
                    "base not specified".as_bytes().to_vec(),
                ));
            }
        }

        let repo = g.snapshot()?;
        for branch in repo.tracked_branches() {
            if repo.get_ref(&branch.refname()).is_none() {
                if let Some(r) = repo.get_ref(&branch.stc_base_refname()) {
                    g.delete_symref(&r.name)?;
                }
                if let Some(r) = repo.get_ref(&branch.stc_start_refname()) {
                    g.delete_ref(&r.name, &r.objectname)?;
                }
                if let Some(r) = repo.get_ref(&branch.stc_remote_refname()) {
                    g.delete_ref(&r.name, &r.objectname)?;
                }
            }
        }

        let repo = g.snapshot()?;
        for branch in repo.tracked_branches() {
            // for each existing branch that's somehow still being tracked:
            let base_symref_name = branch.stc_base_refname();
            let start_refname = branch.stc_start_refname();
            if let Some(base_symref) = repo.get_ref(&base_symref_name) {
                // there's a base symref
                if repo.get_ref(&start_refname).is_none() {
                    // but no start ref
                    if repo.get_ref(&base_symref.symref_target).is_none() {
                        // TODO: base branch doesn't exist (anymore)
                        continue;
                    }
                    // figure out forkpoint from what the base symref points to and the branch
                    // TODO: forkpoint can fail
                    let forkpoint = g.forkpoint(&base_symref.symref_target, &branch.refname())?;
                    // write the commit as start ref
                    g.create_ref(&branch.stc_start_refname(), &forkpoint)?;
                }
            } else {
                // there's no base symref
                if let Some(_start_ref) = repo.get_ref(&start_refname) {
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

        Ok(())
    }
}
