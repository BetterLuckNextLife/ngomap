package scanners

import (
	"encoding/binary"
	"fmt"
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

func GetOutIP(target, protocol string) (net.IP, error) {
	var conn net.Conn
	var err error

	switch protocol {
	case "tcp":
		conn, err = net.Dial("tcp", target+":22")
	case "udp":
		conn, err = net.Dial("udp", target+":53")
	default:
		return nil, fmt.Errorf("unsupported protocol %s", protocol)
	}
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	switch v := conn.LocalAddr().(type) {
	case *net.TCPAddr:
		return v.IP, nil
	case *net.UDPAddr:
		return v.IP, nil
	default:
		return nil, fmt.Errorf("unsupported address type: %T", v)
	}
}
