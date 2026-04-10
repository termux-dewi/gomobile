package wgengine

import (
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun"
)

// StartServer menjalankan stack Wireguard di dalam VpnService
func StartServer(fd int) {
	// 1. Inisialisasi TUN dari File Descriptor Android
	tunDev, _, err := tun.CreateTUNFromFile(uintptr(fd), 1420)
	if err != nil {
		return
	}

	// 2. Buat Device Wireguard
	logger := device.NewLogger(device.LogLevelVerbose, "[WG-MOBILE]")
	dev := device.NewDevice(tunDev, conn.NewDefaultBind(), logger)

	// 3. Konfigurasi Server (Ganti PRIVATE_KEY dengan milik Anda)
	// Ganti 'uE9p...' dengan hasil generate Anda
	config := `private_key=4BHkGskA/cjExGJMwE8PX/gq27hKqU5DZZAyDFbt8FY=
listen_port=51820
replace_peers=true
public_key=GwaS9tFFVnmF01ow6uWjQ/QwIs5TC6I/gL5K3wF+uGE=
allowed_ip=10.0.0.2/32`

	dev.IpcSet(config)
	dev.Up()
}
