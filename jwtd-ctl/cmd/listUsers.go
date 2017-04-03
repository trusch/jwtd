package cmd

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

// listUsersCmd represents the listUsers command
var listUsersCmd = &cobra.Command{
	Use:   "list",
	Short: "list users",
	Long:  `This lists all users and dumps it as yaml.`,
	Run: func(cmd *cobra.Command, args []string) {
		database := getDB()
		users, err := database.ListUsers()
		if err != nil {
			log.Fatal(err)
		}
		bs, err := yaml.Marshal(users)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(bs))
	},
}

func init() {
	userCmd.AddCommand(listUsersCmd)
}
