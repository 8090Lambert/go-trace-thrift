package basepath

import (
	"errors"
	"log"
	"os"
	"strings"
)

var Path *dir

type dir struct {
	root string
}

func init() {
	base, err := os.Getwd()
	if err != nil {
		log.Fatal("load base path failed!")
		os.Exit(-1)
	}

	Path = new(dir)
	Path.root = base
}

func (d *dir) RootPath() (string, error) {
	return strings.TrimRight(d.root, "/"), nil
}

func (d *dir) ConfigPath() (string, error) {
	path := strings.TrimRight(d.root, "/") + "/" + "conf"
	if dirExist(path) != true {
		return "", errors.New("conf path failed")
	}

	return path, nil
}

func dirExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}

	return true
}
