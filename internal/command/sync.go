package command

import (
	"fmt"
	"math"

	"github.com/spf13/cobra"
)

func init() {
	syncCmd := &cobra.Command{
		Use:  "sync [<branch>...]",
		Args: validBranchNames(stacker, 0, math.MaxInt),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("sync")
			return stacker.Sync(cmd.Context(), args...)
		},
		DisableFlagsInUseLine: true,
	}
	rootCmd.AddCommand(syncCmd)
}
