package scanner

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"fmt"
	"log"
	"net"

	"github.com/google/gopacket/pcap"
)

// get the local ip and port based on our destination ip
func localIPPort(dstip net.IP) (net.IP, int, error) {
	serverAddr, err := net.ResolveUDPAddr("udp", dstip.String()+":54321")
	if err != nil {
		return nil, 0, err
	}
	// We don't actually connect to anything, but we can determine
	// based on our destination ip what source ip we should use.
	if con, err := net.DialUDP("udp", nil, serverAddr); err == nil {
		if udpaddr, ok := con.LocalAddr().(*net.UDPAddr); ok {
			return udpaddr.IP, udpaddr.Port, nil
		}
	}
	return nil, -1, err
}
func SynScan(dstIp string, dstPort int) (string, int, error) {
	srcIp, srcPort, err := localIPPort(net.ParseIP(dstIp))
	dstAddrs, err := net.LookupIP(dstIp)
	if err != nil {
		return dstIp, 0, err
	}

	dstip := dstAddrs[0].To4()
	var dstport layers.TCPPort
	dstport = layers.TCPPort(dstPort)
	srcport := layers.TCPPort(srcPort)

	// Our IP header... not used, but necessary for TCP checksumming.
	ip := &layers.IPv4{
		SrcIP:    srcIp,
		DstIP:    dstip,
		Protocol: layers.IPProtocolTCP,
		Version:  4,
		TTL:      64,
	}
	// Our TCP header
	tcp := &layers.TCP{
		SrcPort: srcport,
		DstPort: dstport,
		SYN:     true,
	}
	err = tcp.SetNetworkLayerForChecksum(ip)

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}

	if err := gopacket.SerializeLayers(buf, opts, tcp); err != nil {
		return dstIp, 0, err
	}
	handle, err := pcap.OpenLive("lo0", 1600, false, pcap.BlockForever)
	if err != nil {
		log.Fatalf("pcap error: %v", err)
	}
	defer handle.Close()

	if err := handle.WritePacketData(buf.Bytes()); err != nil {
		return dstIp, 0, err
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			tcpResp, _ := tcpLayer.(*layers.TCP)
			if tcpResp.SYN && tcpResp.ACK {
				fmt.Printf("Port %d is open\n", dstPort)
				return dstIp, dstPort, err
			} else if tcpResp.RST {
				fmt.Printf("Port %d is closed\n", dstPort)
				return dstIp, 0, err
			}
		}
	}

	return dstIp, 0, err
}
