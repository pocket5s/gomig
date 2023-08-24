/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Applies any pending migrations.",
	Long:  `Applies any pending migrations to the desired database`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			err := connect()
			if err != nil {
				log.Fatal(err)
			}
			var mostRecentFile string
			result := DB.QueryRow("SELECT name FROM migrations ORDER BY ran_at DESC LIMIT 1")
			if result != nil {
				err := result.Scan(&mostRecentFile)
				if err != nil && !errors.Is(err, sql.ErrNoRows) {
					log.Fatal("could not query for latest migration file name from db: %v", err)
					return
				}
			}

			migrationDir := "./migrations"
			files, err := ioutil.ReadDir(migrationDir)
			if err != nil {
				log.Fatal(err)
				return
			}

			if mostRecentFile == "" {
				mf := parseFile(migrationDir + "/" + files[0].Name())
				err = executeTransaction(mf.sqlToRun, true, files[0].Name())
				if err != nil {
					log.Fatal(err)
					return
				}
			} else {
				var found bool
				for _, file := range files {
					if file.Name() == mostRecentFile {
						found = true
					}

					if found {
						mf := parseFile(migrationDir + "/" + file.Name())
						err = executeTransaction(mf.sqlToRun, true, file.Name())
						if err != nil {
							log.Fatal(err)
							return
						}
					}
				}
			}

			log.Println("migrations complete")
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
