package main

import (
	"codeberg.org/kasefuchs/extra/projects/xray-subscription-server/internal/app"
	"codeberg.org/kasefuchs/go-kit/log"
)

type Config struct {
	app.Config `koanf:",squash"`
	Log        log.Config `koanf:"log"`
}
