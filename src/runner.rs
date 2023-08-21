use crate::git;
use color_eyre::Result;
use std::{
    collections::HashMap,
    path::PathBuf,
    process::{Command, Stdio},
};

pub struct Runner<'a> {
    gitpath: &'a str,
    workdir: PathBuf,
    env: HashMap<&'a str, &'a str>,
}

impl<'a> Runner<'a> {
    pub fn new(gitpath: &'a str) -> Result<Self> {
        Ok(Runner {
            gitpath,
            workdir: ::std::env::current_dir()?,
            env: HashMap::<&'a str, &'a str>::new(),
        })
    }
}

impl<'a> git::Git for Runner<'a> {
    fn exec(&self, args: &[&str]) -> Result<git::ExecStatus> {
        let output = Command::new(self.gitpath)
            .args(args)
            .current_dir(&self.workdir)
            .envs(self.env.iter())
            .stdin(Stdio::null())
            .stdout(Stdio::piped())
            .stderr(Stdio::piped())
            .output()?;

        if output.status.success() {
            Ok(git::ExecStatus::new(0, &output.stdout, &output.stderr))
        } else if let Some(code) = output.status.code() {
            ::std::assert_ne!(code, 0);
            Err(git::ExecStatus::new(code, &output.stdout, &output.stderr).into())
        } else {
            Err(git::ExecStatus::new(1, &output.stdout, &output.stderr).into())
        }
    }
}
