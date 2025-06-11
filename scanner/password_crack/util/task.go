package util

import (
	"fmt"
	"runtime"
	"sync"
	"strings"
	"time"
	"safedp/scanner/password_crack/logger"
	"safedp/scanner/password_crack/models"
	"safedp/scanner/password_crack/vars"
	"safedp/scanner/password_crack/util/hash"
	"safedp/scanner/password_crack/plugins"
	"gopkg.in/chaeggaa/pb.v2"
	"github.com/urfave/cli"
	"github.com/sirupsen/logrus"
)

func GenerateTask(ipList []models.IpAddr,users []string, passwords []string) (tasks []models.Service,taskNum int) {
	tasks = make([]models.Service, 0) 
	for _, user := range users {
		for _, password := range passwords {
			for _, ipAddr := range ipList {
				task := models.Service{
					Ip:       ipAddr.Ip,
					Port:     ipAddr.Port,
					Protocol: ipAddr.Protocol,
					User:     user,
					Password: password,
				}
				tasks = append(tasks, task)
			}
		}
	}
	taskNum = len(tasks)
	return tasks, taskNum
}

func RunTask(tasks []models.Service){
	totalTask := len(tasks)
	vars.ProgressBar = pb.StartNew(totalTask)
	vars.ProgressBar.SetTemplate(`{{ rndcolor "Scanning progress: " }} {{  percent . "[%.02f%%]" "[?]"| rndcolor}} {{ counters . "[%s/%s]" "[%s/?]" | rndcolor}} {{ bar . "「" "-" (rnd "ᗧ" "◔" "◕" "◷" ) "•" "」" | rndcolor }} {{rtime . | rndcolor}} `)

	wg := &sync.WaitGroup{}
	//创建Channel
	taskChan := make(chan models.Service, vars.ScanNum)
	//创建协程
	for i := 0; i < vars.ScanNum; i++ {
		go crackPassword(taskChan, wg)
	}
	//生产者，将任务放入Channel
	for _, task := range tasks {
		wg.Add(1)
		taskChan <- task
	}

	close(taskChan)
	waitTimeout(wg, vars.TimeOut)
	{
		_ = models.SaveResultToFile()
		models.ResultTotal()
		_ = models.DumpToFile(vars.ResultFile)
	}
}
			
func crackPassword(taskChan chan models.Service, wg *sync.WaitGroup) {
	for task := range taskChan {
		vars.ProgressBar.Increment()
		if vars.DebugMode {
			logger.Log.Debugf("checking: Ip: %v, Port: %v, [%v], Username: %v, Password: %v, goroutineNum: %v",
				task.Ip, task.Port, task.Protocol, task.User, task.Password, runtime.NumGoroutine())
		} 
		var k string
		protocol := strings.ToUpper(task.Protocol)
		if protocol == "REDIS" {
			k = fmt.Sprintf("%v-%v-%v", task.Ip, task.Port, task.Protocol)
		} else {
			k = fmt.Sprintf("%v-%v-%v", task.Ip, task.Port, task.Username)
		}

		h := hash.MakeTaskHash(k)
		if hash.CheckTaskHash(h) {
			wg.Done()
			continue
		}
		fn := plugins.ScanFuncMap[protocol]
		models.SaveResult(fn(task))
		wg.Done()
	}
}

func Scan(ctx *cli.Context) (err error){
	if ctx.IsSet("debug") {
		vars.DebugMode = ctx.Bool("debug")
	}

	if vars.DebugMode {
		logger.Log.Level = logrus.DebugLevel
	}

	if ctx.IsSet("timeout") {
		vars.TimeOut = time.Duration(ctx.Int("timeout")) * time.Second
	}
		if ctx.IsSet("scan_num") {
		vars.ScanNum = ctx.Int("scan_num")
	}

	if ctx.IsSet("ip_list") {
		vars.IpList = ctx.String("ip_list")
	}

	if ctx.IsSet("user_dict") {
		vars.UserDict = ctx.String("user_dict")
	}

	if ctx.IsSet("pass_dict") {
		vars.PassDict = ctx.String("pass_dict")
	}

	if ctx.IsSet("outfile") {
		vars.ResultFile = ctx.String("outfile")
	}
	vars.StartTime = time.Now()

	userDict, uErr := ReadUserDict(vars.UserDict)
	passDict, pErr := ReadPasswordDict(vars.PassDict)

	ipList := ReadIpList(vars.IpList)
	aliveIpList := CheckAlive(ipList)
	if uErr == nil && pErr == nil {
		tasks, _ := GenerateTask(aliveIpList, userDict, passDict)
		RunTask(tasks)
	}
	return	err
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}