package link

import (
	"fmt"
	"net/url"

	"github.com/xtls/xray-core/app/proxyman"
	"github.com/xtls/xray-core/proxy/vless"
)

func vlessURI(account *vless.Account, receiver *proxyman.ReceiverConfig, m Metadata) (*url.URL, error) {
	uri := &url.URL{
		Scheme:   "vless",
		Host:     m.Host(receiver),
		User:     url.User(account.Id),
		Fragment: m.Remark,
	}

	q := uri.Query()

	addType(&q, receiver.StreamSettings)
	addParam(&q, "flow", account.Flow)

	if err := addSecurity(&q, receiver.StreamSettings, account.Id); err != nil {
		return nil, fmt.Errorf("security settings: %w", err)
	}

	if err := addTransport(&q, receiver.StreamSettings); err != nil {
		return nil, fmt.Errorf("transport settings: %w", err)
	}

	uri.RawQuery = q.Encode()
	return uri, nil
}
