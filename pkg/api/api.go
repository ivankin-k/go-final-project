package api

import (
	"log"
	"net/http"
)

const dateFormat = "20060102"

func Init() {
	http.HandleFunc("GET /api/nextdate", nextDateHandler)
	http.HandleFunc("POST /api/task", auth(addTaskHandler))
	http.HandleFunc("GET /api/task", auth(getTaskHandler))
	http.HandleFunc("GET /api/tasks", auth(getTasksHandler))
	http.HandleFunc("PUT /api/task", auth(updateTaskHandler))
	http.HandleFunc("POST /api/task/done", auth(markDoneHandler))
	http.HandleFunc("DELETE /api/task", auth(deleteTaskHandler))

	// Auth
	if len(authPassword) > 0 {
		log.Println("Authentication enabled")
		http.HandleFunc("POST /api/signin", signInHandler)
	}
}
