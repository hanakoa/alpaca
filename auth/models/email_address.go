package models

import (
	"gopkg.in/guregu/null.v3"
	"strconv"
	"fmt"
	"time"
	"github.com/golang-sql/sqlexp"
	"context"
)

// TODO add newEmailAddress method
// TODO manually manage timestamps
// TODO return error if PUT/POST includes more fields than it should? or do we quietly process only the fields we care about

// EmailAddress is a representation of a user's email address.
type EmailAddress struct {
	Id           int64     `json:"id"`
	IdStr        string    `json:"id_str"`
	Created      null.Time `json:"created_at"`
	Deleted      null.Time `json:"deleted_at"`
	LastModified null.Time `json:"last_modified_at"`
	Confirmed    bool      `json:"confirmed"`
	Primary      bool      `json:"primary"`
	EmailAddress string    `json:"email_address"`
	PersonId     int64     `json:"person_id"`
	PersonIdStr  string    `json:"person_id_str"`
}

func (e *EmailAddress) GetEmailAddressByEmailAddress(q sqlexp.Querier) error {
	return q.QueryRowContext(
		context.TODO(),
		"SELECT id, created_timestamp, deleted_timestamp, last_modified_timestamp, confirmed, is_primary, email_address, person_id "+
			"FROM email_address WHERE email_address=$1 "+
			"AND deleted_timestamp IS NULL", e.EmailAddress).Scan(&e.Id, &e.Created, &e.Deleted, &e.LastModified, &e.Confirmed, &e.Primary, &e.EmailAddress, &e.PersonId)
}

func (e *EmailAddress) GetEmailAddress(q sqlexp.Querier) error {
	return q.QueryRowContext(
		context.TODO(),
		"SELECT id, created_timestamp, deleted_timestamp, last_modified_timestamp, confirmed, is_primary, email_address, person_id " +
			"FROM email_address WHERE id=$1 " +
			"AND deleted_timestamp IS NULL", e.Id).Scan(&e.Id, &e.Created, &e.Deleted, &e.LastModified, &e.Confirmed, &e.Primary, &e.EmailAddress, &e.PersonId)
}

func (e *EmailAddress) GetDeletedEmailAddress(q sqlexp.Querier) error {
	return q.QueryRowContext(
		context.TODO(),
		"SELECT id, created_timestamp, deleted_timestamp, last_modified_timestamp, confirmed, email_address, person_id " +
			"FROM email_address WHERE id=$1 " +
			"AND deleted_timestamp IS NOT NULL", e.Id).Scan(&e.Id, &e.Created, &e.Deleted, &e.LastModified, &e.Confirmed, &e.EmailAddress, &e.PersonId)
}

// UpdateEmailAddress updates only the confirmation status of an email address.
func (e *EmailAddress) UpdateEmailAddress(q sqlexp.Querier) error {
	_, err := q.ExecContext(
		context.TODO(),
		"UPDATE email_address SET last_modified_timestamp=$1, confirmed=$2 WHERE id=$3",
		time.Now(), e.Confirmed, e.Id)
	return err
}

func (e *EmailAddress) DeleteEmailAddress(q sqlexp.Querier) error {
	_, err := q.ExecContext(context.TODO(), "DELETE FROM email_address WHERE id=$1", e.Id)
	return err
}

func (e *EmailAddress) CreateEmailAddress(q sqlexp.Querier) error {
	_, err := q.ExecContext(
		context.TODO(),
		"INSERT INTO email_address(id, person_id, email_address, confirmed, is_primary) VALUES($1, $2, $3, $4, $5)",
		e.Id, e.PersonId, e.EmailAddress, e.Confirmed, e.Primary)

	return err
}

func GetEmailAddresses(q sqlexp.Querier, cursor int64, sort string, count int) ([]EmailAddress, error) {
	rows, err := q.QueryContext(
		context.TODO(),
		fmt.Sprintf(
			"SELECT id, created_timestamp, deleted_timestamp, last_modified_timestamp, confirmed, email_address, person_id "+
				"FROM email_address " +
				"WHERE id > $1 " +
				"AND deleted_timestamp IS NULL " +
				"ORDER BY id %s " +
				"FETCH FIRST %d ROWS ONLY", sort, count), cursor)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	emailAddresses := []EmailAddress{}

	for rows.Next() {
		var e EmailAddress
		if err := rows.Scan(&e.Id, &e.Created, &e.Deleted, &e.LastModified, &e.Confirmed, &e.EmailAddress, &e.PersonId); err != nil {
			return nil, err
		}
		e.IdStr = strconv.FormatInt(e.Id, 10)
		e.PersonIdStr = strconv.FormatInt(e.PersonId, 10)
		emailAddresses = append(emailAddresses, e)
	}

	return emailAddresses, nil
}

func (e *EmailAddress) IsConfirmed(q sqlexp.Querier) (bool, error) {
	var count int
	row := q.QueryRowContext(
		context.TODO(),
		"SELECT COUNT(*) AS count " +
		"FROM email_address " +
		"WHERE email_address = $1 " +
		"AND confirmed = $2 " +
		"AND deleted_timestamp IS NULL", e.EmailAddress, true)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 1, nil
}

func (e *EmailAddress) Exists(q sqlexp.Querier) (bool, error) {
	count, err := e.Count(q)
	if err != nil {
		return false, err
	}
	return count == 1, nil
}

func (e *EmailAddress) Count(q sqlexp.Querier) (int, error) {
	var count int
	row := q.QueryRowContext(context.TODO(), "SELECT COUNT(*) AS count FROM email_address WHERE id=$1 AND deleted_timestamp IS NULL", e.Id)
	err := row.Scan(&count)
	return count, err
}