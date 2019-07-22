package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

//Config represents options given in the environment
type Config struct {
	CodeLength int `default:"6"`

	SQLDriver string `required:"true"`
	SQLDSN    string `required:"true"`

	ListenAddr string `default:":8080" required:"true"` //addr format used for net.Dial; required
	Prefix     string //url prefix to mount api to without trailing slash
	Debug      bool   `default:"false"` //return debugging information to client
}

var config = &Config{}

func init() {
	err := envconfig.Process("SHORTENER", config)
	if err != nil {
		log.Fatalln("Error reading configuration from environment:", err)
	}
}
