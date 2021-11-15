package command

import "github.com/spf13/cobra"

func init() {
	fixCmd := &cobra.Command{
		Use:  "fix",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return stacker.Fix(cmd.Context())
		},
		DisableFlagsInUseLine: true,

		Short: "Adds, updates and deletes tracking refs if needed",
	}
	rootCmd.AddCommand(fixCmd)
}
