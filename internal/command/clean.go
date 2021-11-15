package command

import (
	"github.com/spf13/cobra"
)

func init() {
	cleanCmd := &cobra.Command{
		Use:  "clean",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return stacker.Clean(cmd.Context())
		},
		DisableFlagsInUseLine: true,

		Short: "Cleans any stacker related refs and settings from repo.",
	}
	rootCmd.AddCommand(cleanCmd)
}
