package main

import(
	"fmt"
	"os"
	"sec-dev-in-action-src/scanner/tcp_syn_scanner_demo/scan"
	"sec-dev-in-action-src/scanner/tcp_syn_scanner_demo/util"
)

func main() {
	if len(os.Args) == 3 {
		util.CheckRoot()

		ipList := os.Args[1]
		portList := os.Args[2]
		ipList, err := util.GetIpList(ipList)
		ports, err := util.GetPorts(portList)
		_ = err 
		for _, ip := range ipList {
			for _, port := range ports {
				ip1, port1, err1 := scan.SynScan(ip.String(), port)
				if err1 == nil && port1 > 0 {
					fmt.Println("%v:%v is open\n", ip1, port1)
				}
			}
		}
	}else{
		fmt.Println("Usage: %v <ip> <port>\n", os.Args[0])
	}
}