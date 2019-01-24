package config

import (
	"flag"
	"github.com/Unknwon/goconfig"
	"os"
	"brief_framework/util"
)


var (
	conf *goconfig.ConfigFile
	mode, configFile, baseDir string
)

func init() {
	var err error

	flag.StringVar(&mode, "m", "online", "running mode")
	flag.StringVar(&configFile, "c", "../conf/config.ini", "configuration file.")
	flag.Parse()

	_, err = os.Stat(configFile)
	if err != nil {
		conf, err = goconfig.LoadConfigFile(util.GetBaseDirectory() + "/" + configFile)
	} else {
		conf, err = goconfig.LoadConfigFile(configFile)
	}

	if err != nil {
		panic(err)
	}

}

func Instance() *goconfig.ConfigFile {
	return conf
}

func SetMode(m string) {
	mode = m
}

func RunningMode() string {
	return mode
}
