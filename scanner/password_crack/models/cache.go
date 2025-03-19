package models 

import (
	"encoding/gob"
	"fmt"
	"os"
	"strings"
	"time"

	"SafeDP/scanner/password_crack/logger"
	"SafeDP/scanner/password_crack/util/hash"
	"SafeDP/scanner/password_crack/vars"

	"github.com/patrickmn/go-cache"
)

func init() {
	gob.Register(Service{})
	gob.Register(ScanResult{})
}

func SaveResult(result ScanResult, err error){
	if err == nil && result.Result {
		var k string 
		protocol := strings.ToLower(result.Service.Protocol)

		if protocol == "REDIS" {
			k = fmt.Sprintf("%v-%v-%v", result.Service.Ip, result.Service.Port, result.Service.Protocol)
		} else {
			k = fmt.Sprintf("%v-%v-%v", result.Service.Ip, result.Service.Port, result.Service.Username)
		}

		h := hash.MakeTaskHash(k)
		hash.SetTaskHash(h)

		_, found: = vars.CacheService.Get(k)
		if !found {
			logger.Log.Infof("Ip: %v, Port: %v, Protocol: [%v], Username: %v, Password: %v", result.Service.Ip,
				result.Service.Port, result.Service.Protocol, result.Service.Username, result.Service.Password)
		}
		vars.CasheService.Set(k, result, cache.NoExpiration)
	}
}