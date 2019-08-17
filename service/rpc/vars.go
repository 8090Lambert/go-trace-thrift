package rpc

import (
	"fmt"
	"github.com/8090Lambert/tracer_thrift/basepath"
	"github.com/8090Lambert/tracer_thrift/config"
	"log"
	"os"
	"syscall"
)

var vars *runningVars

type runningVars struct {
	ip          string      // ip address, default 0.0.0.0
	port        string      // port, default :8300
	timeout     int64       // Timeout(ms), default 700
	pidFile     string      // pidFile, default
	captureSign []os.Signal // graceful shutdown
}

func init() {
	configPath, _ := basepath.Path.ConfigPath()
	serverConf, err := config.NewWrapperConf(configPath+"/service.ini", "service")
	if err != nil {
		log.Fatal("load service conf failed, err:" + err.Error())
		os.Exit(-1)
	}

	vars = new(runningVars)
	if serverConf.Ip == "" {
		vars.ip = "0.0.0.0"
	}
	vars.ip = serverConf.Ip

	if serverConf.Port == "" {
		vars.port = ":8300"
	}
	vars.port = serverConf.Port

	if serverConf.Timeout == 0 {
		vars.timeout = 700
	}
	vars.timeout = serverConf.Timeout

	pidFile := serverConf.Pid
	if pidFile == "" {
		pidFile = "trace_thrift.pid"
	}
	vars.pidFile = fmt.Sprintf("var/run/%s", pidFile)

	vars.captureSign = []os.Signal{syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT, syscall.SIGUSR2, syscall.SIGQUIT}
}
