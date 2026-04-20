package client

import (
	"errors"
	"fmt"
	"net"

	"codeberg.org/kasefuchs/go-kit/log"
	"github.com/pion/turn/v5"
)

type Client struct {
	cfg Config

	client *turn.Client

	turnConn  *net.UDPConn
	relayConn net.PacketConn
}

func NewClient(cfg Config) *Client {
	return &Client{cfg: cfg}
}

func (c *Client) allocate() error {
	turnAddr, err := net.ResolveUDPAddr("udp", c.cfg.Turn)
	if err != nil {
		return fmt.Errorf("resolve turn addr: %w", err)
	}

	if c.turnConn, err = net.DialUDP("udp", nil, turnAddr); err != nil {
		return fmt.Errorf("dial turn addr: %w", err)
	}

	clientCfg := &turn.ClientConfig{
		Conn:           &packetUDPConn{c.turnConn},
		Username:       c.cfg.Username,
		Password:       c.cfg.Password,
		STUNServerAddr: c.cfg.Stun,
		TURNServerAddr: c.cfg.Turn,
		LoggerFactory:  &pionLoggerFactory{},
	}

	if c.client, err = turn.NewClient(clientCfg); err != nil {
		_ = c.Close()
		return fmt.Errorf("could not create turn client: %w", err)
	}

	if err = c.client.Listen(); err != nil {
		_ = c.Close()
		return fmt.Errorf("could not listen: %w", err)
	}

	if c.relayConn, err = c.client.Allocate(); err != nil {
		_ = c.Close()
		return fmt.Errorf("could not allocate relay connection: %w", err)
	}

	log.Info().Str("relayed-address", c.relayConn.LocalAddr().String()).Msg("Allocation successful")

	return nil
}

func (c *Client) RelayConn() (net.PacketConn, error) {
	if c.relayConn == nil {
		if err := c.allocate(); err != nil {
			return nil, fmt.Errorf("could not allocate relay connection: %w", err)
		}
	}

	return c.relayConn, nil
}

func (c *Client) Close() error {
	var errs []error

	if c.relayConn != nil {
		errs = append(errs, c.relayConn.Close())
	}

	if c.client != nil {
		c.client.Close()
	}

	if c.turnConn != nil {
		errs = append(errs, c.turnConn.Close())
	}

	return errors.Join(errs...)
}
