package httpapi

import (
	"io"

	"github.com/korylprince/simple-url-shortener/db"
)

// Server represents shared resources
type Server struct {
	db     db.DB
	output io.Writer
	prefix string
}

// NewServer returns a new server with the given resources
func NewServer(db db.DB, output io.Writer, prefix string) *Server {
	return &Server{db: db, output: output, prefix: prefix}
}
