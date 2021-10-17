package command

import (
	"context"
	"fmt"

	stackerpkg "github.com/cloneable/stacker/internal/stacker"
	"github.com/spf13/cobra"
)

func overrideRepoValidation(*cobra.Command, []string) error { return nil }

var (
	stacker = stackerpkg.New()

	// used by some subcommands
	forceFlag bool

	rootCmd = &cobra.Command{
		Use: "stacker <command>",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("validate stacker refs")
			return nil
		},
		DisableFlagsInUseLine: true,
	}
)

func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}
