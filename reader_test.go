package seratoparser

import (
	"os/user"
	"path/filepath"
)

var UserHomeDir string
var SeratoDir string

func init() {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	UserHomeDir = filepath.FromSlash(currentUser.HomeDir)
	SeratoDir = filepath.FromSlash(UserHomeDir + "/Music/_Serato_")
}
