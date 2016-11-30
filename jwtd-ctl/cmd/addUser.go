package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// addUserCmd represents the addUser command
var addUserCmd = &cobra.Command{
	Use:   "add",
	Short: "add a user",
	Long:  `This adds a new user to your jwtd server`,
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
		password, _ := cmd.Flags().GetString("password")
		groups, _ := cmd.Flags().GetStringSlice("groups")
		err := database.CreateUser(name, password, groups)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	userCmd.AddCommand(addUserCmd)
	addUserCmd.Flags().StringP("password", "p", "", "user password")
	addUserCmd.Flags().StringSliceP("groups", "g", []string{"default"}, "comma separated list of groups")
}
