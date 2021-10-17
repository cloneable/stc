package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	deleteCmd = &cobra.Command{
		Use:  "delete <branch>",
		Args: validBranchNames(stacker, 1, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("delete")
			return stacker.Delete(cmd.Context(), args[0])
		},
		DisableFlagsInUseLine: true,
	}
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}
