package scanners

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

func IpToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}

func Uint32ToIP(n uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, n)
	return ip
}

type ScanResult struct {
	IP    uint32
	Ports []int
}

func GetOutIP(target string) (net.IP, error) {
	dst := net.ParseIP(target)
	routes, err := netlink.RouteGet(dst)
	if err != nil || len(routes) == 0 {
		return nil, fmt.Errorf("route not found: %v", err)
	}
	return routes[0].Src, nil
}
