package models
import (
	"time"
)

type Agent struct {
	Id 	  int64
	AgentId string `json:"agent_id"`
	Platform string `json:"platform"`
	Architecture string `json:"architecture"`
	UserName string `json:"username"`
	UserGUID string `json:"user_guid"`
	HostName string `json:"host_name"`
	Ips []string `json:"ips" xorm:"text"`
	Pid int `json:"pid"`
	Debug bool `json:"debug"`
	UserAgent string `json:"user_agent"`
	Initial bool `json:"initial"`
	CreateTime time.Time `xorm:"created"`
	UpdateTime time.Time `xorm:"updated"`
	Version int `xorm:"version"`
}

func (a *Agent) Insert() error {
	_, err := Engine.Insert(a)
	return err
}

func ListAgents() ([]Agent, error) {
	agents := make([]Agent,0)
	err := Engine.Find(&agents)
	return agents, err
}

func UpdateAgent(agentId string) error {
	agent := new(Agent)
	has, err := Engine.Where("agent_id = ?", agentId).Get(agent)
	if err != nil {
		return err
	}
	if has {
		_, err = Engine.Id(agent.Id).Update(agent)
	}
	return err
}

func ExistAgentId(agentId string) (bool, error) {
	agent := new(Agent)
	has, err := Engine.Where("agent_id = ?", agentId).Get(agent)
	return has, err
}

func RemoveAll() error {
	_, err := Engine.Exec("DELETE FROM agent")
	return err
}

