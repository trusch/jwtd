package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/trusch/jwtd/db"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "jwtd-ctl",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringP("config", "c", "/etc/jwtd/config.yaml", "config file")
	RootCmd.PersistentFlags().StringP("project", "p", "default", "project to use")
}

func getDB() *db.DB {
	config, _ := RootCmd.Flags().GetString("config")
	db, err := db.New(config)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
