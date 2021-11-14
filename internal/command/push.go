package command

import (
	"github.com/spf13/cobra"
)

func init() {
	pushCmd := &cobra.Command{
		Use:  "push",       // "push [<branch>...]",
		Args: cobra.NoArgs, // validBranchNames(stacker, 0, math.MaxInt),
		RunE: func(cmd *cobra.Command, args []string) error {
			return stacker.Push(cmd.Context(), args...)
		},
		DisableFlagsInUseLine: true,
	}
	rootCmd.AddCommand(pushCmd)
}
