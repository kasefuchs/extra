package subscription

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"codeberg.org/kasefuchs/extra/projects/xray-subscription-server/internal/pkg/link"
	"codeberg.org/kasefuchs/go-kit/log"
	"github.com/xtls/xray-core/app/proxyman/command"
	"github.com/xtls/xray-core/core"
)

type Service struct {
	client   command.HandlerServiceClient
	metadata map[string][]link.Metadata
}

func NewService(client command.HandlerServiceClient, metadata map[string][]link.Metadata) *Service {
	return &Service{
		client,
		metadata,
	}
}

func (s *Service) Subscription(ctx context.Context, email string) ([]byte, error) {
	response, err := s.client.ListInbounds(ctx, &command.ListInboundsRequest{})
	if err != nil {
		return nil, fmt.Errorf("list inbounds: %w", err)
	}

	uris, err := s.buildSubscription(response.Inbounds, email)
	if err != nil {
		return nil, err
	}

	return render(uris), nil
}

func (s *Service) buildSubscription(inbounds []*core.InboundHandlerConfig, email string) ([]*url.URL, error) {
	var uris []*url.URL

	for _, inbound := range inbounds {
		tag := inbound.GetTag()
		meta, ok := s.metadata[tag]
		if !ok {
			continue
		}

		built, err := link.Build(inbound, meta, email)
		if err != nil {
			if !errors.Is(err, link.ErrUserNotFound) {
				log.Warn().Err(err).Str("tag", tag).Msg("failed to build inbound links, skipping")
			}
			
			continue
		}

		uris = append(uris, built...)
	}

	return uris, nil
}
