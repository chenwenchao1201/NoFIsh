package repository

import (
	"database/sql"
	"errors"
	"time"
)

// SQLiteRepository the type for a repository that connects to sqlite database
type SQLiteRepository struct {
	Conn *sql.DB
}

// NewSQLiteRepository returns a new repository with a connection to sqlite
func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		Conn: db,
	}
}

// Migrate creates the table(s) we need
func (repo *SQLiteRepository) Migrate() error {
	err := createPrize(repo)
	if err != nil {
		return err
	}

	err = createSummary(repo)
	if err != nil {
		return err
	}

	return createTask(repo)
}

func createPrize(repo *SQLiteRepository) error {
	query := `
	create table if not exists prizes(
		id integer primary key autoincrement,
		description text not null,
		points int not null,
		is_repeat int not null
		);
	`
	_, err := repo.Conn.Exec(query)
	return err
}

func createTask(repo *SQLiteRepository) error {
	query := `
	create table if not exists tasks(
		id integer primary key autoincrement,
		name varchar(20) not null,
		description text not null,
		due_date int not null,
		completed int not null,
		points int not null,
		is_long_term int not null,
		priority int not null
		);
	`
	_, err := repo.Conn.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func createSummary(repo *SQLiteRepository) error {
	query := `
	create table if not exists summary(
		id integer primary key autoincrement,
		fish_count integer not null,
		finish_count integer not null,
		prize_count integer not null,
		day varchar(10) not null
		);
	`
	_, err := repo.Conn.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// task 相关方法实现
func (repo *SQLiteRepository) InsertTask(t Task) (*Task, error) {
	stmt := "insert into tasks (name, description, due_date, completed, points, is_long_term, priority) values (?, ?, ?, ?, ?, ?, ?)"
	res, err := repo.Conn.Exec(stmt, t.Name, t.Description, t.DueDate.Unix(), t.Completed, t.Points, t.IsLongTerm, t.Priority)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	t.ID = id

	return &t, nil
}
func (repo *SQLiteRepository) AllTasks() ([]Task, error) {
	query := "select id, name, description, due_date, completed, points, is_long_term, priority from tasks order by due_date"
	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Task
	for rows.Next() {
		var t Task
		var unixTime int64
		err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&unixTime,
			&t.Completed,
			&t.Points,
			&t.IsLongTerm,
			&t.Priority,
		)
		if err != nil {
			return nil, err
		}
		t.DueDate = time.Unix(unixTime, 0)
		all = append(all, t)
	}

	return all, nil
}
func (repo *SQLiteRepository) GetTaskByID(id int) (*Task, error) {
	row := repo.Conn.QueryRow("select id, name, description, due_date, completed, points, is_long_term, priority from tasks where id = ?", id)

	var t Task
	var unixTime int64
	err := row.Scan(
		&t.ID,
		&t.Name,
		&t.Description,
		&unixTime,
		&t.Completed,
		&t.Points,
		&t.IsLongTerm,
		&t.Priority,
	)

	if err != nil {
		return nil, err
	}
	t.DueDate = time.Unix(unixTime, 0)

	return &t, nil
}
func (repo *SQLiteRepository) UpdateTask(id int64, updated Task) error {
	if id == 0 {
		return errors.New("id cannot be 0")
	}

	stmt := "update tasks set name = ?, description = ?, due_date = ?, completed = ?, points = ?, is_long_term = ?, priority = ? where id = ?"
	res, err := repo.Conn.Exec(stmt, updated.Name, updated.Description, updated.DueDate.Unix(), updated.Completed, updated.Points, updated.IsLongTerm, updated.Priority, id)
	return updateCheck(err, res)
}

func (repo *SQLiteRepository) DeleteTask(id int64) error {
	res, err := repo.Conn.Exec("delete from tasks where id = ?", id)
	return deleteCheck(err, res)
}

// prize 相关方法实现
func (repo *SQLiteRepository) InsertPrize(p Prize) (*Prize, error) {
	stmt := "insert into prizes (description, points, is_repeat) values (?, ?, ?)"
	res, err := repo.Conn.Exec(stmt, p.Description, p.Points, p.IsRepeat)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	p.ID = id

	return &p, nil
}

func (repo *SQLiteRepository) AllPrizes() ([]Prize, error) {
	query := "select id, description, points, is_repeat from prizes"
	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Prize
	for rows.Next() {
		var p Prize
		err := rows.Scan(
			&p.ID,
			&p.Description,
			&p.Points,
			&p.IsRepeat,
		)
		if err != nil {
			return nil, err
		}
		all = append(all, p)
	}

	return all, nil
}

func (repo *SQLiteRepository) GetPrizeByID(id int) (*Prize, error) {
	row := repo.Conn.QueryRow("select id, description, points, is_repeat from prizes where id = ?", id)

	var p Prize
	err := row.Scan(
		&p.ID,
		&p.Description,
		&p.Points,
		&p.IsRepeat,
	)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (repo *SQLiteRepository) UpdatePrize(id int64, updated Prize) error {
	if id == 0 {
		return errors.New("id cannot be 0")
	}

	stmt := "update prizes set description = ?, points = ?, is_repeat = ? where id = ?"
	res, err := repo.Conn.Exec(stmt, updated.Description, updated.Points, updated.IsRepeat, id)

	return updateCheck(err, res)
}

func (repo *SQLiteRepository) DeletePrize(id int64) error {
	res, err := repo.Conn.Exec("delete from prizes where id = ?", id)
	return deleteCheck(err, res)
}

func deleteCheck(err error, res sql.Result) error {
	if err != nil {
		return err
	}
	rawsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rawsAffected == 0 {
		return errUpdateFailed
	}

	return nil
}

func updateCheck(err error, res sql.Result) error {
	if err != nil {
		return err
	}

	rawsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rawsAffected == 0 {
		return errUpdateFailed
	}

	return nil
}
