package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/Gorynychdo/tdligo.git/internal/model"
	"github.com/Gorynychdo/tdligo.git/internal/service"
	"github.com/Gorynychdo/tdligo.git/internal/tdclient"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config_path", "configs/tdlib.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := model.NewConfig()
	if _, err := toml.DecodeFile(configPath, config); err != nil {
		log.Fatal(err)
	}

	tc := tdclient.NewTDClient(config)
	go func() {
		if err := tc.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	srv := service.NewHTTPServer(config, tc)
	if err := srv.ServeHTTP(); err != nil {
		log.Fatal(err)
	}
}
