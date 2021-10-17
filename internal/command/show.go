package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	showCmd = &cobra.Command{
		Use:  "show",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("show")
			return stacker.Show(cmd.Context())
		},
		DisableFlagsInUseLine: true,
	}
)

func init() {
	rootCmd.AddCommand(showCmd)
}
