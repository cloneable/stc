use crate::git;
use ::anyhow::Result;
use ::std::{
    assert_ne,
    collections::HashMap,
    convert::Into,
    option::Option::Some,
    path::PathBuf,
    process::{Command, Stdio},
    result::Result::{Err, Ok},
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
        let cmd = Command::new(self.gitpath)
            .args(args)
            .current_dir(&self.workdir)
            .envs(self.env.iter())
            .stdin(Stdio::null())
            .stdout(Stdio::piped())
            .stderr(Stdio::piped())
            .spawn()?;

        let output = cmd.wait_with_output()?;
        if output.status.success() {
            ::std::eprintln!("[OK] git {:?}", args);
            Ok(git::ExecStatus::from(0, output.stdout, output.stderr))
        } else if let Some(code) = output.status.code() {
            assert_ne!(code, 0);
            ::std::eprintln!("[ERR {:?}] git {:?}", code, args);
            Err(git::ExecStatus::from(code, output.stdout, output.stderr).into())
        } else {
            ::std::eprintln!("[ERR] git {:?}", args);
            Err(git::ExecStatus::from(1, output.stdout, output.stderr).into())
        }
    }
}
