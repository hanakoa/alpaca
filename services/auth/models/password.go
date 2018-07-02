package models

import (
	"log"
	"encoding/hex"
	"time"
	"golang.org/x/crypto/pbkdf2"
	"crypto/sha1"
	"io"
	"crypto/rand"
	"gopkg.in/guregu/null.v3"
	"github.com/golang-sql/sqlexp"
	"context"
)

// Password is a representation of a user's password.
type Password struct {
	Id             int64       `json:"id"`
	IdStr          string      `json:"id_str"`
	Created        null.Time   `json:"created_at"`
	// TODO add scheme, e.g., PBKDF2-HMAC-SHA1
	IterationCount int
	Salt           []byte
	PasswordHash   []byte
	PersonID       int64       `json:"person_id"`
	PersonIdStr    string      `json:"person_id_str"`
	PasswordText   null.String `json:"password"`
}

// CalibrateIterationCount finds the given
func CalibrateIterationCount(hashTime time.Duration) int {
	salt, err := generateSalt(32)
	if err != nil {
		log.Fatal(err)
	}

	iterationCount := 10000
	log.Printf("Calibrating password iteration count, starting at %d iterations...\n", iterationCount)
	for {
		start := time.Now()
		GenerateHash("MyPassword123!", iterationCount, salt)
		elapsed := time.Since(start)
		if elapsed > hashTime {
			log.Printf("Took %s to do %d iterations\n", elapsed, iterationCount)
			break
		}

		percentage := elapsed.Seconds() / hashTime.Seconds()
		if percentage < 0.2 {
			log.Println("Less than 20% of the way there...")
			iterationCount = iterationCount * 4
		} else if percentage < 0.3 {
			log.Println("Less than 30% of the way there...")
			iterationCount = iterationCount * 3
		} else if percentage < 0.4 {
			log.Println("Less than 40% of the way there...")
			iterationCount = iterationCount * 2
		} else if percentage < 0.5 {
			log.Println("Less than 50% of the way there...")
			iterationCount = int(float64(iterationCount) * 1.75)
		} else if percentage < 0.6 {
			log.Println("Less than 60% of the way there...")
			iterationCount = int(float64(iterationCount) * 1.55)
		} else if percentage < 0.7 {
			log.Println("Less than 70% of the way there...")
			iterationCount = int(float64(iterationCount) * 1.35)
		} else if percentage < 0.8 {
			log.Println("Less than 80% of the way there...")
			iterationCount = int(float64(iterationCount) * 1.20)
		} else if percentage < 0.9 {
			log.Println("Less than 90% of the way there...")
			iterationCount = int(float64(iterationCount) * 1.07)
		} else if percentage < 0.95 {
			log.Println("Less than 95% of the way there...")
			iterationCount = int(float64(iterationCount) * 1.04)
		} else {
			log.Printf("We're close enough... Took %s to do %d iterations\n", elapsed, iterationCount)
			break
		}
	}
	return iterationCount
}

func (p *Password) GetPasswordForPersonID(q sqlexp.Querier) error {
	return q.QueryRowContext(
		context.TODO(),
		"SELECT p.id, p.created_timestamp, p.iteration_count, p.salt, " +
			"p.password_hash, p.person_id "+
			"FROM Password p "+
			"WHERE p.id=$1", p.Id).Scan(&p.Id, &p.Created, &p.IterationCount, &p.Salt, &p.PasswordHash,
		&p.PersonID)
}

func (p *Password) CreatePassword(q sqlexp.Querier, iterationCount int) error {
	now := time.Now()
	var salt, passwordHash []byte
	salt, err := generateSalt(32)
	if err != nil {
		return err
	}

	passwordHash = GenerateHash(p.PasswordText.String, iterationCount, salt)
	_, err = q.ExecContext(
		context.TODO(),
		"INSERT INTO Password(id, created_timestamp, iteration_count, salt, password_hash, person_id) VALUES($1, $2, $3, $4, $5, $6)",
		p.Id, now, iterationCount, salt, passwordHash, p.PersonID)

	return err
}

func (p *Password) UpdatePassword(q sqlexp.Querier) error {
	_, err := q.ExecContext(
		context.TODO(),
		"UPDATE Password SET iteration_count=$1, salt=decode($2, 'hex'), password_hash=decode($3, 'hex') WHERE id=$4",
		p.IterationCount,
		hex.EncodeToString(p.Salt),
		hex.EncodeToString(p.PasswordHash),
		p.Id)
	return err
}

func MatchesHash(passwordText string, password *Password) bool {
	hash := GenerateHash(passwordText, int(password.IterationCount), password.Salt)
	return hex.EncodeToString(password.PasswordHash) == hex.EncodeToString(hash)
}

func GenerateHash(passwordText string, iterationCount int, salt []byte) []byte {
	start := time.Now()
	// TODO use switch statement on scheme string
	hash := pbkdf2.Key([]byte(passwordText), salt, iterationCount, 32, sha1.New)
	log.Printf("Hashing %d iterations took: %s", iterationCount, time.Since(start))
	return hash
}

func generateSalt(byteLength int) ([]byte, error) {
	salt := make([]byte, byteLength)
	_, err := io.ReadFull(rand.Reader, salt)
	return salt, err
}
