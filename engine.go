package vpnbridge

import (
	"github.com/xjasonlyu/tun2socks/v2/engine"
)

// StartEngine menjalankan tun2socks dengan konfigurasi dasar yang stabil
func StartEngine(socksAddr string, tunName string, mtu int) {
	config := &engine.Key{
		Proxy:      "socks5://" + socksAddr,
		Device:     tunName,
		LogLevel:   "info",
		MTU:        mtu,
		UDPTimeout: 60000,
	}

	// Pada versi terbaru tun2socks v2, Insert biasanya tidak mengembalikan error langsung
	// Kita panggil Insert, lalu jalankan Start
	engine.Insert(config)
	engine.Start()
}

func StopEngine() {
	engine.Stop()
}
