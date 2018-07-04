package models

import (
	"github.com/golang-sql/sqlexp"
	"context"
)

type Account struct {
	ID       int64
	Username string `json:"username"`
}

func (p *Account) GetAccountByUsername(q sqlexp.Querier) error {
	return q.QueryRowContext(
		context.TODO(),
		"SELECT id, username "+
			"FROM account WHERE username=$1 "+
			"AND deleted_timestamp IS NULL", p.Username).Scan(&p.ID, &p.Username)
}
