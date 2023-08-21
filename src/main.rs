#![allow(missing_docs)] // TODO: change to warn/deny
#![allow(dead_code)] // TODO: remove

use clap::{Parser, Subcommand};

mod git;
mod runner;
mod stc;

#[derive(Parser)]
#[command(about, version, override_usage = "stc <command>")]
struct Root {
    #[command(subcommand)]
    subcommand: Command,
}

#[derive(Subcommand)]
enum Command {
    /// Cleans any stc related refs and settings from repo.
    Clean,

    /// Adds, updates and deletes tracking refs if needed.
    Fix {
        /// name of the branch to fix
        #[arg(name = "branch")]
        branch: Option<String>,
        /// name of the base branch
        #[arg(name = "base")]
        base: Option<String>,
    },

    /// Initializes the repo and tries to set stc refs for any non-default branches.
    Init,

    /// Sets remote branch head to what local branch head points to.
    Push,

    /// Rebases current branch on top of its base branch.
    Rebase,

    /// Starts a new branch off of current branch.
    Start {
        /// name of the new branch to create
        #[arg(name = "branch")]
        branch: String,
    },

    /// Fetches all branches and tags and prunes deleted ones.
    Sync,
}

fn main() -> color_eyre::Result<()> {
    let root = Root::parse();
    let runner = runner::Runner::new("git")?;
    let stc = stc::Stc::new(runner);
    match root.subcommand {
        Command::Clean => stc.clean(),
        Command::Fix { branch, base } => stc.fix(branch, base),
        Command::Init => stc.init(),
        Command::Push => stc.push(),
        Command::Rebase => stc.rebase(),
        Command::Start { branch } => stc.start(&branch),
        Command::Sync => stc.sync(),
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn verify_commands() {
        use clap::CommandFactory;
        Root::command().debug_assert()
    }
}
