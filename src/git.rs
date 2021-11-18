use ::const_format::concatcp;
use ::std::format;
use ::std::string::String;
use ::std::string::ToString;
use ::std::vec::Vec;

pub const NON_EXISTANT_OBJECT: &'static str = "0000000000000000000000000000000000000000";

pub const STACKER_REF_PREFIX: &'static str = "refs/stacker/";
pub const STACKER_BASE_REF_PREFIX: &'static str = concatcp!(STACKER_REF_PREFIX, "base/");
pub const STACKER_START_REF_PREFIX: &'static str = concatcp!(STACKER_REF_PREFIX, "start/");
pub const STACKER_REMOTE_REF_PREFIX: &'static str = concatcp!(STACKER_REF_PREFIX, "remote/");

pub const BRANCH_REF_PREFIX: &'static str = "refs/heads/";
pub const TAG_REF_PREFIX: &'static str = "refs/tags/";

pub trait Git {
    fn exec(&self, args: &[&str]) -> Result;
    fn fail(&self) -> !;

    fn check_branch_name(&self, name: &String) -> bool {
        let res = self.exec(&["check-ref-format", "--branch", name]);
        res.ok()
    }

    fn create_branch(&self, name: &BranchName, base: &BranchName) {
        let res = self.exec(&["branch", "--create-reflog", name.as_str(), base.as_str()]);
        if !res.ok() {
            self.fail();
        }
    }

    fn switch_branch(&self, b: &BranchName) {
        let res = self.exec(&["switch", "--no-guess", b.as_str()]);
        if !res.ok() {
            self.fail();
        }
    }

    fn create_symref(&self, name: &RefName, target: &RefName, reason: &String) {
        let res = self.exec(&["symbolic-ref", "-m", reason, name.as_str(), target.as_str()]);
        if !res.ok() {
            self.fail();
        }
    }

    fn delete_symref(&self, name: &RefName) {
        let res = self.exec(&["symbolic-ref", "--delete", name.as_str()]);
        if !res.ok() {
            self.fail();
        }
    }

    fn create_ref(&self, name: &RefName, commit: &ObjectName) {
        let res = self.exec(&[
            "update-ref",
            "--no-deref",
            "--create-reflog",
            name.as_str(),
            commit.as_str(),
            NON_EXISTANT_OBJECT,
        ]);
        if !res.ok() {
            self.fail();
        }
    }

    fn update_ref(&self, name: &RefName, new_commit: &ObjectName, cur_commit: &ObjectName) {
        let res = self.exec(&[
            "update-ref",
            "--no-deref",
            "--create-reflog",
            name.as_str(),
            new_commit.as_str(),
            cur_commit.as_str(),
        ]);
        if !res.ok() {
            self.fail();
        }
    }

    fn delete_ref(&self, name: &RefName, cur_commit: &ObjectName) {
        let res = self.exec(&[
            "update-ref",
            "--no-deref",
            "-d",
            name.as_str(),
            cur_commit.as_str(),
        ]);
        if !res.ok() {
            self.fail();
        }
    }

    fn rebase_onto(&self, name: &BranchName) {
        let res = self.exec(&[
            "rebase",
            "--committer-date-is-author-date",
            "--onto",
            name.stacker_base_refname().as_str(),
            name.stacker_start_refname().as_str(),
            name.as_str(),
        ]);
        if !res.ok() {
            self.fail();
        }
    }

    fn push(&self, name: &BranchName, remote: &RemoteName, expect: &ObjectName) {
        let res = self.exec(&[
            "push",
            "--set-upstream",
            format!("--force-with-lease={}:{}", name.as_str(), expect.as_str()).as_str(),
            remote.as_str(),
            format!("{}:{}", name.as_str(), name.as_str()).as_str(),
        ]);
        if !res.ok() {
            self.fail();
        }
    }

    fn config_set(&self, key: &String, value: &String) {
        let res = self.exec(&["config", "--local", key, value]);
        if !res.ok() {
            self.fail();
        }
    }

    fn config_add(&self, key: &String, value: &String) {
        let res = self.exec(&["config", "--local", "--add", key, value]);
        if !res.ok() {
            self.fail();
        }
    }

    fn config_unset_pattern(&self, key: &String, pattern: &String) {
        let res = self.exec(&[
            "config",
            "--local",
            "--fixed-value",
            "--unset-all",
            key,
            pattern,
        ]);
        // 5 means the nothing matched.
        if !res.ok() && res.exitcode != 5 {
            self.fail();
        }
    }

    fn fetch_all_prune(&self) {
        let res = self.exec(&["fetch", "--all", "--prune"]);
        if !res.ok() {
            self.fail();
        }
    }

    fn forkpoint(&self, base: &RefName, branch: &RefName) -> ObjectName {
        let res = self.exec(&["merge-base", "--fork-point", base.as_str(), branch.as_str()]);
        if !res.ok() {
            self.fail();
        }
        // TODO: handle not found
        return ObjectName(String::from_utf8_lossy(&res.stdout).to_string());
    }
}

pub struct Result {
    pub exitcode: i32,
    pub stdout: Vec<u8>,
    pub stderr: Vec<u8>,
}

impl Result {
    pub fn new(exitcode: i32, stdout: Vec<u8>, stderr: Vec<u8>) -> Self {
        Result {
            exitcode,
            stdout,
            stderr,
        }
    }

    pub fn ok(&self) -> bool {
        self.exitcode == 0
    }
}

pub struct BranchName(String);

impl BranchName {
    pub fn as_str(&self) -> &str {
        self.0.as_str()
    }

    pub fn stacker_base_refname(&self) -> RefName {
        RefName(STACKER_BASE_REF_PREFIX.to_string() + &self.0)
    }

    pub fn stacker_start_refname(&self) -> RefName {
        RefName(STACKER_START_REF_PREFIX.to_string() + &self.0)
    }

    pub fn stacker_remote_refname(&self) -> RefName {
        RefName(STACKER_REMOTE_REF_PREFIX.to_string() + &self.0)
    }
}

pub struct RefName(String);

impl RefName {
    pub fn as_str(&self) -> &str {
        self.0.as_str()
    }
}

pub struct ObjectName(String);

impl ObjectName {
    pub fn as_str(&self) -> &str {
        self.0.as_str()
    }
}
pub struct RemoteName(String);

impl RemoteName {
    pub fn as_str(&self) -> &str {
        self.0.as_str()
    }
}
