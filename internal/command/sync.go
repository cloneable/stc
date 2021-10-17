package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	syncCmd = &cobra.Command{
		Use:  "sync [<branch>...]",
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("sync")
			return stacker.Sync(cmd.Context(), args...)
		},
		DisableFlagsInUseLine: true,
	}
)

func init() {
	rootCmd.AddCommand(syncCmd)
}
