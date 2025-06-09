package util

import (
	"bufio"
	"go/scanner"
	"os"
	"strconv"
	"strings"

	"safedp/scanner/password_crack/logger"
	"safedp/scanner/password_crack/models"
	"safedp/scanner/password_crack/vars"
)

func ReadIpList(fileName string) (ipList []models.IpAddr) {
	ipListFile, err := os.Open(fileName)
	if err != nil {
		logger.log.Fatalf("Open ip List file err, %s", err)
	}
	defer func() {
		if ipListFile != nil {
			_ = ipListFile.Close()
		}
	}()
	scanner := bufio.NewScanner(ipListFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		ipPort := strings.TrimSpace(line)
		t := strings.Split(ipPort, ":")
		ip := t[0]
		portProtocol := t[1]
		tmpPort := strings.Split(portProtocol, "|")
		if len(tmpPort) == 2{
			port, _ := strconv.Atoi(tmpPort[0])
			protocol := strings.ToUpper(tmpPort[1])
			if vars.SupportProtocols[protocol] {
				addr := models.IpAddr{Ip: ip, Port: port, Protocol: protocol}
				ipList = append(ipList, addr)
			} else {
				logger.log.Infof("Protocol %s not supported, skip", protocol)
			}
		} else {
			port, err := strconv.Atoi(tmpPort[0])
			if err == nil {
				protocol,ok := vars.PortNames[port]
				if ok && vars.SupportProtocols[protocol] {
					addr := models.IpAddr{Ip: ip, Port: port, Protocol: protocol}
					ipList = append(ipList, addr)
				} else {
					logger.log.Infof("Port %d not supported, skip", port)
				}
			}
		}
		return ipList
	}
}

func ReadUserDict(userDict string)(users []string, err error){
	file, err := os.Open(userDict)
	if err != nil {
		logger.Log.Fatalf("Open user dict file err, %s", err)
	}
	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		user := strings.TrimSpace(scanner.Text())
		if user != "" {
			users = append(users, user)
		}
	}
	return users,err
}

func ReadPasswordDict(passDict string)(password []string,err error){
	file, err := os.Open(passDict)
	if err != nil {
		logger.Log.Fatalf("Open password dict file err, %s", err)
	}
	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		pass := strings.TrimSpace(scanner.Text())
		if pass != "" {
			password = append(password, pass)
		}
	}
	password = append(password,"")
	return password,err
}