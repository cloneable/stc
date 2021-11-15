package command

import (
	"github.com/spf13/cobra"
)

func init() {
	syncCmd := &cobra.Command{
		Use:  "sync",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return stacker.Sync(cmd.Context())
		},
		DisableFlagsInUseLine: true,
	}
	rootCmd.AddCommand(syncCmd)
}
