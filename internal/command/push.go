package command

import (
	"github.com/spf13/cobra"
)

func init() {
	pushCmd := &cobra.Command{
		Use:  "push",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return stacker.Push(cmd.Context())
		},
		DisableFlagsInUseLine: true,
	}
	rootCmd.AddCommand(pushCmd)
}
