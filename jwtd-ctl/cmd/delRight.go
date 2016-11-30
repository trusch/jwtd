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
		subject, _ := cmd.Flags().GetString("subject")
		if service == "" || subject == "" {
			log.Fatal("specifiy --service and --subject")
		}
		group, err := database.GetGroup(name)
		if err != nil {
			log.Fatal(err)
		}
		foundIdx := -1
		for idx, right := range group.Rights {
			if right.Service == service && right.Subject == subject {
				foundIdx = idx
				break
			}
		}
		if foundIdx == -1 {
			log.Fatal("no such right")
		}
		group.Rights = append(group.Rights[:foundIdx], group.Rights[foundIdx+1:]...)
		err = database.UpdateGroup(group)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	groupCmd.AddCommand(delRightCmd)
	delRightCmd.Flags().StringP("subject", "e", "", "affected subject (as regexp)")
	delRightCmd.Flags().StringP("service", "s", "", "affected service")
}
