package scanners

import (
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func BuildSYN(srcIP, dstIP string, srcPort, dstPort uint16) ([]byte, error) {
	ip := &layers.IPv4{
		SrcIP:    net.ParseIP(srcIP),
		DstIP:    net.ParseIP(dstIP),
		Protocol: layers.IPProtocolTCP,
	}

	tcp := &layers.TCP{
		SrcPort: layers.TCPPort(srcPort),
		DstPort: layers.TCPPort(dstPort),
		SYN:     true,
	}

	if err := tcp.SetNetworkLayerForChecksum(ip); err != nil {
		return nil, err
	}

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
	if err := gopacket.SerializeLayers(buf, opts, ip, tcp); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func SendRawPacket(dstIP string, data []byte) error {
	conn, err := net.Dial("tcp", dstIP)
	if err != nil {
		return err
	}

	defer conn.Close()

	n, err := conn.Write(data)
	if err != nil || n != len(data) {
		return err
	}
	return nil
}
