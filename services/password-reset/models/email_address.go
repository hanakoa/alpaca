package models

import (
	"github.com/golang-sql/sqlexp"
	"context"
	"strings"
)

type EmailAddress struct {
	ID           int64  `json:"id"`
	EmailAddress string `json:"email_address"`
	AccountID     int64  `json:"account_id"`
}

func (e *EmailAddress) GetConfirmedEmailAddress(q sqlexp.Querier) error {
	return q.QueryRowContext(
		context.TODO(),
		"SELECT email_address, account_id "+
			"FROM email_address WHERE email_address=$1 " +
			"AND confirmed=$2 "+
			"AND deleted_timestamp IS NULL", e.EmailAddress, true).Scan(&e.EmailAddress, &e.AccountID)
}

func (e *EmailAddress) getMaskedEmailUser() string {
	splits := strings.Split(e.EmailAddress, "@")
	user := splits[0]
	if len(user) == 1 {
		return user[0:1] + strings.Repeat("*", len(user) - 1)
	}
	return user[0:2] + strings.Repeat("*", len(user) - 2)
}

func (e *EmailAddress) getMaskedEmailHost() string {
	emailSplits := strings.Split(e.EmailAddress, "@")
	host := emailSplits[1]
	splits := strings.Split(host, ".")
	splits[0] = splits[0][0:1] + strings.Repeat("*", len(splits[0]) - 1)
	return strings.Join(splits, ".")
}

func (e *EmailAddress) MaskValue() {
	// TODO is it possible for email to be empty or null?
	e.EmailAddress = e.getMaskedEmailUser() + "@" + e.getMaskedEmailHost()
}

func GetEmailAddressesForAccount(accountID int64, q sqlexp.Querier) ([]EmailAddress, error) {
	rows, err := q.QueryContext(
		context.TODO(),
		"SELECT id, email_address, account_id "+
			"FROM email_address " +
			"WHERE confirmed=$1 AND account_id=$2 "+
			"AND deleted_timestamp IS NULL",
		true, accountID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	emailAddresses := []EmailAddress{}

	for rows.Next() {
		var e EmailAddress
		if err := rows.Scan(&e.ID, &e.EmailAddress, &e.AccountID); err != nil {
			return nil, err
		}
		e.MaskValue()
		emailAddresses = append(emailAddresses, e)
	}

	return emailAddresses, nil
}