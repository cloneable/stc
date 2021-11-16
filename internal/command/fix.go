package command

import "github.com/spf13/cobra"

func init() {
	fixCmd := &cobra.Command{
		Use: "fix [<branch> <base>]",
		// TODO: re-enable validation, add completion
		Args: cobra.RangeArgs(0, 2), // validBranchNames(stacker, 0, 2),
		//ValidArgsFunction: nil,
		RunE: func(cmd *cobra.Command, args []string) error {
			return stacker.Fix(cmd.Context(), args...)
		},
		DisableFlagsInUseLine: true,

		Short: "Adds, updates and deletes tracking refs if needed",
	}
	rootCmd.AddCommand(fixCmd)
}
