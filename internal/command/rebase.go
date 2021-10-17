package command

import (
	"fmt"
	"math"

	"github.com/spf13/cobra"
)

var (
	rebaseCmd = &cobra.Command{
		Use:  "rebase [<branch>...]",
		Args: validBranchNames(stacker, 0, math.MaxInt),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("rebase")
			return stacker.Rebase(cmd.Context(), args...)
		},
		DisableFlagsInUseLine: true,
	}
)

func init() {
	rootCmd.AddCommand(rebaseCmd)
}
