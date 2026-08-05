package main

import (
	"codeberg.org/kasefuchs/extra/projects/turn-tunnel/internal/pkg/client"
	"codeberg.org/kasefuchs/extra/projects/turn-tunnel/internal/pkg/tunnel"
	"codeberg.org/kasefuchs/go-kit/log"
)

type Config struct {
	Log    log.Config    `koanf:"log"`
	Client client.Config `koanf:"client"`
	Tunnel tunnel.Config `koanf:"tunnel"`
}
