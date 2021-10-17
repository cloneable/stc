package command

import (
	"fmt"

	"github.com/cloneable/stacker/internal/stacker"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("init")
		return stacker.Init()
	},
	PersistentPreRunE: overrideRepoValidation,
}

func init() {
	rootCmd.AddCommand(initCmd)
}
