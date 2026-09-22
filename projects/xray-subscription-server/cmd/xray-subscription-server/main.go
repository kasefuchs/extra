package main

import (
	"flag"

	"codeberg.org/kasefuchs/extra/projects/xray-subscription-server/internal/app"
	"codeberg.org/kasefuchs/go-kit/config"
	"codeberg.org/kasefuchs/go-kit/log"
	"github.com/knadh/koanf/parsers/yaml"
)

func main() {
	path := flag.String("config", "config.yaml", "Path to the config file")
	flag.Parse()

	if err := config.LoadFile(*path, yaml.Parser()); err != nil {
		log.Fatal().Err(err).Msg("Failed to load config file")
	}

	cfg, err := config.Build[Config]()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to build config")
	}

	if err := log.Init(cfg.Log); err != nil {
		log.Fatal().Err(err).Msg("Failed to init logger")
	}

	app.Run(cfg.Config)
}
