package models

import (
	"github.com/golang-sql/sqlexp"
	"context"
)

type Person struct {
	ID       int64
	Username string `json:"username"`
}

func (p *Person) GetPersonByUsername(q sqlexp.Querier) error {
	return q.QueryRowContext(
		context.TODO(),
		"SELECT id, username "+
			"FROM person WHERE username=$1 "+
			"AND deleted_timestamp IS NULL", p.Username).Scan(&p.ID, &p.Username)
}
