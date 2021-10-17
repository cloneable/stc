package command

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stacker",
	Short: "",
	Long:  "",
}

func Execute() error {
	return rootCmd.Execute()
}
