package config

import (
	"flag"
	"fmt"
	"github.com/larspensjo/config"
	"log"
	"os"
	"path"
	"path/filepath"
)

const (
	paramName         = "conf"
	defaultConfigFile = "app.conf"
	defaultSection    = "DEFAULT"
)

var (
	Conf *Config = nil
)

func init() {
	Conf = new(Config)

	flag.StringVar(&Conf.FilePath, paramName, defaultConfigFile, "Generic Config File")
	flag.Parse()

	if Conf.FilePath == "" {
		fmt.Println("config file path is empty")
		return
	}

	if Conf.FilePath == defaultConfigFile {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}

		Conf.FilePath = path.Join(dir, Conf.FilePath)
	}

	var err error
	Conf.Config, err = config.ReadDefault(Conf.FilePath)
	if err != nil {
		fmt.Println(fmt.Sprintf("failed to load config file %s", Conf.FilePath), err)
	}
}
