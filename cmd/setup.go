/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var migrationTable = `
    CREATE TABLE IF NOT EXISTS migrations (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name varchar(255),
        ran_at int
    );
`

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup migration capabilities.",
	Long:  `Creates the migration tables in the target database so migrations can be tracked.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := connect()
		if err != nil {
			log.Fatal("could not connect to database", err)
			return
		}
		_, err = DB.Exec(migrationTable)
		if err != nil {
			log.Fatal("could not create migration table", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
