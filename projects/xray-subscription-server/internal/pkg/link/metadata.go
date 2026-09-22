package link

import (
	"net"
	"strconv"

	"github.com/xtls/xray-core/app/proxyman"
)

type Metadata struct {
	Remark  string `koanf:"remark"`
	Address string `koanf:"address"`
	Port    uint16 `koanf:"port"`
}

func (m Metadata) Host(receiver *proxyman.ReceiverConfig) string {
	port := m.Port
	if port == 0 {
		if ranges := receiver.PortList.GetRange(); len(ranges) > 0 {
			port = uint16(ranges[0].GetFrom())
		}
	}

	return net.JoinHostPort(m.Address, strconv.Itoa(int(port)))
}
