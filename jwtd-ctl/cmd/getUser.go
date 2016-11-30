package cmd

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

// getUserCmd represents the getUser command
var getUserCmd = &cobra.Command{
	Use:   "get",
	Short: "get a user",
	Long:  `This prints the database content for a given user.`,
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
		user, err := database.GetUser(name)
		if err != nil {
			log.Fatal(err)
		}
		bs, err := yaml.Marshal(user)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(bs))
	},
}

func init() {
	userCmd.AddCommand(getUserCmd)

}
