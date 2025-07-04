package models

import (
	"time"
)

type(
	Command struct {
		Id int64 `json:"id"`
		AgentId string `json:"agent_id"`
		Content string `json:"content"`
		Status int `json:"status"` 
		CreateTime time.Time `xorm:"created"`
		UpdateTime time.Time `xorm:"updated"`
		Result string `json:"result"`
	}
)

func NewCommand(agentId ,content string) *Command {
	return &Command{
		Id: 0,
		AgentId: agentId,
		Content: content,
		Status: 0, 
	}
}

func (c *Command) Insert() error {
	_, err := Engine.Insert(c)
	return err
}

func ListCommandByAgentId(agentId string) ([]Command, error) {
	commands := make([]Command, 0)
	err := Engine.Where("agent_id = ? and status = 0", agentId).Find(&commands)
	return commands, err
}

func ListFinishCommand()([]Command, error) {
	commands := make([]Command, 0)
	err := Engine.Where("status = 1").Find(&commands)
	return commands, err
}

func UpdateCommandResult(cmdId int64, result string) error {
	cmd := new(Command)
	has, err := Engine.Id(cmdId).Get(cmd)
	if err != nil {
		return err
	}
	if has {
		cmd.Result = result
		_, err = Engine.Id(cmdId).Update(cmd)
	}
	return err
}

func UpdateCommandStatus(cmdId int64, status int) error {
	cmd := new(Command)
	has, err := Engine.Id(cmdId).Get(cmd)
	if err != nil {
		return err
	}
	if has {
		cmd.Status = status
		_, err = Engine.Id(cmdId).Update(cmd)
	}
	return err
}

func SetCmdStatusToFinished(cmdId int64) error {
	return UpdateCommandStatus(cmdId, 1) 

}

func SetCmdStatusToEnd(cmdId int64) error {
	return UpdateCommandStatus(cmdId, 2) 
}


