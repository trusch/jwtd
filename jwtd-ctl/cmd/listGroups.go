package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// listGroupsCmd represents the listGroups command
var listGroupsCmd = &cobra.Command{
	Use:   "list",
	Short: "List groups",
	Long:  `This lists all known groups and dumps them as yaml.`,
	Run: func(cmd *cobra.Command, args []string) {
		database := getDB()
		groups, err := database.ListGroups()
		if err != nil {
			log.Fatal(err)
		}
		bs, err := yaml.Marshal(groups)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(bs))
	},
}

func init() {
	groupCmd.AddCommand(listGroupsCmd)
}
