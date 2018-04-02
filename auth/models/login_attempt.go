package models

import (
	"time"
	"database/sql"
)

type LoginAttempt struct {
	Id       int64
	Created  time.Time
	Success  bool
	PersonID int64
}

func (l *LoginAttempt) CreateLoginAttempt(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO login_attempt(id, created_timestamp, success, person_id) VALUES($1, $2, $3, $4)",
		l.Id, l.Created, l.Success, l.PersonID)

	return err
}