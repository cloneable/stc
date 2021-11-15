package command

import (
	"github.com/spf13/cobra"
)

func init() {
	rebaseCmd := &cobra.Command{
		Use:  "rebase",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return stacker.Rebase(cmd.Context())
		},
		DisableFlagsInUseLine: true,
	}
	rootCmd.AddCommand(rebaseCmd)
}
