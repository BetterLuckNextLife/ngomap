package scanners

import (
	"encoding/binary"
	"net"
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

func GetOutIP() (net.IP, error) {
	conn, err := net.Dial("tcp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}
