package models

import (
	"time"
	"github.com/golang-sql/sqlexp"
	"context"
	"github.com/google/uuid"
	"log"
)

// PasswordResetCode a password reset code
type PasswordResetCode struct {
	// Code a randomly generated reset code
	Code       string    `json:"code"`
	// Used indicates whether a password reset code has been used. Used codes are necessarily unusable.
	Used       bool      `json:"used"`
	// Usable indicates whether this password reset code can be used. When a user uses a reset code,
	// all previously issued codes are rendered unusable.
	Usable     bool      `json:"usable"`
	// Expiration when the reset code expires
	Expiration time.Time `json:"expiration"`
	// PersonID the person to which this code belongs
	PersonID   int64     `json:"person_id"`
}

// NewPasswordResetCode generates a fresh, unused code for the given account, with the given expiration time.
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

// CreatePasswordResetCode inserts a reset code record into the database.
func (c *PasswordResetCode) CreatePasswordResetCode(q sqlexp.Querier) error {
	_, err := q.ExecContext(
		context.TODO(),
		"INSERT INTO password_reset_code(code, expiration_timestamp, usable, used, person_id) VALUES($1, $2, $3, $4, $5)",
		c.Code, c.Expiration, c.Usable, c.Used, c.PersonID)

	return err
}

// HasCode returns true if a usable password reset code
func (c *PasswordResetCode) HasCode(q sqlexp.Querier) (bool, error) {
	var count int
	row := q.QueryRowContext(
		context.TODO(),
		"SELECT COUNT(*) AS count " +
			"FROM password_reset_code " +
			"WHERE code = $1 " +
			"AND usable = TRUE " +
			"AND used = FALSE " +
			"AND expiration_timestamp > $2", c.Code, time.Now())
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 1, nil
}

// MarkAsUsed sets a particular reset code as unused.
func (c *PasswordResetCode) MarkAsUsed(q sqlexp.Querier) error {
	_, err := q.ExecContext(
		context.TODO(),
		"UPDATE password_reset_code SET used=TRUE, usable=FALSE WHERE code=$1",
		c.Code)
	return err
}

// MarkAllAsUnusable renders all of a user's extant reset codes as unusable.
// This function is invoked when a user successfully uses one of their reset codes.
func (c *PasswordResetCode) MarkAllAsUnusable(q sqlexp.Querier) error {
	_, err := q.ExecContext(
		context.TODO(),
		"UPDATE password_reset_code SET usable=FALSE WHERE person_id=$1",
		c.PersonID)
	return err
}
