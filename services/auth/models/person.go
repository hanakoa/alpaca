package models

import (
	"gopkg.in/guregu/null.v3"
	"strconv"
	"time"
	"fmt"
	"github.com/golang-sql/sqlexp"
	"context"
)

type Account struct {
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
}

func (p *Account) GetAccountByUsername(q sqlexp.Querier) error {
	return q.QueryRowContext(
		context.TODO(),
		"SELECT p.id, p.created_timestamp, p.deleted_timestamp, p.last_modified_timestamp, p.disabled, "+
			"p.multi_factor_required, p.username, p.current_password_id, p.primary_email_address_id "+
			"FROM Account p "+
			"WHERE p.username=$1 " +
			"AND p.deleted_timestamp IS NULL", p.Username).Scan(&p.Id, &p.Created, &p.Deleted, &p.LastModified, &p.Disabled,
				&p.MultiFactorRequired, &p.Username, &p.CurrentPasswordID, &p.PrimaryEmailAddressID)
}

func GetAccountByEmailAddress(q sqlexp.Querier, emailAddress string) (*Account, error) {
	p := &Account{}
	return p, q.QueryRowContext(
		context.TODO(),
		"SELECT p.id, p.created_timestamp, p.deleted_timestamp, p.last_modified_timestamp, p.disabled, "+
			"p.multi_factor_required, p.username, p.current_password_id, p.primary_email_address_id " +
			"FROM email_address e JOIN Account p ON e.accountId = p.id "+
			"WHERE e.email_address=$1 " +
			"AND p.deleted_timestamp IS NULL", emailAddress).Scan(&p.Id, &p.Created, &p.Deleted, &p.LastModified,
				&p.Disabled, &p.MultiFactorRequired, &p.Username, &p.CurrentPasswordID,
					&p.PrimaryEmailAddressID)
}

func GetAccountByPhoneNumber(q sqlexp.Querier, phoneNumber string) (*Account, error) {
	p := &Account{}
	return p, q.QueryRowContext(
		context.TODO(),
		"SELECT p.id, p.created_timestamp, p.deleted_timestamp, p.last_modified_timestamp, p.disabled, "+
			"p.multi_factor_required, p.username, p.current_password_id, p.primary_email_address_id " +
			"FROM phone_number pn JOIN Account p ON pn.accountId = p.id "+
			"WHERE pn.phone_number=$1 " +
			"AND p.deleted_timestamp IS NULL", phoneNumber).Scan(&p.Id, &p.Created, &p.Deleted, &p.LastModified,
		&p.Disabled, &p.MultiFactorRequired, &p.Username, &p.CurrentPasswordID,
		&p.PrimaryEmailAddressID)
}

func (p *Account) GetDeletedAccount(q sqlexp.Querier) error {
	return q.QueryRowContext(
		context.TODO(),
		"SELECT p.id, p.created_timestamp, p.deleted_timestamp, p.last_modified_timestamp, p.disabled, " +
			"p.multi_factor_required, p.username, p.current_password_id, p.primary_email_address_id "+
			"FROM Account p "+
			"WHERE p.id=$1 "+
			"AND p.deleted_timestamp IS NOT NULL", p.Id).Scan(&p.Id, &p.Created, &p.Deleted, &p.LastModified,
				&p.Disabled, &p.MultiFactorRequired, &p.Username, &p.CurrentPasswordID,
					&p.PrimaryEmailAddressID)
}

func (p *Account) GetAccount(q sqlexp.Querier) error {
	return q.QueryRowContext(
		context.TODO(),
		"SELECT p.id, p.created_timestamp, p.deleted_timestamp, p.last_modified_timestamp, p.disabled, " +
			"p.multi_factor_required, p.username, p.current_password_id, p.primary_email_address_id "+
			"FROM Account p "+
			"WHERE p.id=$1 "+
			"AND p.deleted_timestamp IS NULL", p.Id).Scan(&p.Id, &p.Created, &p.Deleted, &p.LastModified, &p.Disabled,
		&p.MultiFactorRequired, &p.Username, &p.CurrentPasswordID, &p.PrimaryEmailAddressID)
}

func (p *Account) UpdateAccount(q sqlexp.Querier) error {
	_, err := q.ExecContext(
		context.TODO(),
		"UPDATE Account " +
		"SET last_modified_timestamp=$1, username=$2, primary_email_address_id=$3, multi_factor_required=$4 " +
		"WHERE id=$5",
		time.Now(), p.Username, p.PrimaryEmailAddressID, p.MultiFactorRequired, p.Id)
	return err
}

func (p *Account) DeleteAccount(q sqlexp.Querier) error {
	_, err := q.ExecContext(
		context.TODO(),
		"UPDATE Account SET deleted_timestamp=$1 WHERE id=$2",
		time.Now(), p.Id)
	return err
}

func (p *Account) CreateAccount(q sqlexp.Querier) error {
	_, err := q.ExecContext(
		context.TODO(),
		"INSERT INTO Account(id, created_timestamp, username, disabled) VALUES($1, $2, $3, $4)",
		p.Id, time.Now(), p.Username, p.Disabled)

	if err != nil {
		return err
	}

	return nil
}

func GetAccounts(q sqlexp.Querier, cursor int64, sort string, count int) ([]Account, error) {
	rows, err := q.QueryContext(
		context.TODO(),
		fmt.Sprintf(
			"SELECT p.id, p.created_timestamp, p.deleted_timestamp, p.last_modified_timestamp, p.disabled, " +
				"p.multi_factor_required, p.username, p.current_password_id, p.primary_email_address_id "+
			"FROM Account p " +
			"WHERE p.id > $1 " +
			"ORDER BY p.id %s " +
			"FETCH FIRST %d ROWS ONLY", sort, count), cursor)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	accounts := []Account{}

	for rows.Next() {
		var p Account
		if err := rows.Scan(&p.Id, &p.Created, &p.Deleted, &p.LastModified, &p.Disabled, &p.MultiFactorRequired,
			&p.Username, &p.CurrentPasswordID, &p.PrimaryEmailAddressID); err != nil {
			return nil, err
		}
		p.IdStr = strconv.FormatInt(p.Id, 10)
		if p.PrimaryEmailAddressID.Valid && p.PrimaryEmailAddressID.Int64 != 0.0 {
			p.PrimaryEmailAddressIdStr = strconv.FormatInt(p.PrimaryEmailAddressID.Int64, 10)
		}
		accounts = append(accounts, p)
	}

	return accounts, nil
}

func (p *Account) Exists(q sqlexp.Querier) (bool, error) {
	count, err := p.Count(q)
	if err != nil {
		return false, err
	}
	return count == 1, nil
}

func (p *Account) Count(q sqlexp.Querier) (int, error) {
	var count int
	row := q.QueryRowContext(context.TODO(), "SELECT COUNT(*) AS count FROM Account WHERE id=$1 AND deleted_timestamp IS NULL", p.Id)
	err := row.Scan(&count)
	return count, err
}

func (p *Account) UpdateCurrentPassword(q sqlexp.Querier, currentPasswordID int64) error {
	_, err := q.ExecContext(
		context.TODO(),
		"UPDATE Account " +
			"SET last_modified_timestamp=$1, current_password_id=$2 " +
			"WHERE id=$3",
		time.Now(), currentPasswordID, p.Id)
	return err
}