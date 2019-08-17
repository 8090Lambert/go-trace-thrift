package helper

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

func FileOrDirExist(path string) bool {
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

func SetPid2File(pidFile string) {
	dir := filepath.Dir(pidFile)
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	pid := strconv.Itoa(os.Getpid())
	if err := ioutil.WriteFile(pidFile, []byte(pid), 0644); err != nil {
		log.Printf("set pid to file err:%s", err.Error())
	}
}

func CheckAddrAlreadyUse(addr string) (bool, error) {
	if addr == "" {
		return false, errors.New("addr can not be empty")
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	return true, nil
}
