func writeARP(handle *pcap.Handle, iFace *net.Interface, localIp, ip net.IP) error{
	eth := layers.Ethernet{
		SrcMAC:	   iFace.HardwareAddr,
		DstMAC:    net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}
	arp := layers.ARP{
		AddrType:	   layers.LinkTypeEthernet,
		Protocol:	   layers.EthernetTypeIPv4,
		HwAddressSize: 6,
		ProtAddressSize: 4,
		Operation:	   layers.ARPRequest,
		SourceHwAddress: []byte(iFace.HardwareAddr),
		SourceProtAddress []byte(localIp),
		DstHwAddress: []byte{0,0,0,0,0,0}
	}
}