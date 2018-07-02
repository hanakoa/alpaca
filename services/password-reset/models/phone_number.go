package models

import (
	"github.com/golang-sql/sqlexp"
	"context"
)

type PhoneNumber struct {
	ID          int64  `json:"id"`
	PhoneNumber string `json:"phone_number"`
	PersonID    int64  `json:"person_id"`
}

func (p *PhoneNumber) GetPhoneNumberByPhoneNumber(q sqlexp.Querier) error {
	return q.QueryRowContext(
		context.TODO(),
		"SELECT phone_number, person_id "+
			"FROM phone_number WHERE phone_number=$1 "+
			"AND deleted_timestamp IS NULL", p.PhoneNumber).Scan(&p.PhoneNumber, &p.PersonID)
}

func (p *PhoneNumber) MaskValue() {
	p.PhoneNumber = p.PhoneNumber[len(p.PhoneNumber)-2:]
}

func GetPhoneNumbersForPerson(personID int64, q sqlexp.Querier) ([]PhoneNumber, error) {
	rows, err := q.QueryContext(
		context.TODO(),
		"SELECT id, phone_number, person_id "+
			"FROM phone_number " +
			"WHERE confirmed=$1 AND person_id=$2 "+
			"AND deleted_timestamp IS NULL",
			true, personID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	phoneNumbers := []PhoneNumber{}

	for rows.Next() {
		var p PhoneNumber
		if err := rows.Scan(&p.ID, &p.PhoneNumber, &p.PersonID); err != nil {
			return nil, err
		}
		p.MaskValue()
		phoneNumbers = append(phoneNumbers, p)
	}

	return phoneNumbers, nil
}
