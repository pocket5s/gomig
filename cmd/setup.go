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
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var sqliteMigrationTable = `
    CREATE TABLE IF NOT EXISTS migrations (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name VARCHAR(255) NOT NULL,
        ran_at INT NOT NULL
    );
`

var postgresMigrationTable = `
    CREATE TABLE IF NOT EXISTS migrations (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        ran_at INT NOT NULL
    );
`

var mysqlMigrationTable = `
    CREATE TABLE IF NOT EXISTS migrations (
        id MEDIUMINT NOT NULL AUTO_INCREMENT,
        name VARCHAR(255) NOT NULL,
        ran_at BIGINT NOT NULL,
        PRIMARY KEY (id)
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

		migrationTable := ""
		dbType := viper.GetString("GOMIG_DB_TYPE")
		if dbType == "libsql" {
			migrationTable = sqliteMigrationTable
		} else if dbType == "pgx" {
			migrationTable = postgresMigrationTable
		} else if dbType == "mysql" {
			migrationTable = mysqlMigrationTable
		} else {
			log.Fatal("Unknown db type: ", dbType)
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
