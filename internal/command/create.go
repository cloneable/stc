package command

import (
	"fmt"

	"github.com/cloneable/stacker/internal/stacker"
	"github.com/spf13/cobra"
)

var (
	createCmd = &cobra.Command{
		Use:  "create <branch>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("create")
			s := stacker.Stacker{}
			return s.Create(args[0])
		},
		DisableFlagsInUseLine: true,
	}
)

func init() {
	rootCmd.AddCommand(createCmd)
}
