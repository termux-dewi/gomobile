package wgengine

import (
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun"
)

func StartServer(fd int, privateKey string, clientPubKey string) {
	tunDev, _, err := tun.CreateTUNFromFile(uintptr(fd), 1420)
	if err != nil {
		return
	}

	logger := device.NewLogger(device.LogLevelVerbose, "[WG-MOBILE]")
	dev := device.NewDevice(tunDev, conn.NewDefaultBind(), logger)

	// Konfigurasi dinamis
	config := "private_key=" + privateKey + "\n" +
		"listen_port=51820\n" +
		"replace_peers=true\n" +
		"public_key=" + clientPubKey + "\n" +
		"allowed_ip=10.0.0.2/32"

	dev.IpcSet(config)
	dev.Up()
}
