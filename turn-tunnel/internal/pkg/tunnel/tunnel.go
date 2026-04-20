package tunnel

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"

	"golang.org/x/sync/errgroup"
)

type Tunnel struct {
	cfg Config

	wg  *errgroup.Group
	ctx context.Context

	remoteAddr net.Addr
	clientAddr atomic.Value

	relayConn  net.PacketConn
	listenConn net.PacketConn
}

func NewTunnel(pCtx context.Context, cfg Config, relayConn net.PacketConn) *Tunnel {
	wg, ctx := errgroup.WithContext(pCtx)

	return &Tunnel{
		cfg:       cfg,
		wg:        wg,
		ctx:       ctx,
		relayConn: relayConn,
	}
}

func (t *Tunnel) Start() error {
	var err error
	if t.remoteAddr, err = net.ResolveUDPAddr("udp", t.cfg.Remote); err != nil {
		return fmt.Errorf("resolve remote udp addr: %w", err)
	}

	if t.listenConn, err = net.ListenPacket("udp", t.cfg.Listen); err != nil {
		return fmt.Errorf("could not listen on %s: %w", t.cfg.Listen, err)
	}

	t.wg.Go(t.localToRelayLoop)
	t.wg.Go(t.relayToLocalLoop)

	return nil
}

func (t *Tunnel) Wait() error {
	return t.wg.Wait()
}

func (t *Tunnel) Close() error {
	if t.listenConn != nil {
		if err := t.listenConn.Close(); err != nil {
			return fmt.Errorf("could not close listener: %w", err)
		}
	}

	return nil
}

func (t *Tunnel) localToRelayLoop() error {
	buf := make([]byte, t.cfg.MTU)
	for {
		select {
		case <-t.ctx.Done():
			return fmt.Errorf("context canceled: %w", t.ctx.Err())
		default:
		}

		n, addr, err := t.listenConn.ReadFrom(buf)
		if err != nil {
			return fmt.Errorf("could not read from listener: %w", err)
		}

		t.clientAddr.Store(addr)

		_, err = t.relayConn.WriteTo(buf[:n], t.remoteAddr)
		if err != nil {
			return fmt.Errorf("could not write to relay: %w", err)
		}
	}
}

func (t *Tunnel) relayToLocalLoop() error {
	buf := make([]byte, t.cfg.MTU)
	for {
		select {
		case <-t.ctx.Done():
			return fmt.Errorf("context canceled: %w", t.ctx.Err())
		default:
		}

		n, _, err := t.relayConn.ReadFrom(buf)
		if err != nil {
			return fmt.Errorf("could not read from relay: %w", err)
		}

		addrVal := t.clientAddr.Load()
		if addrVal == nil {
			continue
		}

		addr, ok := addrVal.(net.Addr)
		if !ok {
			return fmt.Errorf("invalid address type in store")
		}

		_, err = t.listenConn.WriteTo(buf[:n], addr)
		if err != nil {
			return fmt.Errorf("could not write to relay: %w", err)
		}
	}
}
