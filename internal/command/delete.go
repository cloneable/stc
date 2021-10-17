package command

import (
	"fmt"

	"github.com/cloneable/stacker/internal/stacker"
	"github.com/spf13/cobra"
)

var (
	deleteCmd = &cobra.Command{
		Use:  "delete <branch>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("delete")
			s := stacker.Stacker{}
			return s.Delete(args[0])
		},
		DisableFlagsInUseLine: true,
	}
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}
