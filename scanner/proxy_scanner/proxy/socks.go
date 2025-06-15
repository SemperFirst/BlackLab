package proxy
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"h12.io/socks"
	"SafeDP/scanner/proxy_scanner/models"
	"SafeDP/scanner/proxy_scanner/util"
)
var (
	SockProxyProtocol = map[string]int{
		"socks4": socks.SOCKS4,
		"socks4A": socks.SOCKS4A,
		"socks5": socks.SOCKS5,
	}
)

func CheckSocksProxy(ip string, port int, protocol string) (isProxy bool, proxyInfo models.ProxyInfo, err error) {
	proxyInfo.Addr = ip
	proxyInfo.Port = port
	proxyInfo.Protocol = protocol

	proxy := fmt.Sprintf("%v://%v:%v", protocol, ip, port)
	dialSocksProxy := socks.DialSocksProxy(SockProxyProtocol[protocol], proxy)
	tr := &http.Transport{Dial: dialSocksProxy}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(5) * time.Second,
	}
	resp, err := httpClient.Get(WebUrl)
	if err != nil {
		return false, proxyInfo, err
	}

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, proxyInfo, err
		}
		if strings.Contains(string(body), "<title>网易邮箱") {
			isProxy = true
		}
	}
	util.Log.Debugf("Checking proxy: %v, isProxy: %v", proxy, isProxy)
	return isProxy, proxyInfo, nil
}