package db

import (
	"database/sql"
)

// Charset is the valid charset for a URL code
const Charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-"

// DB represents a database
type DB interface {
	//Begin returns a transaction for the database or an error if one occurred
	Begin() (*sql.Tx, error)

	//CreateURL adds the given URL to the database and returns a short code, or an error if one occurred
	CreateURL(tx *sql.Tx, url string) (string, error)

	//ReadURL returns the URL with the given code, or an error if one occurred. If the returned url is empty, the code wasn't present in the database
	ReadURL(tx *sql.Tx, code string) (string, error)
}
