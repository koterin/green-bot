package config

import (
	"log"

	"github.com/alexflint/go-arg"
)

var Args struct {
	LOG_LEVEL      string `arg:"required,env"`
	TG_BOT_KEY     string `arg:"required,env"`
	API_KEY        string `arg:"required,env"`
	AUTH_URL       string `arg:"required,env"`
	ORIGIN_URL     string `arg:"required,env"`
	NEW_ORIGIN_URL string `arg:"required,env"`
	NEW_USER_URL   string `arg:"required,env"`
}

func Validate() {
	if err := arg.Parse(&Args); err != nil {
		log.Fatal(err)
	}
}
