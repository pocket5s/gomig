package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	_ "modernc.org/sqlite"
	//_ "github.com/lib/pq"
)

var DB *sql.DB

func connect() error {
	if DB != nil {
		return nil
	}
	log.Println("DB_TYPE:", viper.GetString("GOMIG_DB_TYPE"))
	db, err := sql.Open(viper.GetString("GOMIG_DB_TYPE"), viper.GetString("GOMIG_CONN_STR"))
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
		_, txErr := tx.Exec(s)
		if txErr != nil {
			tx.Rollback()
			log.Printf("sql execution err: %v", txErr)
			return fmt.Errorf("could not execute %s", s)
		}
	}

	if up {
		migInsert := fmt.Sprintf("INSERT INTO migrations (name, ran_at) VALUES ('%s', %d);", fileName, time.Now().UnixNano())
		tx.Exec(migInsert)
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
