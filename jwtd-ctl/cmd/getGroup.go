package cmd

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

// getGroupCmd represents the getGroup command
var getGroupCmd = &cobra.Command{
	Use:   "get",
	Short: "get a group",
	Long:  `This dumps group information about a given group from the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		database := getDB()
		project, _ := cmd.Flags().GetString("project")
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			if len(args) > 0 {
				name = args[0]
			} else {
				log.Fatal("specify --name")
			}
		}
		group, err := database.GetGroup(project, name)
		if err != nil {
			log.Fatal(err)
		}
		bs, err := yaml.Marshal(group)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(bs))
	},
}

func init() {
	groupCmd.AddCommand(getGroupCmd)
}
