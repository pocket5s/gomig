/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Applies any pending migrations.",
	Long:  `Applies any pending migrations to the desired database`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			//TODO: read db to see latest migration

			//TODO: scan migrations dir for new files

			//TODO: load migration sql and parse
			fileName := "./migrations/20230822223837_a_new_migration.sql"
			mf := parseFile(fileName)

			err := executeTransaction(mf.sqlToRun)
			if err != nil {
				os.Exit(1)
			}
			log.Printf("Migration %s run", fileName)
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
