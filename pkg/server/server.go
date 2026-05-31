package server

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ivankin-k/go-final-project/pkg/api"
)

const (
	host = "0.0.0.0"
	port = "7540"

	webDir = "web"
)

func todoPort() string {
	portEnv := os.Getenv("TODO_PORT")
	if len(portEnv) > 0 {
		if _, err := strconv.ParseInt(portEnv, 10, 32); err == nil {
			log.Printf("Using custom port: %s", portEnv)
			return portEnv
		}
	}
	log.Printf("Using default port: %s", port)
	return port
}

func Run() error {
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	api.Init()

	port := todoPort()
	log.Printf("Listening at http://%s:%s\n", host, port)
	if err := http.ListenAndServe(host+":"+port, nil); err != nil {
		return err
	}

	return nil
}
