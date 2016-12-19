package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// addRightCmd represents the addRight command
var addRightCmd = &cobra.Command{
	Use:   "add-right",
	Short: "add an access right to a group",
	Long:  `This adds an access right to a group.`,
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
		service, _ := cmd.Flags().GetString("service")
		key, _ := cmd.Flags().GetString("key")
		value, _ := cmd.Flags().GetString("value")
		if service == "" || key == "" || value == "" {
			log.Fatal("specifiy --service, --key and --value")
		}
		group, err := database.GetGroup(project, name)
		if err != nil {
			log.Fatal(err)
		}
		if labels, ok := group.Rights[service]; ok {
			labels[key] = value
		} else {
			group.Rights[service] = map[string]string{key: value}
		}
		err = database.UpdateGroup(group)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	groupCmd.AddCommand(addRightCmd)
	addRightCmd.Flags().StringP("service", "s", "", "affected service")
	addRightCmd.Flags().StringP("key", "k", "", "label key")
	addRightCmd.Flags().StringP("value", "v", "", "label value")
}
