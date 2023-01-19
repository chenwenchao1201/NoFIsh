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
	query, err, err2 := createHoldings(repo)
	if err2 != nil {
		return err2
	}

	query, err, err3 := createSummary(query, err, repo)
	if err3 != nil {
		return err3
	}

	query, err, err4 := createTask(query, err, repo)
	if err4 != nil {
		return err4
	}

	return createPrize(query, err, repo)
}

func createPrize(query string, err error, repo *SQLiteRepository) error {
	query = `
	create table if not exists prizes(
		id integer primary key autoincrement,
		description text not null,
		points int not null,
		is_repeat int not null,
		exchange_time integer not null);
	`
	_, err = repo.Conn.Exec(query)
	return err
}

func createTask(query string, err error, repo *SQLiteRepository) (string, error, error) {
	query = `
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
	_, err = repo.Conn.Exec(query)
	if err != nil {
		return "", nil, err
	}
	return query, err, nil
}

func createSummary(query string, err error, repo *SQLiteRepository) (string, error, error) {
	query = `
	create table if not exists summary(
		id integer primary key autoincrement,
		fish_count integer not null,
		finish_count integer not null,
		prize_count integer not null,
		day varchar(10) not null
		);
	`
	_, err = repo.Conn.Exec(query)
	if err != nil {
		return "", nil, err
	}
	return query, err, nil
}

func createHoldings(repo *SQLiteRepository) (string, error, error) {
	query := `
	create table if not exists holdings(
		id integer primary key autoincrement,
		amount real not null,
		purchase_date integer not null,
		purchase_price integer not null);
	`
	_, err := repo.Conn.Exec(query)
	if err != nil {
		return "", nil, err
	}
	return query, err, nil
}

func (repo *SQLiteRepository) InsertHolding(holdings Holdings) (*Holdings, error) {
	stmt := "insert into holdings (amount, purchase_date, purchase_price) values (?, ?, ?)"
	res, err := repo.Conn.Exec(stmt, holdings.Amount, holdings.PurchaseDate.Unix(), holdings.PurchasePrice)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	holdings.ID = id

	return &holdings, nil

}

// AllHoldings returns all holdings, by purchase date
func (repo *SQLiteRepository) AllHoldings() ([]Holdings, error) {
	query := "select id, amount, purchase_date, purchase_price from holdings order by purchase_date"
	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Holdings
	for rows.Next() {
		var h Holdings
		var unixTime int64
		err := rows.Scan(
			&h.ID,
			&h.Amount,
			&unixTime,
			&h.PurchasePrice,
		)
		if err != nil {
			return nil, err
		}
		h.PurchaseDate = time.Unix(unixTime, 0)
		all = append(all, h)
	}

	return all, nil
}

func (repo *SQLiteRepository) GetHoldingByID(id int) (*Holdings, error) {
	row := repo.Conn.QueryRow("select id, amount, purchase_date, purchase_price from holdings where id = ?", id)

	var h Holdings
	var unixTime int64
	err := row.Scan(
		&h.ID,
		&h.Amount,
		&unixTime,
		&h.PurchasePrice,
	)

	if err != nil {
		return nil, err
	}
	h.PurchaseDate = time.Unix(unixTime, 0)

	return &h, nil

}
func (repo *SQLiteRepository) UpdateHolding(id int64, update Holdings) error {
	if id == 0 {
		return errors.New("id cannot be 0")
	}

	stmt := "update holdings set amount = ?, purchase_date = ?, purchase_price = ? where id = ?"
	res, err := repo.Conn.Exec(stmt, update.Amount, update.PurchaseDate.Unix(), update.PurchasePrice, id)

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

func (repo *SQLiteRepository) DeleteHolding(id int64) error {
	res, err := repo.Conn.Exec("delete from holdings where id = ?", id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errDeleteFailed
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
func (repo *SQLiteRepository) DeleteTask(id int64) error {
	res, err := repo.Conn.Exec("delete from tasks where id = ?", id)
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
