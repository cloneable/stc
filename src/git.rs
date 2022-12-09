use color_eyre::Result;
use const_format::concatcp;
use csv::ReaderBuilder;
use serde::Deserialize;
use std::{
    borrow::{Cow, ToOwned},
    collections::{BTreeSet, HashMap},
    convert::{AsRef, Into},
    string::{String, ToString},
    write,
};
use thiserror::Error;

pub const NON_EXISTANT_OBJECT: ObjectName<'static> =
    ObjectName::new("0000000000000000000000000000000000000000");

pub const STC_REF_PREFIX: &str = "refs/stc/";
pub const STC_BASE_REF_PREFIX: &str = concatcp!(STC_REF_PREFIX, "base/");
pub const STC_START_REF_PREFIX: &str = concatcp!(STC_REF_PREFIX, "start/");
pub const STC_REMOTE_REF_PREFIX: &str = concatcp!(STC_REF_PREFIX, "remote/");

pub const BRANCH_REF_PREFIX: &str = "refs/heads/";

#[derive(Error, Debug)]
#[error("exitcode={exitcode:?}, stdout={stdout:?}, stderr={stderr:?}")]
pub struct ExecStatus {
    pub exitcode: i32,
    pub stdout: String,
    pub stderr: String,
}

impl ExecStatus {
    pub fn new(exitcode: i32, stdout: Vec<u8>, stderr: Vec<u8>) -> Self {
        // TODO: use OsString.
        ExecStatus {
            exitcode,
            stdout: String::from_utf8_lossy(stdout.as_slice()).to_string(),
            stderr: String::from_utf8_lossy(stderr.as_slice()).to_string(),
        }
    }
}

pub trait Git {
    fn exec(&self, args: &[&str]) -> Result<ExecStatus>;

    fn exec_log(&self, args: &[&str]) -> Result<()> {
        match self.exec(args) {
            Ok(_) => {
                ::std::eprintln!("[OK] git {:?}", args);
                Ok(())
            }
            Err(err) => match err.downcast::<ExecStatus>() {
                Ok(status) => {
                    ::std::assert_ne!(status.exitcode, 0);
                    ::std::eprintln!("[ERR {:?}] git {:?}", status.exitcode, args);
                    Err(status.into())
                }
                Err(err) => Err(err),
            },
        }
    }

    fn snapshot(&self) -> Result<Repository> {
        let status = self.exec(&["for-each-ref", "--format", FIELD_FORMATS.join(",").as_str()])?;
        let refs = parse_ref(status.stdout.as_bytes())?;
        let head = refs
            .values()
            .find(|r| r.head)
            .and_then(|r| r.name.branchname())
            .map(|bn| bn.owning_clone());
        Ok(Repository { refs, head })
    }

    fn check_branchname<'a>(&self, name: &'a str) -> Result<BranchName<'a>> {
        self.exec(&["check-ref-format", "--branch", name])?;
        Ok(BranchName(Cow::Owned(name.to_string())))
    }

    fn create_branch(&self, name: &BranchName, base: &BranchName) -> Result<()> {
        self.exec_log(&["branch", "--create-reflog", name.as_str(), base.as_str()])
            .map(|_| {})
    }

    fn switch_branch(&self, b: &BranchName) -> Result<()> {
        self.exec_log(&["switch", "--no-guess", b.as_str()])
            .map(|_| {})
    }

    fn create_symref(&self, name: &RefName, target: &RefName, reason: &'static str) -> Result<()> {
        self.exec_log(&["symbolic-ref", "-m", reason, name.as_str(), target.as_str()])
            .map(|_| {})
    }

    fn delete_symref(&self, name: &RefName) -> Result<()> {
        self.exec_log(&["symbolic-ref", "--delete", name.as_str()])
            .map(|_| {})
    }

    fn create_ref(&self, name: &RefName, commit: &ObjectName) -> Result<()> {
        self.exec_log(&[
            "update-ref",
            "--no-deref",
            "--create-reflog",
            name.as_str(),
            commit.as_str(),
            NON_EXISTANT_OBJECT.as_str(),
        ])
        .map(|_| {})
    }

    fn update_ref(
        &self,
        name: &RefName,
        new_commit: &ObjectName,
        cur_commit: &ObjectName,
    ) -> Result<()> {
        self.exec_log(&[
            "update-ref",
            "--no-deref",
            "--create-reflog",
            name.as_str(),
            new_commit.as_str(),
            cur_commit.as_str(),
        ])
        .map(|_| {})
    }

    fn delete_ref(&self, name: &RefName, cur_commit: &ObjectName) -> Result<()> {
        self.exec_log(&[
            "update-ref",
            "--no-deref",
            "-d",
            name.as_str(),
            cur_commit.as_str(),
        ])
        .map(|_| {})
    }

    fn rebase_onto(&self, name: &BranchName) -> Result<()> {
        self.exec_log(&[
            "rebase",
            "--committer-date-is-author-date",
            "--onto",
            name.stc_base_refname().as_str(),
            name.stc_start_refname().as_str(),
            name.as_str(),
        ])
        .map(|_| {})
    }

    fn push(&self, name: &BranchName, remote: &RemoteName, expect: &ObjectName) -> Result<()> {
        self.exec_log(&[
            "push",
            "--set-upstream",
            ::std::format!("--force-with-lease={}:{}", name.as_str(), expect.as_str()).as_str(),
            remote.as_str(),
            ::std::format!("{}:{}", name.as_str(), name.as_str()).as_str(),
        ])
        .map(|_| {})
    }

    fn config_set(&self, key: &str, value: &str) -> Result<()> {
        self.exec_log(&["config", "--local", key, value])
            .map(|_| {})
    }

    fn config_add(&self, key: &str, value: &str) -> Result<()> {
        self.exec_log(&["config", "--local", "--add", key, value])
            .map(|_| {})
    }

    fn config_unset_pattern(&self, key: &str, pattern: &str) -> Result<()> {
        match self.exec_log(&[
            "config",
            "--local",
            "--fixed-value",
            "--unset-all",
            key,
            pattern,
        ]) {
            // 5 means the nothing matched.
            Err(err) => match err.downcast_ref::<ExecStatus>() {
                Some(status) if status.exitcode != 5 => Err(err),
                _ => Ok(()),
            },
            _ => Ok(()),
        }
    }

    fn fetch_all_prune(&self) -> Result<()> {
        self.exec_log(&["fetch", "--all", "--prune"]).map(|_| {})
    }

    fn forkpoint(&self, base: &RefName, branch: &RefName) -> Result<ObjectName> {
        self.exec(&["merge-base", "--fork-point", base.as_str(), branch.as_str()])
            .map(move |status| {
                // TODO: handle not found
                ObjectName(Cow::Owned(status.stdout))
            })
    }
}

pub struct Repository<'a> {
    refs: HashMap<RefName<'a>, Ref<'a>>,
    head: Option<BranchName<'a>>,
}

impl<'a> Repository<'a> {
    pub fn get_ref(&self, name: &'a RefName) -> Option<&'a Ref> {
        self.refs.get(name)
    }

    pub fn head(&self) -> Option<&'a BranchName> {
        self.head.as_ref()
    }

    pub fn tracked_branches(&self) -> Vec<BranchName> {
        self.refs
            .iter()
            .filter(|(name, _)| name.0.starts_with(STC_REF_PREFIX))
            .map(|(name, _)| name.branchname().unwrap())
            .collect::<BTreeSet<_>>()
            .into_iter()
            .collect()
    }
}

#[derive(Eq, PartialEq, Ord, PartialOrd, Debug)]
pub struct BranchName<'a>(Cow<'a, str>);

impl<'a> BranchName<'a> {
    pub const fn new(name: &'a str) -> Self {
        BranchName(Cow::Borrowed(name))
    }

    pub fn owning_clone<'b: 'a>(&'a self) -> BranchName<'b> {
        BranchName(Cow::Owned(self.0.as_ref().to_owned()))
    }

    pub fn as_str(&self) -> &str {
        &self.0
    }

    pub fn refname(&self) -> RefName {
        RefName(Cow::Owned(BRANCH_REF_PREFIX.to_string() + &self.0))
    }

    pub fn stc_base_refname(&self) -> RefName {
        RefName(Cow::Owned(STC_BASE_REF_PREFIX.to_string() + &self.0))
    }

    pub fn stc_start_refname(&self) -> RefName {
        RefName(Cow::Owned(STC_START_REF_PREFIX.to_string() + &self.0))
    }

    pub fn stc_remote_refname(&self) -> RefName {
        RefName(Cow::Owned(STC_REMOTE_REF_PREFIX.to_string() + &self.0))
    }
}

#[derive(Deserialize, Clone, Hash, Eq, PartialEq, Ord, PartialOrd, Debug)]
pub struct RefName<'a>(Cow<'a, str>);

impl<'a> RefName<'a> {
    pub const fn new(name: &'a str) -> Self {
        RefName(Cow::Borrowed(name))
    }

    pub fn owning_clone<'b: 'a>(&'a self) -> RefName<'b> {
        RefName(Cow::Owned(self.0.as_ref().to_owned()))
    }

    pub fn as_str(&self) -> &str {
        &self.0
    }

    pub fn branchname(&'a self) -> Option<BranchName<'a>> {
        if let Some(suffix) = self.0.strip_prefix(STC_REF_PREFIX) {
            if let Some((_, branchname)) = suffix.split_once('/') {
                return Some(BranchName::new(branchname));
            }
        }
        if let Some(branchname) = self.0.strip_prefix(BRANCH_REF_PREFIX) {
            return Some(BranchName::new(branchname));
        }
        None
    }
}

#[derive(Deserialize, PartialEq, Eq, PartialOrd, Clone, Debug)]
pub struct ObjectName<'a>(pub Cow<'a, str>);

impl<'a> ObjectName<'a> {
    pub const fn new(name: &'a str) -> Self {
        ObjectName(Cow::Borrowed(name))
    }

    pub fn owning_clone<'b: 'a>(&'a self) -> ObjectName<'b> {
        ObjectName(Cow::Owned(self.0.as_ref().to_owned()))
    }

    pub fn as_str(&self) -> &str {
        &self.0
    }
}

#[derive(Deserialize, PartialEq, Eq, PartialOrd, Debug)]
pub struct RemoteName<'a>(Cow<'a, str>);

impl<'a> RemoteName<'a> {
    pub const fn new(name: &'a str) -> Self {
        RemoteName(Cow::Borrowed(name))
    }

    pub fn owning_clone<'b: 'a>(&'a self) -> RemoteName<'b> {
        RemoteName(Cow::Owned(self.0.as_ref().to_owned()))
    }

    pub fn as_str(&self) -> &str {
        &self.0
    }
}

#[derive(Deserialize, PartialEq, Eq, Debug)]
#[serde(rename_all = "lowercase")]
pub enum RefType {
    Commit,
    Tree,
    Blob,
    Tag,
}

#[derive(Deserialize, PartialEq, Eq, Debug)]
pub struct Ref<'a> {
    pub name: RefName<'a>,
    pub head: bool,
    pub objectname: ObjectName<'a>,
    pub objecttype: RefType,
    pub track: String,
    pub remote: RemoteName<'a>,
    pub remote_refname: RefName<'a>,
    pub symref_target: RefName<'a>,
    pub upstream_refname: RefName<'a>,
}

const FIELD_FORMATS: [&str; 9] = [
    "%(refname)",                                // name
    "%(if)%(HEAD)%(then)true%(else)false%(end)", // head
    "%(objectname)",                             // objectname
    "%(objecttype)",                             // objecttype
    "%(upstream:trackshort)",                    // track
    "%(upstream:remotename)",                    // remote
    "%(upstream:remoteref)",                     // remote_refname
    "%(symref)",                                 // symref_target
    "%(upstream)",                               // upstream_refname
];

fn parse_ref<'a, R: ::std::io::Read + ::std::fmt::Debug>(
    csv: R,
) -> Result<HashMap<RefName<'a>, Ref<'a>>, ::csv::Error> {
    let mut reader = ReaderBuilder::new()
        .has_headers(false)
        .delimiter(b',')
        .from_reader(csv);
    let ref_vec = reader
        .deserialize::<Ref>()
        .collect::<Result<Vec<Ref>, ::csv::Error>>()?;
    let refs = ref_vec
        .into_iter()
        .map(move |r| (r.name.clone(), r))
        .collect::<HashMap<RefName, Ref>>();
    Ok(refs)
}

#[cfg(test)]
mod tests {
    use super::*;
    use ::std::assert_eq;
    use ::std::convert::From;

    #[test]
    fn test_parse_ref() {
        let csv = "\
refs/heads/moo1,true,123abc,commit,<>,origin,refs/heads/moo,,refs/remotes/origin/moo
";
        let refs = parse_ref(csv.as_bytes()).expect("cannot parse");
        assert_eq!(refs.len(), 1);
        assert_eq!(
            refs.get(&RefName::new("refs/heads/moo1")).unwrap(),
            &Ref {
                name: RefName::new("refs/heads/moo1"),
                head: true,
                objectname: ObjectName::new("123abc"),
                objecttype: RefType::Commit,
                track: String::from("<>"),
                remote: RemoteName::new("origin"),
                remote_refname: RefName::new("refs/heads/moo"),
                symref_target: RefName::new(""),
                upstream_refname: RefName::new("refs/remotes/origin/moo"),
            }
        )
    }
}
