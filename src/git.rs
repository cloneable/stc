use ::const_format::concatcp;
use ::csv::ReaderBuilder;
use ::serde::Deserialize;
use ::std::{
    borrow::Cow,
    clone::Clone,
    collections::HashMap,
    default::Default,
    error::Error,
    format,
    iter::{IntoIterator, Iterator},
    option::Option::{self, None},
    result::Result::{self, Err, Ok},
    string::{String, ToString},
    todo,
    vec::Vec,
    write,
};

// TODO: use ObjectName as type for const if possibe
pub const NON_EXISTANT_OBJECT: &str = "0000000000000000000000000000000000000000";

pub const STC_REF_PREFIX: &str = "refs/stc/";
pub const STC_BASE_REF_PREFIX: &str = concatcp!(STC_REF_PREFIX, "base/");
pub const STC_START_REF_PREFIX: &str = concatcp!(STC_REF_PREFIX, "start/");
pub const STC_REMOTE_REF_PREFIX: &str = concatcp!(STC_REF_PREFIX, "remote/");

pub const BRANCH_REF_PREFIX: &str = "refs/heads/";

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

    pub fn with(exitcode: i32) -> Self {
        Status {
            exitcode,
            stdout: Default::default(),
            stderr: Default::default(),
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

    fn snapshot(&self) -> Result<Repository, Status> {
        let status = self.exec(&["for-each-ref", "--format", FIELD_FORMATS.join(",").as_str()])?;
        let refs = parse_ref(status.stdout.as_slice()).map_err(|_err| Status::with(1))?;
        let head = refs.values().find(|r| r.head).map(|r| r.name.branchname());
        Ok(Repository { refs, head })
    }

    fn check_branchname<'a>(&self, name: &'a str) -> Result<BranchName<'a>, Status> {
        self.exec(&["check-ref-format", "--branch", name])?;
        Ok(BranchName(Cow::Owned(name.to_string())))
    }

    fn create_branch(&self, name: &BranchName, base: &BranchName) -> Result<(), Status> {
        self.exec(&["branch", "--create-reflog", name.as_str(), base.as_str()])
            .map(|_| {})
    }

    fn switch_branch(&self, b: &BranchName) -> Result<(), Status> {
        self.exec(&["switch", "--no-guess", b.as_str()]).map(|_| {})
    }

    fn create_symref(
        &self,
        name: &RefName,
        target: &RefName,
        reason: &'static str,
    ) -> Result<(), Status> {
        self.exec(&["symbolic-ref", "-m", reason, name.as_str(), target.as_str()])
            .map(|_| {})
    }

    fn delete_symref(&self, name: &RefName) -> Result<(), Status> {
        self.exec(&["symbolic-ref", "--delete", name.as_str()])
            .map(|_| {})
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
        .map(|_| {})
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
        .map(|_| {})
    }

    fn delete_ref(&self, name: &RefName, cur_commit: &ObjectName) -> Result<(), Status> {
        self.exec(&[
            "update-ref",
            "--no-deref",
            "-d",
            name.as_str(),
            cur_commit.as_str(),
        ])
        .map(|_| {})
    }

    fn rebase_onto(&self, name: &BranchName) -> Result<(), Status> {
        self.exec(&[
            "rebase",
            "--committer-date-is-author-date",
            "--onto",
            name.stc_base_refname().as_str(),
            name.stc_start_refname().as_str(),
            name.as_str(),
        ])
        .map(|_| {})
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
        .map(|_| {})
    }

    fn config_set(&self, key: &str, value: &str) -> Result<(), Status> {
        self.exec(&["config", "--local", key, value]).map(|_| {})
    }

    fn config_add(&self, key: &str, value: &str) -> Result<(), Status> {
        self.exec(&["config", "--local", "--add", key, value])
            .map(|_| {})
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
        self.exec(&["fetch", "--all", "--prune"]).map(|_| {})
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
}

#[derive(PartialEq, PartialOrd, Debug)]
pub struct BranchName<'a>(Cow<'a, String>);

impl<'a> BranchName<'a> {
    const fn new(name: String) -> Self {
        BranchName(Cow::Owned(name))
    }

    pub fn as_str(&self) -> &str {
        self.0.as_str()
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
pub struct RefName<'a>(Cow<'a, String>);

impl<'a> RefName<'a> {
    const fn new(name: String) -> Self {
        RefName(Cow::Owned(name))
    }

    pub fn as_str(&self) -> &str {
        self.0.as_str()
    }

    pub fn branchname(&self) -> BranchName<'a> {
        let (_, branchname) = self.0.rsplit_once("/").unwrap();
        BranchName::new(branchname.to_string())
    }
}

#[derive(Deserialize, PartialEq, PartialOrd, Clone, Debug)]
pub struct ObjectName<'a>(pub Cow<'a, String>);

impl<'a> ObjectName<'a> {
    pub const fn new(value: String) -> Self {
        ObjectName(Cow::Owned(value))
    }

    pub fn as_str(&self) -> &str {
        self.0.as_str()
    }
}

#[derive(Deserialize, PartialEq, PartialOrd, Debug)]
pub struct RemoteName<'a>(Cow<'a, String>);

impl<'a> RemoteName<'a> {
    const fn new(value: String) -> Self {
        RemoteName(Cow::Owned(value))
    }

    pub fn as_str(&self) -> &str {
        self.0.as_str()
    }
}

#[derive(Deserialize, PartialEq, Debug)]
#[serde(rename_all = "lowercase")]
pub enum RefType {
    Commit,
    Tree,
    Blob,
    Tag,
}

#[derive(Deserialize, PartialEq, Debug)]
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
            refs.get(&RefName::new("refs/heads/moo1".to_string()))
                .unwrap(),
            &Ref {
                name: RefName::new("refs/heads/moo1".to_string()),
                head: true,
                objectname: ObjectName::new("123abc".to_string()),
                objecttype: RefType::Commit,
                track: String::from("<>"),
                remote: RemoteName::new("origin".to_string()),
                remote_refname: RefName::new("refs/heads/moo".to_string()),
                symref_target: RefName::new("".to_string()),
                upstream_refname: RefName::new("refs/remotes/origin/moo".to_string()),
            }
        )
    }
}
