package command

import (
	"github.com/spf13/cobra"
)

func init() {
	initCmd := &cobra.Command{
		Use:  "init [--force]",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return stacker.Init(cmd.Context(), forceFlag)
		},
		PersistentPreRunE:     overrideRepoValidation,
		DisableFlagsInUseLine: true,

		Short: "Initializes the repo and tries to set stacker refs for any non-default branches.",
	}
	initCmd.Flags().BoolVar(&forceFlag, "force", false, "Force removal or overwriting of broken stacker refs")
	rootCmd.AddCommand(initCmd)
}
