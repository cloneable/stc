package command

import (
	"fmt"

	"github.com/cloneable/stacker/internal/stacker"
	"github.com/spf13/cobra"
)

var (
	showCmd = &cobra.Command{
		Use:  "show",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("show")
			s := stacker.Stacker{}
			return s.Show()
		},
		DisableFlagsInUseLine: true,
	}
)

func init() {
	rootCmd.AddCommand(showCmd)
}
