package httpapi

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

func (s *Server) createURLLegacy(r *http.Request, tx *sql.Tx) (int, interface{}) {
	url := mux.Vars(r)["url"]

	if !(govalidator.IsURL(url)) {
		return http.StatusBadRequest, errors.New("Invalid URL")
	}

	code, err := s.db.CreateURL(tx, url)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	scheme := "http"
	if s := r.Header.Get("X-Forwarded-Proto"); s != "" {
		scheme = s
	}

	u, err := r.URL.Parse(fmt.Sprintf("%s://%s%s/%s", scheme, r.Host, s.prefix, code))
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Unable to create URL: %v", err)
	}

	return http.StatusOK, u.String()
}

func (s *Server) readURLLegacy(r *http.Request, tx *sql.Tx) (int, interface{}) {
	code := mux.Vars(r)["code"]

	(r.Context().Value(contextKeyLogData)).(*logData).ActionID = code

	url, err := s.db.ReadURL(tx, code)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if url == "" {
		return http.StatusNotFound, fmt.Errorf("Unknown code: %s", code)
	}

	return http.StatusOK, url
}
