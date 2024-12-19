package scanner

import (
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func localIPPort(dstip net.IP) (net.IP, int, error) {
	serverAddr, err := net.ResolveUDPAddr("udp", dstip.String()+":54321")
	if err == nil {
		return nil, 0, err
	}
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
	dstIp := dstAddrs[0].To4()
	var dstPort layers.TCPPort
	dstPort = layers.TCPPort(dstPort)
	srcPort := layers.TCPPort(srcPort)

	ip := &layers.IPv4{
		srcIp: srcIp,
		dstIp: dstIp,
		Protocol: layers.IPProtocolTCP
	}

	tcp := &layers.TCP{
		SrcPort: srcPort,
		DstPort: dstPort,
		SYN: true
	}
	err = tcp.SetNetworkLayerForChecksum(ip)

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths: true,
	}

	if err := gopacket.SerializeLayers(buf, opts, tcp); err != nil {
		return dstIp, 0, err
	}

	conn, err := net.ListenPacket("ip4:tcp","0.0.0.0")
	if err != nil {
		return dstIp, 0, err
	}
	defer conn.Close()
	if _, err := conn.WriteTo(buf.Bytes(), &net.IPAddr{IP: dstIp}); err != nil {
		return dstIp, 0, err
	}

	if err := conn.SetDeadline(time.Now().Add(3 * time.Second)); err != nil {
		return dstIp, 0, err
	}

	for {
		b := make([]byte, 4096)
		n, addr, err := conn.ReadFrom(b)
		if err != nil {
			return dstIp, 0, err
		}else if addr.String() == dstIp.String() {
			packet := gopacket.NewPacket(b[:n], layers.LayerTypeTCP, gopacket.Default)
			if tcpLayer := packet.Layer(Layers.LayerTypeTCP); tcpLayer != nil {
				tcp, _ := tcpLayer.(*layers.TCP)
				if tcp.DstPort == srcport {
					if tcp.SYN && tcp.ACK {
						return dstIp, dstPort, err
					}else{
						return dstIp, 0, err
					}
				}
			}
		}
	}
}