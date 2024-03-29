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
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/libsql/libsql-client-go/libsql"

	//_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"

	"github.com/spf13/viper"
)

var DB *sql.DB

func connect() error {
	if DB != nil {
		return nil
	}
	//db, err := sql.Open(viper.GetString("GOMIG_DB_TYPE"), viper.GetString("GOMIG_CONN_STR"))
	db, err := sql.Open("libsql", viper.GetString("GOMIG_CONN_STR"))
	if err == nil {
		DB = db
		return DB.Ping()
	}
	return err
}

func executeTransaction(statements []string, up bool, fileName string) error {
	if DB == nil {
		return errors.New("No database connected")
	}
	tx, err := DB.Begin()
	if err != nil {
		log.Printf("err: %v", err)
		return fmt.Errorf("could not start db transaction")
	}
	for _, s := range statements {
		if s != "" {
			log.Printf("executing: %s", s)
			_, txErr := tx.Exec(s)
			if txErr != nil {
				tx.Rollback()
				log.Printf("sql execution err: %v", txErr)
				return fmt.Errorf("could not execute %s", s)
			}
		}
	}

	if up {
		tx.Exec("INSERT INTO migrations (name, ran_at) VALUES (?, ?);", fileName, time.Now().UnixMilli())
	} else {
		tx.Exec("DELETE FROM migrations WHERE name = ?", fileName)
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("commit err: %v", err)
		return fmt.Errorf("could not commit transaction")
	}
	return nil
}
