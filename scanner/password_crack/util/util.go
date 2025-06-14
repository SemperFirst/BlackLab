package util

import (
	"sync"
	"fmt"
	"net"
	"gopkg.in/cheggaaa/pb.v2"
	"SafeDp/scanner/password_crack/models"
	"SafeDp/scanner/password_crack/logger"
	"SafeDp/scanner/password_crack/vars"
)

var (
	AliveAddr []models.IpAddr
	mutex sync.Mutex
)

func init() {
	AliveAddr = make([]models.IpAddr, 0)
}

func CheckAlive(ipList []models.IpAddr) []models.IpAddr {
	logger.Log.Infoln("Start checking alive hosts")
	vars.ProcessBarActive = pb.StartNew(len(ipList))
	vars.ProcessBarActive.SetTemplate(`{{ rndcolor "Checking progress: " }} {{  percent . "[%.02f%%]" "[?]"| rndcolor}} {{ counters . "[%s/%s]" "[%s/?]" | rndcolor}} {{ bar . "「" "-" (rnd "ᗧ" "◔" "◕" "◷" ) "•" "」" | rndcolor}}  {{rtime . | rndcolor }}`)

	var wg sync.WaitGroup
	wg.Add(len(ipList))

	for _, addr := range ipList {
		go func(addr models.IpAddr) {
			defer wg.Done()
			SaveAddr(check(addr))
		}(addr)
	}
	wg.Wait()
	vars.ProcessBarActive.Finish()
	return AliveAddr
}

func check(ipAddr models.IpAddr) (bool,models.IpAddr){
	alive := false
	_,err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ipAddr.Ip, ipAddr.Port), vars.TimeOut)
	if err == nil {
		alive = true
	}
	vars.ProcessBarActive.Increment()
	return alive, ipAddr
}


func SaveAddr(alive bool, ipAddr models.IpAddr) {
	if alive {
		mutex.Lock()
		AliveAddr = append(AliveAddr, ipAddr)
		mutex.Unlock()
	}
}
