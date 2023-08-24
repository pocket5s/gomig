package cmd

import (
	"bufio"
	"os"
	"strings"
)

type migrationFile struct {
	sqlToRun []string
	undoSql  []string
}

func parseFile(file string) migrationFile {
	mf := migrationFile{}
	f, err := os.Open(file)

	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	sql := ""
	undo := ""
	var inUndo bool
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "-- //@UNDO") {
			inUndo = true
		}

		if !strings.HasPrefix(line, "--") {
			if inUndo {
				undo += line
			} else {
				sql += line
			}
		}
	}
	mf.sqlToRun = strings.Split(sql, ";")
	mf.undoSql = strings.Split(undo, ";")

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return mf
}
