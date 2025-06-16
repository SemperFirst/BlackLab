package proxy
import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"SafeDp/scanner/proxy_scanner/models"
	"SafeDp/scanner/proxy_scanner/util"
)

var (
	DebugMode = false
	ScanNum = 100
	IpList = "iplist.txt"
	Timeout = 5
)
type CheckProxyFunc func(ip string, port int, protocol string) (isProxy bool, proxyInfo models.ProxyInfo, err error)

var (
	httpProxyFunc CheckProxyFunc = CheckHttpProxy
	sockProxyFunc CheckProxyFunc = CheckSockProxy
)

func CheckProxy(proxyAddr []util.ProxyAddr) {
	var wg sync.WaitGroup
	wg.Add(len(proxyAddr) * (len(HttpProxyProtocol) + len(SockProxyProtocol)))
	for _, addr := range proxyAddr {
		for _, proto := range HttpProxyProtocol {
			go func(ip string, port int, protocol string) {
				defer wg.Done()
				 _= models.SaveProxies(httpProxyFunc(ip, port, protocol))
			}(addr.IP, addr.Port, proto)
		}
		for proto := range SockProxyProtocol {
			go func(ip string, port int, protocol string) {
				defer wg.Done()
				_= models.SaveProxies(sockProxyFunc(ip, port, protocol))
			}(addr.IP, addr.Port, proto)
		}
	}
	wg.Wait()
}