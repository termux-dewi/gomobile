package vpnengine

import (
	"crypto/rand"
	"encoding/base64"
	"net/netip"

	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun/netstack"
)

var wgDev *device.Device

func GenerateKey() string {
	k := make([]byte, 32)
	rand.Read(k)
	return base64.StdEncoding.EncodeToString(k)
}

func StartClient(domain string, privKey string, peerPub string, clientIP string) string {
	// Konversi string IP ke tipe netip.Addr (Standard terbaru WireGuard)
	addr, err := netip.ParseAddr(clientIP)
	if err != nil {
		return "ERR_IP: " + err.Error()
	}
	dns := netip.MustParseAddr("1.1.1.1")

	// Perbaikan netstack.CreateNetTUN
	tun, _, err := netstack.CreateNetTUN(
		[]netip.Addr{addr},
		[]netip.Addr{dns},
		1280,
	)
	if err != nil {
		return "FAIL_TUN: " + err.Error()
	}

	// Perbaikan device.NewDevice (Sekarang butuh conn.NewDefaultBind)
	wgDev = device.NewDevice(tun, conn.NewDefaultBind(), device.NewLogger(device.LogLevelSilent, ""))
	
	config := "private_key=" + privKey + "\n" +
		"public_key=" + peerPub + "\n" +
		"endpoint=127.0.0.1:10000\n" +
		"allowed_ip=0.0.0.0/0\n"

	wgDev.IpcSet(config)
	wgDev.Up()
	return "STARTED"
}

func StartServer(privKey string, port string) string {
	addr := netip.MustParseAddr("10.0.0.1")
	dns := netip.MustParseAddr("1.1.1.1")

	tun, _, _ := netstack.CreateNetTUN(
		[]netip.Addr{addr},
		[]netip.Addr{dns},
		1280,
	)
	
	// Perbaikan device.NewDevice
	wgDev = device.NewDevice(tun, conn.NewDefaultBind(), device.NewLogger(device.LogLevelSilent, ""))
	
	wgDev.IpcSet("private_key=" + privKey + "\nlisten_port=" + port)
	wgDev.Up()
	return "SERVER_UP"
}

func Stop() {
	if wgDev != nil {
		wgDev.Close()
	}
}
