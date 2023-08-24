/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
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
