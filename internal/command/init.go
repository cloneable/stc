package command

import (
	"fmt"

	"github.com/cloneable/stacker/internal/stacker"
	"github.com/spf13/cobra"
)

var (
	initCmd = &cobra.Command{
		Use:  "init [--force]",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("init")
			s := stacker.Stacker{}
			return s.Init(forceFlag)
		},
		PersistentPreRunE:     overrideRepoValidation,
		DisableFlagsInUseLine: true,
	}
)

func init() {
	initCmd.Flags().BoolVar(&forceFlag, "force", false, "Source directory to read from")
	rootCmd.AddCommand(initCmd)
}
