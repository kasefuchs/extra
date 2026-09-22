package app

import (
	"codeberg.org/kasefuchs/extra/projects/xray-subscription-server/internal/pkg/subscription"
	"codeberg.org/kasefuchs/go-kit/grpc/client"
	"codeberg.org/kasefuchs/go-kit/log"
	fiberzerolog "github.com/gofiber/contrib/v3/zerolog"
	"github.com/gofiber/fiber/v3"
	"github.com/xtls/xray-core/app/proxyman/command"
)

func Run(cfg Config) {
	srv := fiber.New()

	srv.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: log.Logger(),
	}))

	conn := client.MustNew(cfg.Client)
	defer conn.Close()

	handler := command.NewHandlerServiceClient(conn)
	sub := subscription.NewService(handler, cfg.Inbounds)
	srv.Get("/:email", handleSubscription(sub))

	log.Info().Str("address", cfg.Listen).Msg("server started")
	if err := srv.Listen(cfg.Listen, fiber.ListenConfig{
		DisableStartupMessage: true,
	}); err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}
}
