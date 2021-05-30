package serato_parser

import (
	"os/user"
	"path/filepath"
)

var USER_HOME_DIR string
var SERATO_DIR string

func init() {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	USER_HOME_DIR = filepath.FromSlash(currentUser.HomeDir)
	SERATO_DIR = filepath.FromSlash(USER_HOME_DIR + "/Music/_Serato_")
}
