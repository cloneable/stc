package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	initCmd := &cobra.Command{
		Use:  "init [--force]",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("init")
			return stacker.Init(cmd.Context(), forceFlag)
		},
		PersistentPreRunE:     overrideRepoValidation,
		DisableFlagsInUseLine: true,
	}
	initCmd.Flags().BoolVar(&forceFlag, "force", false, "Source directory to read from")
	rootCmd.AddCommand(initCmd)
}
