package scanner

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"SafeDp/scanner/tcp_scanner_1/vars"
)

func GenerateTask(ipList []net.IP, ports[]int) ([]map[string]int,int) {
	tasks := make([]map[string]int, 0)
	
	for _, ip := range ipList {
		for _, port := range ports {
			ipPort := map[string]int{ip.String(): port}
			tasks = append(tasks, ipPort)
		}
	}
	return tasks, len(tasks)
}

func AssigningTasks (tasks []map[string]int) {
	scanBatch := len(tasks) / vars.ThreadNum
	for i := 0; i < scanBatch; i++ {
		curTask := tasks[vars.ThreadNum * i:vars.ThreadNum * (i + 1)]
		RunTask(curTask)
	}

	if len(tasks) % vars.ThreadNum > 0 {
		lastTask := tasks[vars.ThreadNum * scanBatch:]
		RunTask(lastTask)
	}
}

func RunTask(tasks []map[string]int) {
	var wg sync.WaitGroup
	wg.Add(len(tasks))
	for _, task := range tasks {
		for ip, port := range task {
			go func(ip string,port int) {
				err := SaveResult(Connect(ip, port))
				_ = err
				wg.Done()
			}(ip, port)
		}
	}
	wg.Wait()
}
func SaveResult(ip string, port int, err error) error {
	if err != nil {
		return err
	}
	
	v, ok := vars.Result.Load(ip)
	if ok {
		ports,ok1 := v.([]int)
		if ok1 {
			ports = append(ports, port)
			vars.Result.Store(ip, ports)
		}
	} else {
		ports := make([]int, 0)
		ports = append(ports, port)
		vars.Result.Store(ip, ports)
	}
	return err
}


func PrintResult() {
	vars.Result.Range(func(key, value interface{}) bool {
		fmt.Println("IP:%v \n", key)
		fmt.Println("Port:%v \n", value)
		fmt.Println(strings.Repeat("-", 100))
		return true
	})
}