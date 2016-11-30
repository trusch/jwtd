package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// delGroupCmd represents the delGroup command
var delGroupCmd = &cobra.Command{
	Use:   "del",
	Short: "delete a group",
	Long:  `This deletes a group from the db.`,
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
		err := database.DelGroup(name)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	groupCmd.AddCommand(delGroupCmd)
}
