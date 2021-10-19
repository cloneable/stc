package command

import (
	"fmt"
	"math"

	"github.com/spf13/cobra"
)

func init() {
	cleanCmd := &cobra.Command{
		Use:  "clean [--force] [<branch>...]",
		Args: validBranchNames(stacker, 0, math.MaxInt),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("clean")
			return stacker.Clean(cmd.Context(), forceFlag, args...)
		},
		DisableFlagsInUseLine: true,
	}
	cleanCmd.Flags().BoolVar(&forceFlag, "force", false, "Source directory to read from")
	rootCmd.AddCommand(cleanCmd)
}
