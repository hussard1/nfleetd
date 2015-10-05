package nfleetd

import (
	"os"
)

type Configuration struct {

}

func (conf Configuration) Load() error {
	path, err := os.Getwd()
	if err != nil {

	}
	return err
}