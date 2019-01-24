package logger

import (
	"brief_framework/plugin/log4go"
	"os"
	"brief_framework/config"
	"brief_framework/util"
)

type Logger struct {
	log4go.Logger
}

var serverLogger *Logger
var sessionLogger *Logger

func (l *Logger) init() {
	l.Logger = make(log4go.Logger)
}

func (l *Logger) Write(p []byte) (n int, err error) {
	l.Error("%s", string(p))
	return len(p), nil
}

func getDefaultLogSet(file string) *log4go.FileLogWriter {
	flw := log4go.NewFileLogWriter(file, false)
	flw.SetFormat("[%D %T] [%L] (%S) %M")
	flw.SetRotate(true)
	flw.SetRotateLines(0)
	flw.SetRotateSize(0)
	flw.SetRotateDaily(true)
	return flw
}

func init() {
	serverLogger = new(Logger)
	serverLogger.init()

	mode := config.RunningMode()
	serverLogConf := config.Instance().MustValue(mode + ".log", "server_log_conf", "./conf/server.xml")
	_, err := os.Stat(serverLogConf)
	if err != nil {
		//fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not open %q for reading: %s\n", filename, err)
		serverLogger.LoadConfiguration(util.GetBaseDirectory() + "/" + serverLogConf)
	} else {
		serverLogger.LoadConfiguration(serverLogConf)
	}

	sessionLogger = new(Logger)
	sessionLogger.init()
	sessionLogConf := config.Instance().MustValue(mode + ".log", "session_log_conf", "./conf/session.xml")
	_, err = os.Stat(sessionLogConf)
	if err != nil {
		//fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not open %q for reading: %s\n", filename, err)
		serverLogger.LoadConfiguration(util.GetBaseDirectory() + "/"  + sessionLogConf)
	} else {
		sessionLogger.LoadConfiguration(sessionLogConf)
	}
	//sessionLogger.AddFilter("sessionlog", log4go.FINE, getDefaultLogSet("./logs/session.log"))
}

func Instance() *Logger {
	return serverLogger
}

//Do not use this
func GetSessionLogger() *Logger {
	return sessionLogger
}
