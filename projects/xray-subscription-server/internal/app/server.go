package app

import (
	"codeberg.org/kasefuchs/extra/projects/xray-subscription-server/internal/pkg/subscription"
	"github.com/gofiber/fiber/v3"
)

func handleSubscription(sub *subscription.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		data, err := sub.Subscription(c.Context(), c.Params("email"))
		if err != nil {
			return err
		}

		return c.Send(data)
	}
}
