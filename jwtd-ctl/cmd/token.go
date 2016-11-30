package cmd

import "github.com/spf13/cobra"

// tokenCmd represents the token command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "token related stuff",
	Long:  `This contains token related subcommands (create and verify).`,
}

func init() {
	RootCmd.AddCommand(tokenCmd)
}
