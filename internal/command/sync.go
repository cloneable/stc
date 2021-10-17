package command

import (
	"fmt"

	"github.com/cloneable/stacker/internal/stacker"
	"github.com/spf13/cobra"
)

var (
	syncCmd = &cobra.Command{
		Use:  "sync [<branch>...]",
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("sync")
			s := stacker.Stacker{}
			return s.Sync(args...)
		},
		DisableFlagsInUseLine: true,
	}
)

func init() {
	rootCmd.AddCommand(syncCmd)
}
