package vpnengine

import (
	"net/netip"
	"strconv"

	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun/netstack"
)

// StartServer menjalankan WireGuard server di userspace via gVisor netstack
func StartServer(privateKey string, listenPort int, localAddress string) {
	// Konversi string IP ke netip.Addr (API terbaru mewajibkan netip)
	addr, err := netip.ParseAddr(localAddress)
	if err != nil {
		return
	}

	// CreateNetTUN membutuhkan []netip.Addr
	tun, tnet, err := netstack.CreateNetTUN(
		[]netip.Addr{addr},
		[]netip.Addr{netip.MustParseAddr("8.8.8.8")},
		1420)
	
	if err != nil {
		return
	}

	// tnet harus digunakan untuk routing atau minimal di-ignore
	_ = tnet 

	// Device baru sekarang butuh 3 argumen: TUN, Bind, dan Logger
	dev := device.NewDevice(
		tun, 
		conn.NewDefaultBind(), 
		device.NewLogger(device.LogLevelError, "wg-go: "),
	)

	// Inisialisasi konfigurasi
	config := "private_key=" + privateKey + "\nlisten_port=" + strconv.Itoa(listenPort)
	dev.IpcSet(config)
	
	// Jalankan interface
	dev.Up()
}
