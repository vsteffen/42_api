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
	flags := []interface{}{}
	flags = append(flags, flag.Bool("refresh", false, "force to refresh token"))
	flags = append(flags, flag.Bool("check-default-values", false, "send a request to verify the default values"))
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Stamp})
	api42 := reqApi42.New(flags)
	// api42.GetCursusProjects()
	api42.UpdateLocations()
	api42.UpdateLocations()
	api42.UpdateLocations()
	api42.UpdateLocations()
	api42.UpdateLocations()
	time.Sleep(time.Second * 2)
	api42.UpdateLocations()
	api42.UpdateLocations()
}
