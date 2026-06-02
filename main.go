package main

import (
	"log"

	"github.com/ivankin-k/go-final-project/pkg/db"
	"github.com/ivankin-k/go-final-project/pkg/server"
)

func main() {
	var err error

	defer db.DB.Close()

	if err = db.Connect(); err != nil {
		log.Print(err)
	} else if err := server.Run(); err != nil {
		log.Print(err)
	}
}
