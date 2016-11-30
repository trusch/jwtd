package cmd

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
	"github.com/trusch/jwtd/jwt"
)

// verifyTokenCmd represents the verifyToken command
var verifyTokenCmd = &cobra.Command{
	Use:   "verify",
	Short: "verify a token",
	Long:  `This verifies a given token and prints its claims.`,
	Run: func(cmd *cobra.Command, args []string) {
		token, _ := cmd.Flags().GetString("token")
		keyFile, _ := cmd.Flags().GetString("key")
		if token == "" || keyFile == "" {
			log.Fatal("specify --token and --key")
		}
		key, err := jwt.LoadPublicKey(keyFile)
		if err != nil {
			log.Fatal(err)
		}
		claims, err := jwt.ValidateToken(token, key)
		if err != nil {
			log.Fatal(err)
		}
		bs, err := yaml.Marshal(claims)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(bs))
	},
}

func init() {
	tokenCmd.AddCommand(verifyTokenCmd)
	verifyTokenCmd.Flags().StringP("token", "t", "", "token to verify")
	verifyTokenCmd.Flags().StringP("key", "k", "", "public key to use")
}
