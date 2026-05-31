package main

import (
	"log"

	"github.com/ivankin-k/go-final-project/pkg/db"
	"github.com/ivankin-k/go-final-project/pkg/server"
)

func main() {
	defer db.DB.Close()

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
