package cmd

import (
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/spf13/cobra"
)

// updateUserCmd represents the updateUser command
var updateUserCmd = &cobra.Command{
	Use:   "update",
	Short: "update a user",
	Long:  `This updates a user in your database.`,
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

		user, err := database.GetUser(name)
		if err != nil {
			log.Fatal(err)
		}
		if password != "" {
			hash, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if e != nil {
				log.Fatal(e)
			}
			user.PasswordHash = string(hash)
		}
		if len(groups) > 0 {
			user.Groups = groups
		}
		err = database.UpdateUser(user)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	userCmd.AddCommand(updateUserCmd)
	updateUserCmd.Flags().String("password", "", "user password")
	updateUserCmd.Flags().StringSliceP("groups", "g", []string{"default"}, "comma separated list of groups")
}
