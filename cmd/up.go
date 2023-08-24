/*
Copyright © 2023 Robert McIntosh pocket5s@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
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
}
