use crate::git::{Git, Result};
use ::std::option::Option::Some;
use ::std::process::exit;
use ::std::process::Command;
use ::std::process::Stdio;

pub struct Runner<'a> {
    gitpath: &'a str,
}

impl<'a> Runner<'a> {
    pub fn new(gitpath: &'a str) -> Self {
        Runner { gitpath }
    }
}

impl<'a> Git for Runner<'a> {
    fn exec(&self, args: &[&str]) -> Result {
        let cmd = Command::new(self.gitpath)
            .args(args)
            .stdin(Stdio::null())
            .stdin(Stdio::piped())
            .stdout(Stdio::piped())
            .spawn()
            .expect("failed to start git");

        let output = cmd.wait_with_output().expect("failed to wait on git");
        if let Some(code) = output.status.code() {
            Result::new(code, output.stdout, output.stderr)
        } else {
            Result::new(
                if output.status.success() { 0 } else { 1 },
                output.stdout,
                output.stderr,
            )
        }
    }

    fn fail(&self) -> ! {
        exit(1)
    }
}
