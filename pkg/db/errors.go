package db

import "fmt"

var ErrTaskNotFound = &TaskNotFoundError{}

type TaskNotFoundError struct {
	id int64
}

func (err *TaskNotFoundError) Error() string {
	return fmt.Sprintf("No task found with ID=%d", err.id)
}
