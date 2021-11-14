package command

import (
	"github.com/spf13/cobra"
)

func init() {
	cleanCmd := &cobra.Command{
		Use:  "clean",      // "clean [--force] [<branch>...]",
		Args: cobra.NoArgs, // validBranchNames(stacker, 0, math.MaxInt),
		RunE: func(cmd *cobra.Command, args []string) error {
			return stacker.Clean(cmd.Context(), forceFlag, args...)
		},
		DisableFlagsInUseLine: true,

		Short: "Cleans any stacker related refs and settings from repo.",
	}
	// cleanCmd.Flags().BoolVar(&forceFlag, "force", false, "Source directory to read from")
	rootCmd.AddCommand(cleanCmd)
}
