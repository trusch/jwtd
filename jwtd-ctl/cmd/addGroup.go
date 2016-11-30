package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/trusch/jwtd/db"
)

// addGroupCmd represents the addGroup command
var addGroupCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a group",
	Long:  `This adds a group to your jwtd server.`,
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
		err := database.CreateGroup(name, []*db.AccessRight{})
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	groupCmd.AddCommand(addGroupCmd)
}
