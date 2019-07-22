package httpapi

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	packr "github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/korylprince/simple-url-shortener/db"
)

//API is the current API version
const API = "1.0"
const apiPath = "/api/" + API

//Router returns a new API router
func (s *Server) Router() http.Handler {
	r := mux.NewRouter()
	r.SkipClean(true)

	api := r.PathPrefix(apiPath).Subrouter()

	api.NotFoundHandler = withJSONResponse(func(r *http.Request) (int, interface{}) {
		return http.StatusNotFound, nil
	})

	api.Methods("POST").Path("/urls").Handler(
		withLogging("CreateURL", s.output,
			withJSONResponse(
				withTX(s.db, s.createURL))))

	api.Methods("GET").Path(fmt.Sprintf("/urls/{code:[%s]+}", regexp.QuoteMeta(db.Charset))).Handler(
		withLogging("ReadURL", s.output,
			withJSONResponse(
				withTX(s.db, s.readURL))))

	r.Methods("GET").Path("/in/{url:.*}").Handler(
		withLogging("LegacyCreateURL", s.output,
			withTextResponse(
				withTX(s.db, s.createURLLegacy))))

	r.Methods("GET").Path(fmt.Sprintf("/{code:[%s]+}", regexp.QuoteMeta(db.Charset))).Handler(
		withLogging("LegacyReadURL", s.output,
			withRedirectResponse(
				withTX(s.db, s.readURLLegacy))))

	box := packr.New("Static", "./static")
	r.Methods("GET").Path("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idx, err := box.Find("index.html")
		if err != nil {
			panic("expected to find index.html")
		}
		if _, err := w.Write(idx); err != nil {
			log.Println("Unable to write HTTP response:", err)
		}
	}))

	return r
}
