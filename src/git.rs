use ::const_format::concatcp;
use ::csv::ReaderBuilder;
use ::serde::Deserialize;
use ::std::borrow::Cow;
use ::std::clone::Clone;
use ::std::collections::HashMap;
use ::std::default::Default;
use ::std::error::Error;
use ::std::format;
use ::std::iter::IntoIterator;
use ::std::iter::Iterator;
use ::std::option::Option::{self, None};
use ::std::result::Result::{self, Err, Ok};
use ::std::string::String;
use ::std::string::ToString;
use ::std::todo;
use ::std::vec::Vec;
use ::std::write;

// TODO: use ObjectName as type for const if possibe
pub const NON_EXISTANT_OBJECT: &'static str = "0000000000000000000000000000000000000000";

pub const STACKER_REF_PREFIX: &'static str = "refs/stacker/";
pub const STACKER_BASE_REF_PREFIX: &'static str = concatcp!(STACKER_REF_PREFIX, "base/");
pub const STACKER_START_REF_PREFIX: &'static str = concatcp!(STACKER_REF_PREFIX, "start/");
pub const STACKER_REMOTE_REF_PREFIX: &'static str = concatcp!(STACKER_REF_PREFIX, "remote/");

pub const BRANCH_REF_PREFIX: &'static str = "refs/heads/";

#[derive(Debug)]
pub struct Status {
    pub exitcode: i32,
    pub stdout: Vec<u8>,
    pub stderr: Vec<u8>,
}

impl Status {
    pub fn new(exitcode: i32, stdout: Vec<u8>, stderr: Vec<u8>) -> Self {
        Status {
            exitcode,
            stdout,
            stderr,
        }
    }
}

impl ::std::fmt::Display for Status {
    fn fmt(&self, f: &mut ::std::fmt::Formatter<'_>) -> ::std::fmt::Result {
        write!(f, "SuperError is here!")
    }
}

impl Error for Status {
    fn source(&self) -> Option<&(dyn Error + 'static)> {
        // XXX
        None //Some(&self.side)
    }
}

pub trait Git {
    fn exec(&self, args: &[&str]) -> Result<Status, Status>;

    fn head(&self) -> Result<BranchName, Status> {
        todo!()
    }

    fn snapshot(&self) -> Result<HashMap<RefName, Ref>, Status> {
        let status = self.exec(&["for-each-ref", "--format", FIELD_FORMATS.join(",").as_str()])?;
        parse_ref(status.stdout.as_slice())
            .map_err(move |_err| Status::new(1, Default::default(), Default::default()))
    }

    fn get_ref(&self, _name: &RefName) -> Result<Ref, Status> {
        todo!()
    }

    fn branch(&self, _name: &String) -> Result<BranchName, Status> {
        todo!()
    }

    fn check_branchname<'a>(&self, name: &'a String) -> Result<BranchName<'a>, Status> {
        self.exec(&["check-ref-format", "--branch", name])?;
        Ok(BranchName(Cow::Borrowed(name)))
    }

    fn create_branch(&self, name: &BranchName, base: &BranchName) -> Result<(), Status> {
        self.exec(&["branch", "--create-reflog", name.as_str(), base.as_str()])
            .map(move |_| -> () {})
    }

    fn switch_branch(&self, b: &BranchName) -> Result<(), Status> {
        self.exec(&["switch", "--no-guess", b.as_str()])
            .map(move |_| -> () {})
    }

    fn create_symref(
        &self,
        name: &RefName,
        target: &RefName,
        reason: &'static str,
    ) -> Result<(), Status> {
        self.exec(&["symbolic-ref", "-m", reason, name.as_str(), target.as_str()])
            .map(move |_| -> () {})
    }

    fn delete_symref(&self, name: &RefName) -> Result<(), Status> {
        self.exec(&["symbolic-ref", "--delete", name.as_str()])
            .map(move |_| -> () {})
    }

    fn create_ref(&self, name: &RefName, commit: &ObjectName) -> Result<(), Status> {
        self.exec(&[
            "update-ref",
            "--no-deref",
            "--create-reflog",
            name.as_str(),
            commit.as_str(),
            NON_EXISTANT_OBJECT,
        ])
        .map(move |_| -> () {})
    }

    fn update_ref(
        &self,
        name: &RefName,
        new_commit: &ObjectName,
        cur_commit: &ObjectName,
    ) -> Result<(), Status> {
        self.exec(&[
            "update-ref",
            "--no-deref",
            "--create-reflog",
            name.as_str(),
            new_commit.as_str(),
            cur_commit.as_str(),
        ])
        .map(move |_| -> () {})
    }

    fn delete_ref(&self, name: &RefName, cur_commit: &ObjectName) -> Result<(), Status> {
        self.exec(&[
            "update-ref",
            "--no-deref",
            "-d",
            name.as_str(),
            cur_commit.as_str(),
        ])
        .map(move |_| -> () {})
    }

    fn rebase_onto(&self, name: &BranchName) -> Result<(), Status> {
        self.exec(&[
            "rebase",
            "--committer-date-is-author-date",
            "--onto",
            name.stacker_base_refname().as_str(),
            name.stacker_start_refname().as_str(),
            name.as_str(),
        ])
        .map(move |_| -> () {})
    }

    fn push(
        &self,
        name: &BranchName,
        remote: &RemoteName,
        expect: &ObjectName,
    ) -> Result<(), Status> {
        self.exec(&[
            "push",
            "--set-upstream",
            format!("--force-with-lease={}:{}", name.as_str(), expect.as_str()).as_str(),
            remote.as_str(),
            format!("{}:{}", name.as_str(), name.as_str()).as_str(),
        ])
        .map(move |_| -> () {})
    }

    fn config_set(&self, key: &str, value: &str) -> Result<(), Status> {
        self.exec(&["config", "--local", key, value])
            .map(move |_| -> () {})
    }

    fn config_add(&self, key: &str, value: &str) -> Result<(), Status> {
        self.exec(&["config", "--local", "--add", key, value])
            .map(move |_| -> () {})
    }

    fn config_unset_pattern(&self, key: &str, pattern: &str) -> Result<(), Status> {
        match self.exec(&[
            "config",
            "--local",
            "--fixed-value",
            "--unset-all",
            key,
            pattern,
        ]) {
            // 5 means the nothing matched.
            Err(status) if status.exitcode != 5 => Err(status),
            _ => Ok(()),
        }
    }

    fn fetch_all_prune(&self) -> Result<(), Status> {
        self.exec(&["fetch", "--all", "--prune"])
            .map(move |_| -> () {})
    }

    fn tracked_branches(&self) -> Result<Vec<BranchName>, Status> {
        todo!()
    }

    fn forkpoint(&self, base: &RefName, branch: &RefName) -> Result<ObjectName, Status> {
        self.exec(&["merge-base", "--fork-point", base.as_str(), branch.as_str()])
            .map(move |status| {
                // TODO: handle not found
                ObjectName(Cow::Owned(
                    String::from_utf8_lossy(&status.stdout).to_string(),
                ))
            })
    }
}

#[derive(PartialEq, PartialOrd, Debug)]
pub struct BranchName<'a>(Cow<'a, String>);

impl<'a> BranchName<'a> {
    pub fn as_str(&self) -> &str {
        self.0.as_str()
    }

    pub fn refname(&self) -> RefName {
        RefName(Cow::Owned(BRANCH_REF_PREFIX.to_string() + &self.0))
    }

    pub fn stacker_base_refname(&self) -> RefName {
        RefName(Cow::Owned(STACKER_BASE_REF_PREFIX.to_string() + &self.0))
    }

    pub fn stacker_start_refname(&self) -> RefName {
        RefName(Cow::Owned(STACKER_START_REF_PREFIX.to_string() + &self.0))
    }

    pub fn stacker_remote_refname(&self) -> RefName {
        RefName(Cow::Owned(STACKER_REMOTE_REF_PREFIX.to_string() + &self.0))
    }
}

#[derive(Deserialize, Clone, Hash, Eq, PartialEq, Ord, PartialOrd, Debug)]
pub struct RefName<'a>(Cow<'a, String>);

impl<'a> RefName<'a> {
    pub fn as_str(&self) -> &str {
        self.0.as_str()
    }
}

#[derive(Deserialize, PartialEq, PartialOrd, Clone, Debug)]
pub struct ObjectName<'a>(Cow<'a, String>);

impl<'a> ObjectName<'a> {
    pub const fn new(value: String) -> ObjectName<'a> {
        ObjectName(Cow::Owned(value))
    }

    pub fn as_str(&self) -> &str {
        self.0.as_str()
    }
}

#[derive(Deserialize, PartialEq, PartialOrd, Debug)]
pub struct RemoteName<'a>(Cow<'a, String>);

impl<'a> RemoteName<'a> {
    pub fn as_str(&self) -> &str {
        self.0.as_str()
    }
}

#[derive(Deserialize, Debug)]
#[serde(rename_all = "lowercase")]
pub enum RefType {
    Commit,
    Tree,
    Blob,
    Tag,
}

#[derive(Deserialize, Debug)]
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

const FIELD_FORMATS: [&'static str; 9] = [
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
    let mut reader = ReaderBuilder::new().delimiter(b',').from_reader(csv);
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

    #[test]
    fn test_parse_ref() {
        let csv = "\
name,head,objectname,objecttype,track,remote,remote_refname,symref_target,upstream_refname
refs/heads/moo1,true,123abc,commit,<>,origin,refs/heads/moo,,refs/remotes/origin/moo
refs/heads/moo2,true,123abc,commit,<>,origin,refs/heads/moo,,refs/remotes/origin/moo
refs/heads/moo3,true,123abc,commit,<>,origin,refs/heads/moo,,refs/remotes/origin/moo
";
        let refs = parse_ref(csv.as_bytes()).expect("cannot parse");
        assert_eq!(refs.len(), 3);
    }
}
