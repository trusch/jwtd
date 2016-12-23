package cmd

import (
	"fmt"
	"log"
	"strings"
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
		project, _ := cmd.Flags().GetString("project")
		keyFile, _ := cmd.Flags().GetString("key")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		service, _ := cmd.Flags().GetString("service")
		labels := parseLabels(cmd)
		force, _ := cmd.Flags().GetBool("force")

		if force {
			if keyFile == "" || username == "" || service == "" || len(labels) == 0 {
				log.Fatal("specify --key --username --service and --labels")
			}
		} else {
			if keyFile == "" || username == "" || password == "" || service == "" || len(labels) == 0 {
				log.Fatal("specify --key --username --password --service and --labels")
			}
		}

		if !force {
			user, err := database.GetUser(project, username)
			if err != nil {
				log.Fatalf("failed request: no such user (%v)", username)
			}
			if ok, e := user.CheckPassword(password); e != nil || !ok {
				log.Fatalf("failed request: wrong password (user: %v)", username)
			}
			if ok, e := user.CheckRights(database, project, service, labels); e != nil || !ok {
				log.Fatalf("failed request: no rights (user: %v service: %v, labels: %v)", username, service, labels)
			}
		}

		claims := jwt.Claims{
			"user":    username,
			"service": service,
			"project": project,
			"labels":  labels,
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

func parseLabels(cmd *cobra.Command) map[string]string {
	labelsSlice, _ := cmd.Flags().GetStringSlice("labels")
	labels := make(map[string]string)
	for _, labelStr := range labelsSlice {
		parts := strings.Split(labelStr, "=")
		if len(parts) == 2 {
			labels[parts[0]] = parts[1]
		}
	}
	return labels
}

func init() {
	tokenCmd.AddCommand(createTokenCmd)
	createTokenCmd.Flags().StringP("key", "k", "", "private rsa key to sign the token")
	createTokenCmd.Flags().StringP("username", "u", "", "username to use")
	createTokenCmd.Flags().String("password", "", "password to use")
	createTokenCmd.Flags().BoolP("force", "f", false, "force token creation (no auth checks)")
	createTokenCmd.Flags().String("service", "", "service to use")
	createTokenCmd.Flags().StringSlice("labels", []string{}, "comma separated list of labels (foo=abc,bar=baz)")
}
