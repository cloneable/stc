use crate::git;
use ::std::{
    assert_ne,
    collections::HashMap,
    option::Option::Some,
    path::PathBuf,
    process::{Command, Stdio},
    result::Result::{self, Err, Ok},
};

pub struct Runner<'a> {
    gitpath: &'a str,
    workdir: PathBuf,
    env: HashMap<&'a str, &'a str>,
}

impl<'a> Runner<'a> {
    pub fn new(gitpath: &'a str) -> Self {
        Runner {
            gitpath,
            workdir: ::std::env::current_dir().expect("cannot determine current working directory"),
            env: HashMap::<&'a str, &'a str>::new(),
        }
    }
}

impl<'a> git::Git for Runner<'a> {
    fn exec(&self, args: &[&str]) -> Result<git::Status, git::Status> {
        let cmd = Command::new(self.gitpath)
            .args(args)
            .current_dir(&self.workdir)
            .envs(self.env.iter())
            .stdin(Stdio::null())
            .stdout(Stdio::piped())
            .stderr(Stdio::piped())
            .spawn()
            .expect("failed to start git");

        let output = cmd.wait_with_output().expect("failed to wait on git");
        if output.status.success() {
            ::std::eprintln!("[OK] git {:?}", args);
            Ok(git::Status::new(0, output.stdout, output.stderr))
        } else if let Some(code) = output.status.code() {
            assert_ne!(code, 0);
            ::std::eprintln!("[ERR {:?}] git {:?}", code, args);
            Err(git::Status::new(code, output.stdout, output.stderr))
        } else {
            ::std::eprintln!("[ERR] git {:?}", args);
            Err(git::Status::new(1, output.stdout, output.stderr))
        }
    }
}
