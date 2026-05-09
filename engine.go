package vpnbridge

import (
	"github.com/xjasonlyu/tun2socks/v2/engine"
)

func StartEngine(socksAddr string, tunName string, mtu int) error {
	config := &engine.Key{
		Proxy:                "socks5://" + socksAddr,
		Device:               tunName,
		LogLevel:             "info",
		MTU:                  mtu,
		UDPTimeout:           60000,
		Sniffing:             true,
		AllowSecondaryRoute:  true,
	}

	return engine.Insert(config)
}

func StopEngine() {
	engine.Stop()
}
