package proxy

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"SafeDp/scanner/proxy_scanner/models"
	"SafeDp/scanner/proxy_scanner/util"
)

var (
	    HttpProxyProtocol = []string{"http", "https"}
		WebUrl            = "http://email.163.com"
)

func CheckHttpProxy(ip string, port int, protocol string) (isProxy bool, proxyInfo models.ProxyInfo, err error) {
	proxyInfo.Addr = ip
	proxyInfo.Port = port
	proxyInfo.Protocol = protocol

	rawProxyUrl := fmt.Sprintf("%v://%v:%v", protocol, ip, port)
	proxyUrl, err := url.Parse(rawProxyUrl)
	if err != nil {
		return
	}

	Transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	client := &http.Client{Transport: Transport, Timeout: time.Duration(Timeout) * time.Second}

	resp, err := client.Get(WebUrl)
	if err != nil {
		return false, proxyInfo, err
	}

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			return false, proxyInfo, err
		}

		if strings.Contains(string(body), "<title>网易邮箱") {
			isProxy = true
		}
	}

	util.Log.Debugf("Checking proxy: %v, isProxy: %v", rawProxyUrl, isProxy)
	return
}