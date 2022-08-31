package config

import (
	"log"

	"github.com/alexflint/go-arg"
)

var Args struct {
	TG_BOT_KEY string `arg:"required,env"`
	API_KEY    string `arg:"required,env"`
	AUTH_URL   string `arg:"required,env"`
}

func Validate() {
	if err := arg.Parse(&Args); err != nil {
		log.Fatal(err)
	}
}