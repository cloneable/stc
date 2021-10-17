package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	createCmd = &cobra.Command{
		Use:  "create <branch>",
		Args: validBranchNames(stacker, 1, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("create")
			return stacker.Create(cmd.Context(), args[0])
		},
		DisableFlagsInUseLine: true,
	}
)

func init() {
	rootCmd.AddCommand(createCmd)
}
