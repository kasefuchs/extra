package client

import "net"

var _ net.PacketConn = (*packetUDPConn)(nil)

type packetUDPConn struct {
	*net.UDPConn
}

func (c *packetUDPConn) WriteTo(p []byte, _ net.Addr) (int, error) {
	return c.Write(p)
}
