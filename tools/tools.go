package tools

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
)

// ReadAndHideData read on stdin and hide user input
func ReadAndHideData() string {
	byteRead, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal().Err(err)
	}
	return string(byteRead)
}
