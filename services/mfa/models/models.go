package models

import (
	"time"
	"github.com/google/uuid"
	"context"
	"github.com/golang-sql/sqlexp"
)

type MFACode struct {
	ID         uuid.UUID `json:"id"`
	Code       string    `json:"code"`
	Created    time.Time `json:"created_at"`
	Expiration time.Time `json:"expiration"`
	Usable     bool      `json:"usable"`
	Used       bool      `json:"used"`
	AccountID   int64     `json:"account_id"`
}

func (c *MFACode) Create(q sqlexp.Querier) error {
	_, err := q.ExecContext(
		context.TODO(),
		"INSERT INTO authentication_code(id, code, created_timestamp, expiration_timestamp, usable, used, account_id) VALUES($1, $2, $3, $4, $5, $6, $7)",
		c.ID, c.Code, c.Created, c.Expiration, c.Usable, c.Used, c.AccountID)

	return err
}