package command

import (
	"github.com/spf13/cobra"
)

func init() {
	deleteCmd := &cobra.Command{
		Use:  "delete <branch>",
		Args: validBranchNames(stacker, 1, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return stacker.Delete(cmd.Context(), args[0])
		},
		DisableFlagsInUseLine: true,
		Hidden:                true, // TODO: unhide when ready
	}
	rootCmd.AddCommand(deleteCmd)
}
