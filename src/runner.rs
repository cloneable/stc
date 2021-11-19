use crate::git::{Git, Status};
use ::std::assert_ne;
use ::std::option::Option::Some;
use ::std::process::Command;
use ::std::process::Stdio;
use ::std::result::Result::{self, Err, Ok};

pub struct Runner<'a> {
    gitpath: &'a str,
}

impl<'a> Runner<'a> {
    pub fn new(gitpath: &'a str) -> Self {
        Runner { gitpath }
    }
}

impl<'a> Git for Runner<'a> {
    fn exec(&self, args: &[&str]) -> Result<Status, Status> {
        let cmd = Command::new(self.gitpath)
            .args(args)
            .stdin(Stdio::null())
            .stdin(Stdio::piped())
            .stdout(Stdio::piped())
            .spawn()
            .expect("failed to start git");

        let output = cmd.wait_with_output().expect("failed to wait on git");
        if output.status.success() {
            Ok(Status::new(0, output.stdout, output.stderr))
        } else {
            if let Some(code) = output.status.code() {
                assert_ne!(code, 0);
                Err(Status::new(code, output.stdout, output.stderr))
            } else {
                Err(Status::new(1, output.stdout, output.stderr))
            }
        }
    }
}
