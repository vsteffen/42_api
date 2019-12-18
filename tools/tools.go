package tools

import (
	"golang.org/x/crypto/ssh/terminal"
	"github.com/rs/zerolog/log"
	"syscall"
)

func ReadAndHideData() (string) {
	byteRead, err := terminal.ReadPassword(int(syscall.Stdin))
	if (err != nil) {
		log.Fatal().Err(err)
	}
	return string(byteRead)
}
