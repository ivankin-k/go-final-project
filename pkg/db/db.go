package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

const (
	dbFile = "scheduler.db"
	dbName = "scheduler"

	schema = `
CREATE TABLE IF NOT EXISTS ` + dbName + ` (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "date" CHAR(8) NOT NULL DEFAULT "",
    "title" VARCHAR(128) NOT NULL,
    "comment" TEXT(250) NOT NULL DEFAULT '',
    "repeat" VARCHAR(128) NOT NULL
);

CREATE index tasks_date ON ` + dbName + ` ("date");
`
)

var DB *sql.DB

func todoDBFile() string {
	dbEnv := os.Getenv("TODO_DBFILE")
	if len(dbEnv) > 0 {
		return dbEnv
	}
	return dbFile
}

func Connect() error {
	var (
		installDB bool
		err       error
		dbFile    string
	)
	dbFile = todoDBFile()

	_, err = os.Stat(dbFile)
	if errors.Is(err, os.ErrNotExist) {
		installDB = true
	}

	// Open DB file
	if DB, err = sql.Open("sqlite", dbFile); err != nil {
		return fmt.Errorf("Can't open DB (%s):\n%v\n", dbFile, err)
	}

	log.Printf("DB opened (%s)", dbFile)

	// InstallDB if needed
	if installDB {
		if _, err = DB.Exec(schema); err != nil {
			return fmt.Errorf("Can't install DB, exec schema (%s):\n%v\n", dbFile, err)
		}
		log.Print("DB schema installed successfully")
	}

	return nil
}
