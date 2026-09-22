package link

import (
	"fmt"
	"net/url"

	"github.com/xtls/xray-core/app/proxyman"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/proxy/vless"
)

func Build(inbound *core.InboundHandlerConfig, metadata []Metadata, email string) ([]*url.URL, error) {
	user, err := findUser(inbound, email)
	if err != nil {
		return nil, err
	}

	receiver, err := receiverSettings(inbound)
	if err != nil {
		return nil, err
	}

	account, err := user.Account.GetInstance()
	if err != nil {
		return nil, err
	}

	uris := make([]*url.URL, 0, len(metadata))
	for _, meta := range metadata {
		uri, err := linkURI(account, receiver, meta)
		if err != nil {
			return nil, err
		}

		uris = append(uris, uri)
	}

	return uris, nil
}

func linkURI(account any, receiver *proxyman.ReceiverConfig, m Metadata) (*url.URL, error) {
	switch acc := account.(type) {
	case *vless.Account:
		return vlessURI(acc, receiver, m)
	default:
		return nil, fmt.Errorf("%w: %T", ErrUnsupportedProtocol, acc)
	}
}
