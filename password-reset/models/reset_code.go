package models

import (
	"time"
	"github.com/golang-sql/sqlexp"
	"context"
	"github.com/google/uuid"
	"log"
)

type PasswordResetCode struct {
	Code       string    `json:"code"`
	Used       bool      `json:"used"`
	Usable     bool      `json:"usable"`
	Expiration time.Time `json:"expiration"`
	PersonID   int64     `json:"person_id"`
}

func NewPasswordResetCode(personID int64, expiration time.Time) (*PasswordResetCode, error) {
	var id string
	if u, err := uuid.NewRandom(); err != nil {
		return nil, err
	} else {
		id = u.String()
		log.Printf("Generated password reset code: %s", id)
	}

	return &PasswordResetCode{
		Code:       id,
		Used:       false,
		Usable:     true,
		Expiration: expiration,
		PersonID:   personID}, nil
}

func (c *PasswordResetCode) CreatePasswordResetCode(q sqlexp.Querier) error {
	_, err := q.ExecContext(
		context.TODO(),
		"INSERT INTO password_reset_code(code, expiration_timestamp, usable, used, person_id) VALUES($1, $2, $3, $4, $5)",
		c.Code, c.Expiration, c.Usable, c.Used, c.PersonID)

	return err
}

func (c *PasswordResetCode) HasCode(q sqlexp.Querier) (bool, error) {
	var count int
	row := q.QueryRowContext(
		context.TODO(),
		"SELECT COUNT(*) AS count " +
			"FROM password_reset_code " +
			"WHERE code = $1 " +
			"AND usable = $2 " +
			"AND used = $3 " +
			"AND expiration_timestamp > $4", c.Code, true, false, time.Now())
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 1, nil
}

func (c *PasswordResetCode) MarkAsUsed(q sqlexp.Querier) error {
	_, err := q.ExecContext(
		context.TODO(),
		"UPDATE password_reset_code SET used=$1, usable=$2 WHERE code=$3",
		true, false, c.Code)
	return err
}

func (c *PasswordResetCode) MarkAllAsUnusable(q sqlexp.Querier) error {
	_, err := q.ExecContext(
		context.TODO(),
		"UPDATE password_reset_code SET usable=$1 WHERE person_id=$2",
		false, c.PersonID)
	return err
}
