package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// delUserCmd represents the delUser command
var delUserCmd = &cobra.Command{
	Use:   "del",
	Short: "delete a user from db",
	Long:  `This deletes a user from the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		database := getDB()
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			if len(args) > 0 {
				name = args[0]
			} else {
				log.Fatal("specify --name")
			}
		}
		err := database.DelUser(name)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	userCmd.AddCommand(delUserCmd)
}
