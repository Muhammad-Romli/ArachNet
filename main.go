package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/muhammad-romli/ArachNet/internal/config"
	"github.com/muhammad-romli/ArachNet/internal/database"
)

var maxRetries = 2

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Enter the command to start using the CLI applications")
		fmt.Println("usage: app <command> <param>")
		os.Exit(1)
	}
	commandName := os.Args[1]
	commandParams := os.Args[2:]

	data, err := config.Read()
	if err != nil {
		fmt.Printf("file is not in destination")
	}
	if data == nil {
		fmt.Printf("There is no data inside targeted file")
		os.Exit(1)
	}

	db, err := OpenSQL("postgres", data.DBURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	mainCommands := commands{CommandsMap: make(map[string]func(*state, command) error)}
	mainCommands.register("login", handlerLogin)
	mainCommands.register("register", handlerRegister)
	mainState := state{}
	mainState.cfg = data
	mainState.db = dbQueries

	cmd := command{commandName, commandParams}
	if err = mainCommands.run(&mainState, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func OpenSQL(usn string, dbURL string) (*sql.DB, error) {
	retries := 0
	for retries < maxRetries {
		db, err := sql.Open(usn, dbURL)
		if err == nil {
			return db, err // CHANGED: return the pointer, not *db
		}
		retries++ // CHANGED: was never incremented, so it looped forever on failure
		fmt.Printf("Retrying connecting to databases")
	}
	return nil, fmt.Errorf("can't connect to database")
}
