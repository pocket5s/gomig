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
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Shows the status of the current migrations.",
	Long:  `Shows the status of the current migrations.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := connect()
		if err != nil {
			log.Fatal(err)
		}

		rows, err := DB.Query("SELECT name, ran_at FROM migrations")
		defer rows.Close()
		if err != nil {
			log.Fatal("could not query for migrations to generate status: %v", err)
			return
		}

		fileNames := make([]string, 0)
		fmt.Println("")
		fmt.Printf("+%s+\n", strings.Repeat("-", 78))
		fmt.Printf("| %-45s | %-28s |\n", "File Name", "Migration Ran At")
		fmt.Printf("+%s+\n", strings.Repeat("-", 78))
		for rows.Next() {
			var name string
			var ranAt int64
			err = rows.Scan(&name, &ranAt)
			if err != nil {
				log.Fatal("could not scan row for migration information")
				return
			}
			fileNames = append(fileNames, name)
			fmt.Printf("| %-45s | %-28s |\n", name, time.UnixMilli(ranAt).Format(time.UnixDate))
		}

		migrationDir := "./migrations"
		files, err := ioutil.ReadDir(migrationDir)
		if err != nil {
			log.Fatal(err)
			return
		}
		for _, f := range files {
			if fileNotRan(f.Name(), fileNames) {
				fmt.Printf("| %-45s | %-28s |\n", f.Name(), "Pending...")
			}
		}
		fmt.Printf("+%s+\n", strings.Repeat("-", 78))
	},
}

// yes brute force, but let's face it, the list won't be _that_ big and array scans are fast enough
func fileNotRan(fileName string, fileList []string) bool {
	for _, name := range fileList {
		if name == fileName {
			return false
		}
	}

	return true
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
