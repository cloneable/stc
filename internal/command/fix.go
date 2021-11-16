package command

import "github.com/spf13/cobra"

func init() {
	fixCmd := &cobra.Command{
		Use:  "fix [<branch> <base>]",
		Args: validBranchNames(stacker, 0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return stacker.Fix(cmd.Context(), args...)
		},
		DisableFlagsInUseLine: true,

		Short: "Adds, updates and deletes tracking refs if needed",
	}
	rootCmd.AddCommand(fixCmd)
}
