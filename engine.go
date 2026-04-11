package vpnengine

import (
	"log"
	"net"
	"sync"

	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun"
	"golang.zx2c4.com/wireguard/tun/netstack"
	"google.golang.org/grpc"
)

var (
	wgDevice  *device.Device
	isStarted bool
	mu        sync.Mutex
)

func GenerateKey() string {
	// Implementasi dummy untuk testing
	return "PRIVATE_KEY_SAMPEL"
}

func Stop() {
	mu.Lock()
	defer mu.Unlock()
	if wgDevice != nil {
		wgDevice.Close()
		wgDevice = nil
	}
	isStarted = false
	log.Println("VPN Stopped")
}

// StartClient - Mode Client
func StartClient(endpoint string, privKey string, pubKey string, localIP string) string {
	mu.Lock()
	defer mu.Unlock()
	if isStarted {
		return "ALREADY_RUNNING"
	}

	// Perbaikan: netstack.CreateNetStack (S kapital di Stack)
	tunDev, _, err := netstack.CreateNetStack(
		[]net.IP{net.ParseIP(localIP)},
		[]net.IP{net.ParseIP("8.8.8.8")},
		1420,
	)
	if err != nil {
		return "TUN_ERROR: " + err.Error()
	}

	// Perbaikan: NewDevice butuh 3 argumen (tun, bind, logger)
	logger := device.NewLogger(device.LogLevelError, "vpn:")
	wgDevice = device.NewDevice(tunDev, conn.NewDefaultBind(), logger)

	isStarted = true
	return "CLIENT_STARTED"
}

// StartServer - Mode Server
func StartServer(port string, privKey string, localIP string) string {
	mu.Lock()
	defer mu.Unlock()
	if isStarted {
		return "ALREADY_RUNNING"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return "BIND_ERROR: " + err.Error()
	}

	s := grpc.NewServer()
	// Di sini nanti register service gRPC kamu

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Printf("Server failed: %v", err)
		}
	}()

	isStarted = true
	return "SERVER_STARTED"
}
