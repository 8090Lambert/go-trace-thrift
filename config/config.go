package config

import (
	"errors"
	"fmt"
	"github.com/8090Lambert/tracer_thrift/helper"
	"github.com/go-ini/ini"
	"strings"
)

type Config struct {
	Ip      string `ini:"ip"`
	Port    string `ini:"port"`
	Timeout int64  `ini:"timeout"`
	Pid     string `ini:"pid"`
}

func NewWrapperConf(file string, section string) (*Config, error) {
	confFile := fmt.Sprintf("%s.ini", strings.TrimRight(file, ".ini"))
	if helper.FileOrDirExist(confFile) == false {
		return nil, errors.New("must set " + confFile)
	}

	handler, err := ini.Load(confFile)
	config := new(Config)
	err = handler.Section(section).MapTo(config)

	return config, err
}
