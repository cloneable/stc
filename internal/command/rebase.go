package command

import (
	"fmt"
	"math"

	"github.com/spf13/cobra"
)

func init() {
	rebaseCmd := &cobra.Command{
		Use:  "rebase [<branch>...]",
		Args: validBranchNames(stacker, 0, math.MaxInt),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("rebase")
			return stacker.Rebase(cmd.Context(), args...)
		},
		DisableFlagsInUseLine: true,
	}
	rootCmd.AddCommand(rebaseCmd)
}
