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
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var tpl = `-- // %s
-- Migration SQL here


-- //@UNDO
-- SQL to undo the migration here`

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new migration file.",
	Long:  `Creates a new migration file using the supplied name.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Name of migration required")
		} else {
			return nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		n := time.Now().In(time.UTC)
		fn := strings.ReplaceAll(args[0], " ", "_")
		filename := "%d%02d%02d%02d%02d%02d_" + fn + ".sql"
		filename = fmt.Sprintf(filename, n.Year(), n.Month(), n.Day(), n.Hour(), n.Minute(), n.Second())
		output := fmt.Sprintf(tpl, args[0])
		err := os.WriteFile("./migrations/"+filename, []byte(output), 0644)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("migration file created")
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringP("dir", "d", "./migrations", "Directory to output migration")
}
