package tunnel

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"codeberg.org/kasefuchs/go-kit/log"
)

type Tunnel struct {
	cfg Config

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	remoteAddr net.Addr
	clientAddr atomic.Value

	relayConn  net.PacketConn
	listenConn net.PacketConn
}

func NewTunnel(pCtx context.Context, cfg Config, relayConn net.PacketConn) *Tunnel {
	ctx, cancel := context.WithCancel(pCtx)

	return &Tunnel{
		cfg:       cfg,
		ctx:       ctx,
		cancel:    cancel,
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

	context.AfterFunc(t.ctx, func() {
		_ = t.listenConn.Close()
	})

	t.wg.Go(t.localToRelayLoop)
	t.wg.Go(t.relayToLocalLoop)

	return nil
}

func (t *Tunnel) Wait() {
	t.wg.Wait()
}

func (t *Tunnel) Close() error {
	t.cancel()
	if t.listenConn != nil {
		if err := t.listenConn.Close(); err != nil {
			return fmt.Errorf("could not close listener: %w", err)
		}
	}

	return nil
}

func (t *Tunnel) localToRelayLoop() {
	defer t.cancel()

	buf := make([]byte, t.cfg.MTU)
	for {
		if t.ctx.Err() != nil {
			return
		}

		n, addr, err := t.listenConn.ReadFrom(buf)
		if err != nil {
			return
		}

		t.clientAddr.Store(addr)

		_, err = t.relayConn.WriteTo(buf[:n], t.remoteAddr)
		if err != nil {
			return
		}
	}
}

func (t *Tunnel) relayToLocalLoop() {
	defer t.cancel()

	buf := make([]byte, t.cfg.MTU)
	for {
		if t.ctx.Err() != nil {
			return
		}

		n, _, err := t.relayConn.ReadFrom(buf)
		if err != nil {
			return
		}

		addrVal := t.clientAddr.Load()
		if addrVal == nil {
			continue
		}

		addr, ok := addrVal.(net.Addr)
		if !ok {
			log.Warn().Msg("Invalid address type in store")
			return
		}

		_, err = t.listenConn.WriteTo(buf[:n], addr)
		if err != nil {
			return
		}
	}
}
