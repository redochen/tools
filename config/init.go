package config

import (
	"flag"
	"fmt"
	"github.com/larspensjo/config"
)

const (
	parameName     = "conf"
	defaultCfgFile = "app.conf"
	defaultSection = "DEFAULT"
)

var (
	Conf *Config = nil
)

func init() {

	Conf = new(Config)

	flag.StringVar(&Conf.FilePath, parameName, defaultCfgFile, "Generic Config File")
	flag.Parse()

	if Conf.FilePath == "" {
		fmt.Println("config file path is empty")
		return
	}

	var err error
	Conf.Config, err = config.ReadDefault(Conf.FilePath)
	if err != nil {
		fmt.Println(fmt.Sprintf("failed to load config file %s", Conf.FilePath), err)
	}
}
