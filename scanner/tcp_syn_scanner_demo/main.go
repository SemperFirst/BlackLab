package main

import (
	"SafeDp/scanner/tcp_syn_scanner_demo/scan"
	"SafeDp/scanner/tcp_syn_scanner_demo/util"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 3 {
		util.CheckRoot()

		ipList := os.Args[1]
		portList := os.Args[2]
		ips, err := util.GetIpList(ipList)
		ports, err := util.GetPorts(portList)
		_ = err
		for _, ip := range ips {
			for _, port := range ports {
				ip1, port1, err1 := scan.SynScan(ip.String(), port)
				if err1 == nil && port1 > 0 {
					fmt.Printf("%v:%v is open\n", ip1, port1)
				} else {
					fmt.Printf("Scan failed for %v:%v - %v\n", ip, port, err1)
				}
			}
		}
	} else {
		fmt.Printf("%v iplist port\n", os.Args[0])
	}
}
