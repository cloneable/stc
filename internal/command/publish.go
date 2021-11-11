package command

import (
	"fmt"
	"math"

	"github.com/spf13/cobra"
)

func init() {
	publishCmd := &cobra.Command{
		Use:  "publish [<branch>...]",
		Args: validBranchNames(stacker, 1, math.MaxInt),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("publish")
			return stacker.Publish(cmd.Context(), args...)
		},
		DisableFlagsInUseLine: true,
		Hidden:                true, // TODO: unhide when ready
	}
	rootCmd.AddCommand(publishCmd)
}
