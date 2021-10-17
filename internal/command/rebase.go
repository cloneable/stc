package command

import (
	"fmt"

	"github.com/cloneable/stacker/internal/stacker"
	"github.com/spf13/cobra"
)

var (
	rebaseCmd = &cobra.Command{
		Use:  "rebase [<branch>...]",
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("rebase")
			s := stacker.Stacker{}
			return s.Rebase(args...)
		},
		DisableFlagsInUseLine: true,
	}
)

func init() {
	rootCmd.AddCommand(rebaseCmd)
}
