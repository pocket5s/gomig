package cmd

import (
	"database/sql"
	"fmt"
	"log"
	//_ "github.com/lib/pq"
)

var DB *sql.DB

func connect(dbType string, connStr string) error {
	if DB != nil {
		return nil
	}
	db, err := sql.Open(dbType, connStr)
	if err == nil {
		DB = db
		return DB.Ping()
	}
	return err
}

func executeTransaction(statements []string) error {
	tx, err := DB.Begin()
	if err != nil {
		log.Printf("err: %v", err)
		return fmt.Errorf("could not start db transaction")
	}
	for _, s := range statements {
		_, txErr := tx.Exec(s)
		if txErr != nil {
			tx.Rollback()
			log.Printf("err: %v", txErr)
			return fmt.Errorf("could not execute %s", s)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("err: %v", err)
		return fmt.Errorf("could not commit transaction")
	}
	return nil
}
