package main 

import (
	"fmt"
	"os"
	"runtime"

	"SafeDp/scanner/tcp_syn_scanner/util"
	"SafeDp/scanner/tcp_syn_scanner/scanner"
)

func main() {
	if len(os.Args) == 3 {
		util.CheckRoot()
		ipList := os.Args[1]
		portList := os.Args[2]
		ips, err := util.GetIpList(ipList)
		ports, err := util.GetPorts(portList)
		_ = err
		tasks, _ := scanner.GenerateTask(ips, ports)
		scanner.RunTask(tasks)
		scanner.PrintResult()
	} else {
		fmt.Printf("%v iplist port\n", os.Args[0])
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}