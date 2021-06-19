package vert

import (
	"fmt"
	"log"
	"os"
)

//MakeDir creates dir and logs with panic if it couldn't
func MakeDir(path string, mode uint32) {
	err := os.Mkdir(path, os.FileMode(mode))
	if err != nil {
		log.Fatal(fmt.Errorf("could not create path: %w", err))
	}
}
