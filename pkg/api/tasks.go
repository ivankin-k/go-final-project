package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/ivankin-k/go-final-project/pkg/db"
)

func verifyID(data string) (int64, error) {
	// Check if ID is received
	if len(data) == 0 {
		return 0, errors.New(`"id" is missing in request`)
	}

	var (
		id  int64
		err error
	)
	// Check if ID is an integer
	if id, err = strconv.ParseInt(data, 10, 64); err != nil {
		return 0, fmt.Errorf(`Invalid "id" received, not a number: %s`, data)
	}

	return id, nil
}

func verifyTitleAndDate(task *db.Task) error {

	// Check if title is received
	if len(task.Title) == 0 {
		return errors.New("\"title\" field is missing for new task request")
	}

	// Set to now, if not date received
	var now = time.Now()
	if len(task.Date) == 0 {
		task.Date = now.Format(dateFormat)
		return nil
	}

	// Parse date received
	var (
		taskDate time.Time
		err      error
	)
	if taskDate, err = time.Parse(dateFormat, task.Date); err != nil {
		return err
	}

	// Get next date if repeat rule is received
	var next string
	if len(task.Repeat) > 0 {
		next, err = nextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	// Use date, now or next as a final date value
	if after(now, taskDate) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(dateFormat)
		} else {
			task.Date = next
		}
	}
	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {

	// Read and json-decode new task request
	var (
		err     error
		newTask *db.Task
	)
	newTask = &db.Task{}
	defer r.Body.Close()
	if err = json.NewDecoder(r.Body).Decode(newTask); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check title and date
	if err = verifyTitleAndDate(newTask); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return error or ID
	var id int64
	if id, err = db.AddTask(newTask); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, &struct {
		ID int64 `json:"id"`
	}{
		ID: id,
	})
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {

	// Read and json-decode request
	var (
		err  error
		task *db.Task
	)
	task = &db.Task{}
	defer r.Body.Close()
	if err = json.NewDecoder(r.Body).Decode(task); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var id int64

	if id, err = verifyID(task.ID); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Make other required checks
	if err = verifyTitleAndDate(task); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update task data
	if err = db.UpdateTask(id, task); err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeError(w, err.Error(), http.StatusNotFound)
		} else {
			writeError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, struct{}{})
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		tasks []*db.Task
	)

	var (
		date,
		search string
	)
	// If search is received
	if len(r.FormValue("search")) != 0 {
		var (
			err error
			t   time.Time
		)

		search = r.FormValue("search")

		// If search is a valid date
		if t, err = time.Parse("02.01.2006", search); err == nil {
			date = t.Format(dateFormat)
		} else
		// If search is NOT a valid date
		{
			search = strings.NewReplacer("%", `\%`, "_", `\_`).Replace(search)
		}
	}

	// Get tasks from DB
	if tasks, err = db.GetTasks(50, date, search); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Потом уже вспомнил, что сортировать надо в SQL-запросе, но удалять стало жалко :)
	slices.SortFunc(tasks, func(a, b *db.Task) int {
		if a.Date < b.Date {
			return -1
		}
		if a.Date > b.Date {
			return 1
		}
		return 0
	})

	// Return tasks
	writeJSON(w, &struct {
		Tasks []*db.Task `json:"tasks"`
	}{
		Tasks: tasks,
	})

}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		id  int64
	)
	// Verify ID
	if id, err = verifyID(r.FormValue("id")); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get task from DB
	var task *db.Task
	if task, err = db.GetTask(id); err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeError(w, err.Error(), http.StatusNotFound)
		} else {
			writeError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Return task
	writeJSON(w, task)
}

func markDoneHandler(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		id  int64
	)
	if id, err = verifyID(r.FormValue("id")); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var task *db.Task
	if task, err = db.GetTask(id); err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeError(w, err.Error(), http.StatusNotFound)
		} else {
			writeError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if task.Repeat == "" {
		if err = db.DeleteTask(id); err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, struct{}{})
		return
	}

	var newDate string
	if newDate, err = nextDate(time.Now(), task.Date, task.Repeat); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = db.UpdateDate(id, newDate); err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeError(w, err.Error(), http.StatusNotFound)
		} else {
			writeError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, struct{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		id  int64
	)
	// Verify ID
	if id, err = verifyID(r.FormValue("id")); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Delete task
	if err = db.DeleteTask(id); err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeError(w, err.Error(), http.StatusNotFound)
		} else {
			writeError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, struct{}{})
}
