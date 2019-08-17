package config

import (
	"github.com/8090Lambert/tracer_thrift/basepath"
	"testing"
)

func TestNewWrapperConf(t *testing.T) {
	configPath, _ := basepath.Path.RootPath()
	config, err := NewWrapperConf(configPath+"/../conf/service.ini", "service")
	if err != nil || config.Ip != "0.0.0.0" || config.Port != "8300" {
		t.Errorf("test newwrapperconf failed")
	}
}
