package main

import (
	_ "github.com/vsteffen/42_api/tools/constants"
	"os"
	"time"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vsteffen/42_api/reqApi42"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Stamp})
	reqApi42.New()
}
