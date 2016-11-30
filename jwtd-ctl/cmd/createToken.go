package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/trusch/jwtd/jwt"
)

// createTokenCmd represents the createToken command
var createTokenCmd = &cobra.Command{
	Use:   "create",
	Short: "create a token",
	Long:  `This creates a token.`,
	Run: func(cmd *cobra.Command, args []string) {
		database := getDB()
		keyFile, _ := cmd.Flags().GetString("key")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		service, _ := cmd.Flags().GetString("service")
		subject, _ := cmd.Flags().GetString("subject")
		force, _ := cmd.Flags().GetBool("force")

		if force {
			if keyFile == "" || username == "" || service == "" || subject == "" {
				log.Fatal("specify --key --username --service and --subject")
			}
		} else {
			if keyFile == "" || username == "" || password == "" || service == "" || subject == "" {
				log.Fatal("specify --key --username --password --service and --subject")
			}
		}

		if !force {
			user, err := database.GetUser(username)
			if err != nil {
				log.Fatalf("failed request: no such user (%v)", username)
			}
			if ok, e := user.CheckPassword(password); e != nil || !ok {
				log.Fatalf("failed request: wrong password (user: %v)", username)
			}
			if ok, e := user.CheckRights(database, service, subject); e != nil || !ok {
				log.Fatalf("failed request: no rights (user: %v service: %v, subject: %v)", username, service, subject)
			}
		}

		claims := jwt.Claims{
			"user":    username,
			"service": service,
			"subject": subject,
			"nbf":     time.Now(),
			"exp":     time.Now().Add(10 * time.Minute),
		}
		key, err := jwt.LoadPrivateKey(keyFile)
		if err != nil {
			log.Fatal(err)
		}
		token, err := jwt.CreateToken(claims, key)
		if err != nil {
			log.Fatal("failed request: can not generate token (wtf?!)")
		}
		fmt.Println(token)
	},
}

func init() {
	tokenCmd.AddCommand(createTokenCmd)
	createTokenCmd.Flags().StringP("key", "k", "", "private rsa key to sign the token")
	createTokenCmd.Flags().StringP("username", "u", "", "username to use")
	createTokenCmd.Flags().StringP("password", "p", "", "password to use")
	createTokenCmd.Flags().BoolP("force", "f", false, "force token creation (no auth checks)")
	createTokenCmd.Flags().String("service", "", "service to use")
	createTokenCmd.Flags().String("subject", "", "subject to use")
}
