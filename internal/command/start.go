package command

import (
	"github.com/spf13/cobra"
)

func init() {
	startCmd := &cobra.Command{
		Use:  "start <branch>",
		Args: validBranchNames(stacker, 1, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return stacker.Start(cmd.Context(), args[0])
		},
		DisableFlagsInUseLine: true,

		Short: "Starts a new branch off of current branch.",
	}
	rootCmd.AddCommand(startCmd)
}
