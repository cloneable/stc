package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	startCmd := &cobra.Command{
		Use:  "start <branch>",
		Args: validBranchNames(stacker, 1, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("start")
			return stacker.Start(cmd.Context(), args[0])
		},
		DisableFlagsInUseLine: true,
	}
	rootCmd.AddCommand(startCmd)
}
