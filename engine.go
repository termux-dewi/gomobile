package wgserver

import (
    "net"
    "gvisor.dev/gvisor/pkg/tcpip/stack"
    "golang.zx2c4.com/wireguard/device"
    "golang.zx2c4.com/wireguard/tun/netstack"
)

// StartServer menjalankan WireGuard server di userspace
func StartServer(privateKey string, listenPort int, localAddress string) {
    tun, tnet, err := netstack.CreateNetTUN(
        []net.IP{net.ParseIP(localAddress)},
        []net.IP{net.ParseIP("8.8.8.8")},
        1420)
    
    if err != nil {
        return
    }

    dev := device.NewDevice(tun, device.NewLogger(device.LogLevelError, "wg-go: "))
    config := "private_key=" + privateKey + "\nlisten_port=" + strconv.Itoa(listenPort)
    dev.IpcSet(config)
    dev.Up()
}
