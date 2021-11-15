package command

import (
	"github.com/spf13/cobra"
)

func init() {
	initCmd := &cobra.Command{
		Use:  "init",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return stacker.Init(cmd.Context())
		},
		PersistentPreRunE:     overrideRepoValidation,
		DisableFlagsInUseLine: true,

		Short: "Initializes the repo and tries to set stacker refs for any non-default branches.",
	}
	rootCmd.AddCommand(initCmd)
}
