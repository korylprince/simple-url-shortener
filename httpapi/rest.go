package httpapi

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

func (s *Server) createURL(r *http.Request, tx *sql.Tx) (int, interface{}) {
	type request struct {
		URL string `json:"URL"`
	}

	type response struct {
		Code string `json:"code"`
	}

	req := new(request)

	if err := jsonRequest(r, req); err != nil {
		return http.StatusBadRequest, err
	}

	if !(govalidator.IsURL(req.URL)) {
		return http.StatusBadRequest, errors.New("Invalid URL")
	}

	code, err := s.db.CreateURL(tx, req.URL)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, &response{Code: code}
}

func (s *Server) readURL(r *http.Request, tx *sql.Tx) (int, interface{}) {
	type response struct {
		URL string `json:"url"`
	}

	code := mux.Vars(r)["code"]

	(r.Context().Value(contextKeyLogData)).(*logData).ActionID = code

	url, err := s.db.ReadURL(tx, code)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if url == "" {
		return http.StatusNotFound, fmt.Errorf("Unknown code: %s", code)
	}

	return http.StatusOK, &response{URL: url}
}
