#![no_implicit_prelude]
#![allow(missing_docs)] // TODO: change to warn/deny
#![allow(dead_code)] // TODO: remove

use ::clap::{self, Parser, Subcommand};
use ::std::option::Option::{self, None, Some};
use ::std::result::Result;
use ::std::string::String;

mod git;
mod runner;
mod stc;

#[derive(Parser)]
#[clap(about, version)]
struct Root {
    #[clap(subcommand)]
    subcommand: Command,
}

#[derive(Subcommand)]
#[clap()]
enum Command {
    #[clap(
        name = "clean",
        about = "Cleans any stacker related refs and settings from repo.",
        override_usage = "stacker clean"
    )]
    Clean,

    #[clap(
        name = "fix",
        about = "Adds, updates and deletes tracking refs if needed.",
        override_usage = "stacker fix [<branch> [<base]]"
    )]
    Fix {
        #[clap(name = "branch", about = "name of the branch to fix")]
        branch: Option<String>,
        #[clap(name = "base", about = "name of the base branch")]
        base: Option<String>,
    },

    #[clap(
        name = "init",
        about = "Initializes the repo and tries to set stacker refs for any non-default branches.",
        override_usage = "stacker init"
    )]
    Init,

    #[clap(
        name = "push",
        about = "Sets remote branch head to what local branch head points to.",
        override_usage = "stacker push"
    )]
    Push,

    #[clap(
        name = "rebase",
        about = "Rebases current branch on top of its base branch.",
        override_usage = "stacker rebase"
    )]
    Rebase,

    #[clap(
        name = "start",
        about = "Starts a new branch off of current branch.",
        override_usage = "stacker start <branch>"
    )]
    Start {
        #[clap(name = "branch", about = "name of the new branch to create")]
        branch: String,
    },

    #[clap(
        name = "sync",
        about = "Fetches all branches and tags and prunes deleted ones.",
        override_usage = "stacker sync"
    )]
    Sync,
}

fn main() -> Result<(), git::Status> {
    let root = Root::parse();
    let runner = runner::Runner::new("git");
    let stc = stc::STC::new(runner);
    match root.subcommand {
        Command::Clean => stc.clean(),
        Command::Fix { branch, base } => stc.fix(branch, base),
        Command::Init => stc.init(),
        Command::Push => stc.push(),
        Command::Rebase => stc.rebase(),
        Command::Start { branch } => stc.start(branch),
        Command::Sync => stc.sync(),
    }
}
