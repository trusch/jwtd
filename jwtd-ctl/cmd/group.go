package cmd

import "github.com/spf13/cobra"

// groupCmd represents the group command
var groupCmd = &cobra.Command{
	Use:     "group",
	Aliases: []string{"groups"},
	Short:   "group related stuff",
	Long:    `This is the top level command for group related stuff.`,
}

func init() {
	RootCmd.AddCommand(groupCmd)
	groupCmd.PersistentFlags().StringP("name", "n", "", "group name")
}
