package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/korylprince/simple-url-shortener/db/sql"
	"github.com/korylprince/simple-url-shortener/httpapi"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.New(config.SQLDriver, config.SQLDSN, config.CodeLength)
	if err != nil {
		log.Fatalln("Could not open database:", err)
	}

	httpapi.Debug = config.Debug

	s := httpapi.NewServer(db, os.Stdout, config.Prefix)

	log.Println("Listening on:", config.ListenAddr)

	log.Println(http.ListenAndServe(config.ListenAddr, http.StripPrefix(config.Prefix, s.Router())))
}
