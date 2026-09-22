package app

import (
	"codeberg.org/kasefuchs/extra/projects/xray-subscription-server/internal/pkg/link"
	"codeberg.org/kasefuchs/go-kit/grpc/client"
)

type Config struct {
	Listen   string                     `koanf:"listen"`
	Client   client.Config              `koanf:"client"`
	Inbounds map[string][]link.Metadata `koanf:"inbounds"`
}
