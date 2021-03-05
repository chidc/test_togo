package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"togo/internal/storages"
)

// LiteDB for working with sqllite
type PostgresDB struct {
	DB *sql.DB
}

// RetrieveTasks returns tasks if match userID AND createDate.
func (l *PostgresDB) RetrieveTasks(ctx context.Context, userID, createdDate sql.NullString) ([]*storages.Task, error) {
	stmt := fmt.Sprintf(`SELECT id, content, user_id, created_date FROM tasks WHERE user_id = '%s' AND created_date = '%s'`, userID.String, createdDate.String)
	rows, err := l.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*storages.Task
	for rows.Next() {
		t := &storages.Task{}
		err := rows.Scan(&t.ID, &t.Content, &t.UserID, &t.CreatedDate)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// AddTask adds a new task to DB
func (l *PostgresDB) AddTask(ctx context.Context, t *storages.Task) (string, error) {
	lmt := fmt.Sprintf(`SELECT count(id),(select max_todo from users where id = '%s') from tasks where user_id= '%s' and created_date= '%s'`, t.UserID, t.UserID, t.CreatedDate)
	//fmt.Println("lmt:", lmt)
	rows, err := l.DB.QueryContext(ctx, lmt)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var count, max int
		err := rows.Scan(&count, &max)
		fmt.Println("count: ", count, "max: ", max)

		if err != nil {
			return "", err
		}
		if count == max {
			return "gioi han", nil
		}
	}
	stmt := fmt.Sprintf(`INSERT INTO tasks (id, content, user_id, created_date) VALUES ('%s', '%s', '%s', '%s')`, t.ID, t.Content, t.UserID, t.CreatedDate)
	//fmt.Println("Stmt:", stmt)
	_, err = l.DB.ExecContext(ctx, stmt)
	if err != nil {
		return "", err
	}

	return "", nil
}

// ValidateUser returns tasks if match userID AND password
func (l *PostgresDB) ValidateUser(ctx context.Context, userID, pwd sql.NullString) bool {
	stmt := fmt.Sprintf(`SELECT id FROM users WHERE id = '%s' AND password = '%s'`, userID.String, pwd.String)
	row := l.DB.QueryRowContext(ctx, stmt)
	u := &storages.User{}

	err := row.Scan(&u.ID)
	if err != nil {
		return false
	}

	return true
}
