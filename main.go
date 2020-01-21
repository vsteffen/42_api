package main

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vsteffen/42_api/reqApi42"
	_ "github.com/vsteffen/42_api/tools/constants"
	"os"
	"time"
)

func main() {
	flagRefresh := flag.Bool("refresh", false, "force to refresh token")
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Stamp})
	argRefresh := *flagRefresh
	api42 := reqApi42.New(argRefresh)
	_ = api42
}
