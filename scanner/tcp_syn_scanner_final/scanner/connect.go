package scanner 

import (
	"net"
	"time"
	"fmt"
	"SafeDp/scanner/tcp_syn_scanner_final/vars"
)

func Connect(ip string, port int)(string, int, error){
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ip, port), time.Duration(vars.Timeout)*time.Second)
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()

	return ip, port, err
}