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
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Runs the UNDO of the most recent migration.",
	Long:  `Runs the UNDO of the most recent migration against the desired database.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := connect()
		if err != nil {
			log.Fatal(err)
		}
		var downCount int
		if len(args) == 0 {
			downCount = 1
		} else {
			i, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatal("arg needs to be a number")
				return
			}
			downCount = i
		}

		rows, err := DB.Query("SELECT name FROM migrations ORDER BY ran_at DESC LIMIT ?", downCount)
		if err != nil {
			log.Fatal("could not query for migrations to downgrade: %v", err)
			return
		}

		fileNames := make([]string, 0)
		for rows.Next() {
			var fileName string
			err := rows.Scan(&fileName)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				log.Fatal("could not get migration file name from db: %v", err)
				return
			} else {
				fileNames = append(fileNames, fileName)
			}
		}

		migrationDir := "./migrations"
		//sort.Sort(sort.Reverse(sort.StringSlice(fileNames)))

		for _, fileName := range fileNames {
			mf := parseFile(migrationDir + "/" + fileName)
			err := executeTransaction(mf.undoSql, false, fileName)
			if err != nil {
				log.Fatal("could not run down migration on file ", fileName)
				return
			}
		}
		log.Println("migrations complete")
	},
}

func init() {
	rootCmd.AddCommand(downCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
