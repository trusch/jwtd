package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// delRightCmd represents the delRight command
var delRightCmd = &cobra.Command{
	Use:   "del-right",
	Short: "Delete a right from a group",
	Long:  `This deletes an accessright from a group. The right is specified by --service and --subject.`,
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
		service, _ := cmd.Flags().GetString("service")
		key, _ := cmd.Flags().GetString("key")
		if service == "" || key == "" {
			log.Fatal("specifiy --service and --key")
		}
		group, err := database.GetGroup(name)
		if err != nil {
			log.Fatal(err)
		}
		if _, ok := group.Rights[service][key]; ok {
			delete(group.Rights[service], key)
		} else {
			log.Fatal("no such label key")
		}
		err = database.UpdateGroup(group)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	groupCmd.AddCommand(delRightCmd)
	delRightCmd.Flags().StringP("service", "s", "", "affected service")
	delRightCmd.Flags().StringP("key", "k", "", "label key")
}
