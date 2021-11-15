package command

import (
	"context"
	"errors"
	"os"

	stackerpkg "github.com/cloneable/stacker/internal/stacker"
	"github.com/spf13/cobra"
)

func overrideRepoValidation(*cobra.Command, []string) error { return nil }

func validBranchNames(s *stackerpkg.Stacker, min, max int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < min || len(args) > max {
			return errors.New("invalid number of branch names")
		}
		if !s.ValidBranchNames(args...) {
			return errors.New("invalid branch name")
		}
		return nil
	}
}

// TODO: fine a better way to init Stacker
func must(v string, err error) string {
	if err != nil {
		panic(err)
	}
	return v
}

var (
	workdir = must(os.Getwd())
	stacker = stackerpkg.New(workdir)

	printCommandsFlag bool

	rootCmd = &cobra.Command{
		Use: "stacker <command>",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// TODO: do refs integrity checks
			return nil
		},
		DisableFlagsInUseLine: true,
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&printCommandsFlag, "print-commands", "x", false, "Print commands that are executed")
}

func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}
