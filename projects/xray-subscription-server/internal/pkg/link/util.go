package link

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"net/url"
	"strings"

	"github.com/xtls/xray-core/app/proxyman"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/core"
	vless "github.com/xtls/xray-core/proxy/vless/inbound"
	"github.com/xtls/xray-core/transport/internet"
	"github.com/xtls/xray-core/transport/internet/grpc"
	"github.com/xtls/xray-core/transport/internet/httpupgrade"
	"github.com/xtls/xray-core/transport/internet/reality"
	"github.com/xtls/xray-core/transport/internet/splithttp"
	"github.com/xtls/xray-core/transport/internet/tls"
	"github.com/xtls/xray-core/transport/internet/websocket"
	"golang.org/x/crypto/curve25519"
)

var linkTypes = map[string]string{
	"splithttp": "xhttp",
}

func findUser(inbound *core.InboundHandlerConfig, email string) (*protocol.User, error) {
	msg, err := inbound.ProxySettings.GetInstance()
	if err != nil {
		return nil, fmt.Errorf("decode proxy settings: %w", err)
	}

	var users []*protocol.User
	switch cfg := msg.(type) {
	case *vless.Config:
		users = cfg.Clients
	default:
		return nil, fmt.Errorf("unsupported protocol: %T", msg)
	}

	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, fmt.Errorf("user not found")
}

func receiverSettings(inbound *core.InboundHandlerConfig) (*proxyman.ReceiverConfig, error) {
	msg, err := inbound.ReceiverSettings.GetInstance()
	if err != nil {
		return nil, fmt.Errorf("decode receiver settings: %w", err)
	}

	receiver, ok := msg.(*proxyman.ReceiverConfig)
	if !ok {
		return nil, fmt.Errorf("unexpected receiver settings type %T", msg)
	}

	return receiver, nil
}

func addParam(q *url.Values, key, value string) {
	if value != "" {
		q.Set(key, value)
	}
}

func addType(q *url.Values, stream *internet.StreamConfig) {
	t := stream.GetProtocolName()
	if mapped, ok := linkTypes[t]; ok {
		t = mapped
	}

	addParam(q, "type", t)
}

func addHostPath(q *url.Values, host, path string) {
	addParam(q, "host", host)
	addParam(q, "path", path)
}

func shortID(uuid string, ids [][]byte) string {
	h := fnv.New32a()
	h.Write([]byte(uuid))
	return hex.EncodeToString(ids[h.Sum32()%uint32(len(ids))])
}

func addSecurity(q *url.Values, stream *internet.StreamConfig, uuid string) error {
	secs := stream.GetSecuritySettings()
	if len(secs) == 0 {
		return nil
	}

	msg, err := secs[0].GetInstance()
	if err != nil {
		return err
	}

	switch s := msg.(type) {
	case *tls.Config:
		q.Set("security", "tls")
		addParam(q, "sni", s.ServerName)

		if len(s.NextProtocol) > 0 {
			q.Set("alpn", strings.Join(s.NextProtocol, ","))
		}

	case *reality.Config:
		q.Set("security", "reality")
		q.Set("sni", s.ServerNames[0])

		pbk, err := curve25519.X25519(s.PrivateKey, curve25519.Basepoint)
		if err != nil {
			return err
		}

		q.Set("pbk", base64.RawURLEncoding.EncodeToString(pbk))
		q.Set("sid", shortID(uuid, s.ShortIds))
	}

	return nil
}

func addTransport(q *url.Values, stream *internet.StreamConfig) error {
	trans := stream.GetTransportSettings()
	if len(trans) == 0 {
		return nil
	}

	msg, err := trans[0].GetSettings().GetInstance()
	if err != nil {
		return err
	}

	switch s := msg.(type) {
	case *websocket.Config:
		addHostPath(q, s.Host, s.Path)
	case *httpupgrade.Config:
		addHostPath(q, s.Host, s.Path)
	case *splithttp.Config:
		addHostPath(q, s.Host, s.Path)
	case *grpc.Config:
		addParam(q, "authority", s.Authority)
		addParam(q, "serviceName", s.ServiceName)
	}

	return nil
}
