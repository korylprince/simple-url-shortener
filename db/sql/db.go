package sql

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/korylprince/simple-url-shortener/db"
)

// DB represents a SQL database
type DB struct {
	db      *sql.DB
	codeLen int
}

// New returns a new DB with the given driver, DSN, and code length, or an error if one occurred
func New(driver, dsn string, codeLen int) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return &DB{db: db, codeLen: codeLen}, nil
}

func (d *DB) newCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, d.codeLen)
	for i := range b {
		b[i] = db.Charset[r.Intn(len(db.Charset))]
	}
	return string(b)
}

// Begin returns a transaction for the database or an error if one occurred
func (d *DB) Begin() (*sql.Tx, error) {
	return d.db.Begin()
}

// CreateURL adds the given URL to the database and returns a short code, or an error if one occurred
func (d *DB) CreateURL(tx *sql.Tx, url string) (string, error) {
	code := d.newCode()

	//check if code already exists
	u, err := d.ReadURL(tx, code)
	if err != nil {
		return "", fmt.Errorf("Unable to check code: %v", err)
	}
	//recurse until non-used code is generated
	if u != "" {
		return d.CreateURL(tx, url)
	}

	if _, err := tx.Exec("INSERT INTO url(code, url) VALUES(?, ?);", code, url); err != nil {
		return "", fmt.Errorf("Unable to insert URL: %v", err)
	}

	return code, nil
}

// ReadURL returns the URL with the given code, or an error if one occurred. If the returned url is empty, the code wasn't present in the database
func (d *DB) ReadURL(tx *sql.Tx, code string) (string, error) {

	row := tx.QueryRow("SELECT url FROM url WHERE code=?;", code)

	var url string
	err := row.Scan(&(url))

	switch {
	case err == sql.ErrNoRows:
		return "", nil
	case err != nil:
		return "", fmt.Errorf("Unable to query URL(%s): %v", code, err)
	}

	return url, nil
}
