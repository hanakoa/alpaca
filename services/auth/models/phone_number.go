package models

import (
	"github.com/golang-sql/sqlexp"
	"context"
)

type PhoneNumber struct {
	ID          int64  `json:"id"`
	PhoneNumber string `json:"phone_number"`
	AccountID    int64  `json:"account_id"`
}

func (p *PhoneNumber) GetPhoneNumberByPhoneNumber(q sqlexp.Querier) error {
	return q.QueryRowContext(
		context.TODO(),
		"SELECT phone_number, account_id "+
			"FROM phone_number WHERE phone_number=$1 "+
			"AND deleted_timestamp IS NULL", p.PhoneNumber).Scan(&p.PhoneNumber, &p.AccountID)
}

func (p *PhoneNumber) MaskValue() {
	p.PhoneNumber = p.PhoneNumber[len(p.PhoneNumber)-2:]
}

func GetPhoneNumbersForAccount(accountID int64, q sqlexp.Querier) ([]PhoneNumber, error) {
	rows, err := q.QueryContext(
		context.TODO(),
		"SELECT id, phone_number, account_id "+
			"FROM phone_number " +
			"WHERE confirmed=$1 AND account_id=$2 "+
			"AND deleted_timestamp IS NULL",
			true, accountID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	phoneNumbers := []PhoneNumber{}

	for rows.Next() {
		var p PhoneNumber
		if err := rows.Scan(&p.ID, &p.PhoneNumber, &p.AccountID); err != nil {
			return nil, err
		}
		p.MaskValue()
		phoneNumbers = append(phoneNumbers, p)
	}

	return phoneNumbers, nil
}
