package db

import (
	"database/sql"
	"errors"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var (
		err   error
		res   sql.Result
		id    int64
		query string
	)

	query = `INSERT INTO ` + dbName + `
			  (date, title, comment, repeat)
			  VALUES
			  (:date, :title, :comment, :repeat);`

	if res, err = DB.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	); err != nil {
		return 0, err
	}
	id, err = res.LastInsertId()

	return id, err
}

func GetTasks(limit int, date, search string) ([]*Task, error) {
	var (
		err   error
		rows  *sql.Rows
		query string
	)

	query = "SELECT id, date, title, comment, repeat FROM " + dbName
	if len(date) > 0 {
		query += " WHERE date = :date"
	} else if len(search) > 0 {
		search = "%" + search + "%"
		query = query + " WHERE title LIKE :search OR comment LIKE :search"
	}
	query = query + " ORDER BY date ASC LIMIT :limit"

	if rows, err = DB.Query(query,
		sql.Named("date", date),
		sql.Named("search", search),
		sql.Named("limit", limit),
	); err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*Task{}

	for rows.Next() {
		var task Task
		task = Task{}

		if err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func GetTask(id int64) (*Task, error) {
	var (
		err   error
		query string
		task  Task
	)

	query = "SELECT id, date, title, comment, repeat FROM " + dbName + " WHERE id = :id"

	task = Task{}
	if err = DB.QueryRow(query, sql.Named("id", id)).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &TaskNotFoundError{id: id}
		}
		return nil, err
	}

	return &task, nil
}

func UpdateTask(id int64, task *Task) error {
	var (
		err   error
		query string
		res   sql.Result
	)

	query = `UPDATE ` + dbName + `
			 SET date=:date, title=:title, comment=:comment, repeat=:repeat
			 WHERE id=:id`

	if res, err = DB.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", id),
	); err != nil {
		return err
	}

	var count int64
	if count, err = res.RowsAffected(); err != nil {
		return err
	}

	if count == 0 {
		return &TaskNotFoundError{id: id}
	}

	return nil
}

func UpdateDate(id int64, date string) error {
	var (
		err   error
		query string
		res   sql.Result
	)
	query = `UPDATE ` + dbName + `
			 SET date=:date
			 WHERE id=:id`

	if res, err = DB.Exec(query, sql.Named("date", date), sql.Named("id", id)); err != nil {
		return err
	}
	var count int64
	if count, err = res.RowsAffected(); err != nil {
		return err
	}
	if count == 0 {
		return &TaskNotFoundError{id: id}
	}
	return nil
}

func DeleteTask(id int64) error {
	var (
		err   error
		query string
		res   sql.Result
	)

	query = "DELETE FROM " + dbName + " WHERE id=:id"

	if res, err = DB.Exec(query, sql.Named("id", id)); err != nil {
		return err
	}
	var count int64
	if count, err = res.RowsAffected(); err != nil {
		return err
	}
	if count == 0 {
		return &TaskNotFoundError{id: id}
	}

	return nil
}
