package todo

import (
	slogger "cicdServer/log/wlog"
	doMysql "cicdServer/mypack_mysql"
	sp "cicdServer/struct_pack"
	"database/sql"
	"encoding/json"
)

type DeployTaskState sp.DeployTaskState
type DeployNodeState sp.DeployNodeState

func (cl *DeployNodeState) DealHeartbeat(sqldb *sql.DB) {
	doMysql.UpdateNodeState(sqldb, (*doMysql.DeployNodeState)(cl))
}

func (cl *DeployTaskState) DealDeployResult(sqldb *sql.DB) {
	doMysql.UpdateTaskState(sqldb, (*doMysql.DeployTaskState)(cl))

}

func DealClientPackage(cinfoJson string, sqldb *sql.DB) (ct, cp string, num float64) {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(cinfoJson), &m)
	if err != nil {
		slogger.Error(err.Error())
	}
	category := m["category"]

	if category == "heartbeat" {
		clinfo := new(DeployNodeState)
		clinfo.NodeId = m["nodeId"].(string)
		clinfo.HostName = m["hostName"].(string)
		clinfo.NodeUpdateTs = m["lastUpdate"].(float64)
		n := m["deployRunningNum"]
		if n == nil {
			clinfo.LoadNum = 0
		} else {
			clinfo.LoadNum = n.(float64)
		}
		clinfo.DealHeartbeat(sqldb)
		cp = category.(string)
		num = clinfo.LoadNum
		ct = m["clientToken"].(string)
		return ct, cp, num
	} else if category == "deployState" {
		ct = m["clientToken"].(string)
		dlinfo := new(DeployTaskState)
		agentInfoMap := m["agentDeployInfo"].(map[string]interface{})
		dlinfo.TaskId = agentInfoMap["taskId"].(string)
		dlinfo.StartTime = agentInfoMap["startTime"].(string)
		dlinfo.NodeId = m["nodeId"].(string)
		dlinfo.DeployDetails = agentInfoMap["taskInfo"].(string)
		dlinfo.DeployState = agentInfoMap["taskState"].(float64)
		dlinfo.PushImageName = agentInfoMap["pushImageName"].(string)
		dlinfo.DealDeployResult(sqldb)
	}
	return
}
