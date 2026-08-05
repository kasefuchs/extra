package main

import (
	"context"
	"flag"

	"codeberg.org/kasefuchs/extra/projects/turn-tunnel/internal/pkg/client"
	"codeberg.org/kasefuchs/extra/projects/turn-tunnel/internal/pkg/tunnel"
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

	c := client.NewClient(cfg.Client)
	defer func(c *client.Client) {
		if err := c.Close(); err != nil {
			log.Fatal().Err(err).Msg("Failed to close client")
		}
	}(c)

	conn, err := c.RelayConn()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get relay connection")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tun := tunnel.NewTunnel(ctx, cfg.Tunnel, conn)
	defer func(tun *tunnel.Tunnel) {
		if err := tun.Close(); err != nil {
			log.Fatal().Err(err).Msg("Failed to close tunnel")
		}
	}(tun)

	if err := tun.Start(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start tunnel")
	}

	if err := tun.Wait(); err != nil {
		log.Fatal().Err(err).Msg("Failed to wait tunnel")
	}
}
