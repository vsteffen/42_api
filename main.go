package main

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vsteffen/42_api/reqApi42"
	cst "github.com/vsteffen/42_api/tools/constants"
	"os"
	"time"
)

func main() {
	flagRefresh := flag.Bool("refresh", false, "force to refresh token")
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Stamp})
	api42 := reqApi42.New()
	if *flagRefresh == true {
		log.Info().Msg("Force to refresh token")
		api42.RefreshToken()
	}
	api42.GetCampusID(cst.CampusName)
	api42.UpdateLocations()
}
