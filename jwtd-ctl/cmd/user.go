package cmd

import "github.com/spf13/cobra"

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:     "user",
	Aliases: []string{"users"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func init() {
	RootCmd.AddCommand(userCmd)
	userCmd.PersistentFlags().StringP("name", "n", "", "user name")
}
