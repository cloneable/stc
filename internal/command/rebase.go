package command

import (
	"github.com/spf13/cobra"
)

func init() {
	rebaseCmd := &cobra.Command{
		Use:  "rebase",     //"rebase [<branch>...]",
		Args: cobra.NoArgs, // validBranchNames(stacker, 0, math.MaxInt),
		RunE: func(cmd *cobra.Command, args []string) error {
			return stacker.Rebase(cmd.Context(), args...)
		},
		DisableFlagsInUseLine: true,
	}
	rootCmd.AddCommand(rebaseCmd)
}
