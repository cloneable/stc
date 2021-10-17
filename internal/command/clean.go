package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cleanCmd = &cobra.Command{
		Use:  "clean [--force] [<branch>...]",
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("clean")
			return stacker.Clean(cmd.Context(), forceFlag, args...)
		},
		DisableFlagsInUseLine: true,
	}
)

func init() {
	cleanCmd.Flags().BoolVar(&forceFlag, "force", false, "Source directory to read from")
	rootCmd.AddCommand(cleanCmd)
}
