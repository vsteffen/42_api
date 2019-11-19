package main

import (
	"log"
	"os"
	"github.com/vsteffen/42_api/constants"
)

type Keys struct {
	uid string
	secret string
	token string
}

var (
	g_keys = Keys{constants.UID, "",""}
)

func read_keys(path_file string) (string, error) {
	buff := make([]byte, constants.SizeKeys)

	file, err := os.Open(path_file)
	if (err == nil) {
		read_size, err := file.Read(buff)
		if (err != nil) {
			log.Fatal(err)
		} else if (read_size != constants.SizeKeys) {
			log.Fatal("wrong size for " + path_file + "\n")
		}
		log.Print(constants.ProjectName + ": " + path_file + ": hash found")
		return string(buff), nil
	}
	return "", err
}

func init_keys() {
	var err error
	if g_keys.secret, err = read_keys(constants.PathSecret); err != nil {
		log.Fatal(err)
	}
	if g_keys.token, err = read_keys(constants.PathToken); err != nil {
		log.Print(err)
	}
}