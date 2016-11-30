package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/trusch/jwtd/db"
)

// addRightCmd represents the addRight command
var addRightCmd = &cobra.Command{
	Use:   "add-right",
	Short: "add an access right to a group",
	Long:  `This adds an access right to a group.`,
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
		subject, _ := cmd.Flags().GetString("subject")
		if service == "" || subject == "" {
			log.Fatal("specifiy --service and --subject")
		}
		group, err := database.GetGroup(name)
		if err != nil {
			log.Fatal(err)
		}
		group.Rights = append(group.Rights, &db.AccessRight{Service: service, Subject: subject})
		err = database.UpdateGroup(group)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	groupCmd.AddCommand(addRightCmd)
	addRightCmd.Flags().StringP("service", "s", "", "affected service")
	addRightCmd.Flags().StringP("subject", "e", "", "affected subject (as regexp)")
}
