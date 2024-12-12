package scanner

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"SafeDp/scanner/tcp_scanner_2/vars"
)

func GenerateTask(ipList []net.IP, ports []int) ([]map[string]int, int) {
	tasks := make([]map[string]int, 0)

	for _, ip := range ipList {
		for _, port := range ports {
			ipPort := map[string]int{ip.String(): port}
			tasks = append(tasks, ipPort)
		}
	}
	return tasks, len(tasks)
}

func RunTask(tasks []map[string]int) {
	wg := &sync.WaitGroup{}

	//创建一个buffer为vars.threadNum*2的channel
	taskChan := make(chan map[string]int, vars.ThreadNum*2)

	//创建vars.ThreadNum个协程
	for i := 0; i < vars.ThreadNum; i++ {
		go Scan(taskChan, wg)
	}

	//生产者，不断地往taskChan channnel中写入数据
	for _, task := range tasks {
		wg.Add(1)
		taskChan <- task
	}
	close(taskChan)
	wg.Wait()
}

func Scan(taskChan chan map[string]int, wg *sync.WaitGroup) {
	for task := range taskChan {
		for ip, port := range task {
			err := SaveResult(Connect(ip, port))
			_ = err
			wg.Done()
		}
	}
}

func SaveResult(ip string, port int, err error) error {
	if err != nil {
		return err
	}

	v, ok := vars.Result.Load(ip)
	if ok {
		ports, ok1 := v.([]int)
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
