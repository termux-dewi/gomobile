package vpnengine

import (
	"crypto/rand"
	"encoding/base64"
	"net"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun/netstack"
)

var wgDev *device.Device

// GenerateKey untuk tombol refresh di UI Sketchware
func GenerateKey() string {
	k := make([]byte, 32)
	rand.Read(k)
	return base64.StdEncoding.EncodeToString(k)
}

// StartClient menghubungkan ke Cloudflare Tunnel via gRPC
func StartClient(domain string, privKey string, peerPub string, clientIP string) string {
	tun, _, err := netstack.CreateNetTUN(
		[]net.IP{net.ParseIP(clientIP)},
		[]net.IP{net.ParseIP("1.1.1.1")},
		1280,
	)
	if err != nil {
		return "ERROR: " + err.Error()
	}

	wgDev = device.NewDevice(tun, device.NewLogger(device.LogLevelSilent, ""))
	// Endpoint diarahkan ke bridge internal
	config := "private_key=" + privKey + "\n" +
		"public_key=" + peerPub + "\n" +
		"endpoint=127.0.0.1:10000\n" +
		"allowed_ip=0.0.0.0/0\n"

	wgDev.IpcSet(config)
	wgDev.Up()
	return "CONNECTED"
}

// StartServer mode standby menerima traffic gRPC
func StartServer(privKey string, port string) string {
	tun, _, _ := netstack.CreateNetTUN(
		[]net.IP{net.ParseIP("10.0.0.1")},
		[]net.IP{net.ParseIP("1.1.1.1")},
		1280,
	)
	wgDev = device.NewDevice(tun, device.NewLogger(device.LogLevelSilent, ""))
	wgDev.IpcSet("private_key=" + privKey + "\nlisten_port=" + port)
	wgDev.Up()
	return "SERVER_READY"
}

func Stop() {
	if wgDev != nil {
		wgDev.Close()
	}
}
