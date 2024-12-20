package scanner

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"fmt"
	"log"
	"net"
	"time"

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

func SynScan2(dstIp string, dstPort int) (string, int, error) {
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

	conn, err := net.ListenPacket("ip4:tcp", "0.0.0.0")
	if err != nil {
		return dstIp, 0, err
	}
	defer conn.Close()

	if _, err := conn.WriteTo(buf.Bytes(), &net.IPAddr{IP: dstip}); err != nil {
		return dstIp, 0, err
	}

	// Set deadline so we don't wait forever.
	if err := conn.SetDeadline(time.Now().Add(4 * time.Second)); err != nil {
		return dstIp, 0, err
	}

	for {
		b := make([]byte, 8192)
		n, addr, err := conn.ReadFrom(b)
		if err != nil {
			fmt.Printf("Error reading from connection: %v\n", err)
			return dstIp, 0, err
		} else if addr.String() == dstip.String() {
			// Decode a packet
			packet := gopacket.NewPacket(b[:n], layers.LayerTypeTCP, gopacket.Default)
			// Get the TCP layer from this packet
			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
				tcp, _ := tcpLayer.(*layers.TCP)

				if tcp.DstPort == srcport {
					if tcp.SYN && tcp.ACK {
						// log.Printf("%v:%d is OPEN\n", dstIp, dstport)
						return dstIp, dstPort, err
					} else {
						return dstIp, 0, err
					}
				}
			}
		}
	}
}
