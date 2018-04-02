package models

import (
	"gopkg.in/guregu/null.v3"
	"strconv"
	"time"
	"fmt"
	"github.com/golang-sql/sqlexp"
	"context"
)

type Person struct {
	Id                       int64       `json:"id"`
	IdStr                    string      `json:"id_str"`
	Created                  null.Time   `json:"created_at"`
	Deleted                  null.Time   `json:"deleted_at"`
	LastModified             null.Time   `json:"last_modified_at"`
	Disabled                 bool        `json:"disabled"`
	MultiFactorRequired      bool        `json:"multi_factor_required"`
	Username                 null.String `json:"username"`
	CurrentPasswordID        null.Int    `json:"current_password_id"`
	CurrentPasswordIdStr     string      `json:"current_password_id_str"`
	PrimaryEmailAddressID    null.Int    `json:"primary_email_address_id"`
	PrimaryEmailAddressIdStr string      `json:"primary_email_address_id_str"`
	EmailAddress             string      `json:"email_address"`
}

func (p *Person) GetPersonByUsername(q sqlexp.Querier) error {
	return q.QueryRowContext(
		context.TODO(),
		"SELECT p.id, p.created_timestamp, p.deleted_timestamp, p.last_modified_timestamp, p.disabled, "+
			"p.multi_factor_required, p.username, p.current_password_id, p.primary_email_address_id "+
			"FROM Person p "+
			"WHERE p.username=$1 " +
			"AND p.deleted_timestamp IS NULL", p.Username).Scan(&p.Id, &p.Created, &p.Deleted, &p.LastModified, &p.Disabled,
				&p.MultiFactorRequired, &p.Username, &p.CurrentPasswordID, &p.PrimaryEmailAddressID)
}

func (p *Person) GetPersonByEmailAddress(q sqlexp.Querier) error {
	return q.QueryRowContext(
		context.TODO(),
		"SELECT p.id, p.created_timestamp, p.deleted_timestamp, p.last_modified_timestamp, p.disabled, "+
			"p.multi_factor_required, p.username, p.current_password_id, p.primary_email_address_id, " +
			"e.email_address "+
			"FROM email_address e JOIN Person p ON e.personId = p.id "+
			"WHERE e.email_address=$1 " +
			"AND p.deleted_timestamp IS NULL", p.EmailAddress).Scan(&p.Id, &p.Created, &p.Deleted, &p.LastModified,
				&p.Disabled, &p.MultiFactorRequired, &p.Username, &p.CurrentPasswordID,
					&p.PrimaryEmailAddressID, &p.EmailAddress)
}

func (p *Person) GetDeletedPerson(q sqlexp.Querier) error {
	return q.QueryRowContext(
		context.TODO(),
		"SELECT p.id, p.created_timestamp, p.deleted_timestamp, p.last_modified_timestamp, p.disabled, " +
			"p.multi_factor_required, p.username, p.current_password_id, p.primary_email_address_id "+
			"FROM Person p"+
			"WHERE p.id=$1 "+
			"AND p.deleted_timestamp IS NOT NULL", p.Id).Scan(&p.Id, &p.Created, &p.Deleted, &p.LastModified,
				&p.Disabled, &p.MultiFactorRequired, &p.Username, &p.CurrentPasswordID,
					&p.PrimaryEmailAddressID)
}

func (p *Person) GetPerson(q sqlexp.Querier) error {
	return q.QueryRowContext(
		context.TODO(),
		"SELECT p.id, p.created_timestamp, p.deleted_timestamp, p.last_modified_timestamp, p.disabled, " +
			"p.multi_factor_required, p.username, p.current_password_id, p.primary_email_address_id "+
			"FROM Person p "+
			"WHERE p.id=$1 "+
			"AND p.deleted_timestamp IS NULL", p.Id).Scan(&p.Id, &p.Created, &p.Deleted, &p.LastModified, &p.Disabled,
		&p.MultiFactorRequired, &p.Username, &p.CurrentPasswordID, &p.PrimaryEmailAddressID)
}

func (p *Person) UpdatePerson(q sqlexp.Querier) error {
	_, err := q.ExecContext(
		context.TODO(),
		"UPDATE Person " +
		"SET last_modified_timestamp=$1, username=$2, primary_email_address_id=$3, multi_factor_required=$4 " +
		"WHERE id=$5",
		time.Now(), p.Username, p.PrimaryEmailAddressID, p.MultiFactorRequired, p.Id)
	return err
}

func (p *Person) DeletePerson(q sqlexp.Querier) error {
	_, err := q.ExecContext(
		context.TODO(),
		"UPDATE Person SET deleted_timestamp=$1 WHERE id=$2",
		time.Now(), p.Id)
	return err
}

func (p *Person) CreatePerson(q sqlexp.Querier) error {
	_, err := q.ExecContext(
		context.TODO(),
		"INSERT INTO Person(id, created_timestamp, username, disabled) VALUES($1, $2, $3, $4)",
		p.Id, time.Now(), p.Username, p.Disabled)

	if err != nil {
		return err
	}

	return nil
}

func GetPersons(q sqlexp.Querier, cursor int64, sort string, count int) ([]Person, error) {
	rows, err := q.QueryContext(
		context.TODO(),
		fmt.Sprintf(
			"SELECT p.id, p.created_timestamp, p.deleted_timestamp, p.last_modified_timestamp, p.disabled, " +
				"p.multi_factor_required, p.username, p.current_password_id, p.primary_email_address_id "+
			"FROM Person p " +
			"WHERE p.id > $1 " +
			"ORDER BY p.id %s " +
			"FETCH FIRST %d ROWS ONLY", sort, count), cursor)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	people := []Person{}

	for rows.Next() {
		var p Person
		if err := rows.Scan(&p.Id, &p.Created, &p.Deleted, &p.LastModified, &p.Disabled, &p.MultiFactorRequired,
			&p.Username, &p.CurrentPasswordID, &p.PrimaryEmailAddressID); err != nil {
			return nil, err
		}
		p.IdStr = strconv.FormatInt(p.Id, 10)
		if p.PrimaryEmailAddressID.Valid && p.PrimaryEmailAddressID.Int64 != 0.0 {
			p.PrimaryEmailAddressIdStr = strconv.FormatInt(p.PrimaryEmailAddressID.Int64, 10)
		}
		people = append(people, p)
	}

	return people, nil
}

func (p *Person) Exists(q sqlexp.Querier) (bool, error) {
	count, err := p.Count(q)
	if err != nil {
		return false, err
	}
	return count == 1, nil
}

func (p *Person) Count(q sqlexp.Querier) (int, error) {
	var count int
	row := q.QueryRowContext(context.TODO(), "SELECT COUNT(*) AS count FROM Person WHERE id=$1 AND deleted_timestamp IS NULL", p.Id)
	err := row.Scan(&count)
	return count, err
}

func (p *Person) UpdateCurrentPassword(q sqlexp.Querier, currentPasswordID int64) error {
	_, err := q.ExecContext(
		context.TODO(),
		"UPDATE Person " +
			"SET last_modified_timestamp=$1, current_password_id=$2 " +
			"WHERE id=$3",
		time.Now(), currentPasswordID, p.Id)
	return err
}