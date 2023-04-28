package main

import (
	dimg "cicdAgent/deal_image"
	dcode "cicdAgent/deploy_code"
	alogger "cicdAgent/log/wlog"
	sp "cicdAgent/struct_pack"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

const clientToken = "hhHA6VcuTg8QIjVpsUB7eHhWmc7A0pTi"

var deployRunningNum int64 = 0

type AgentDeployInfo sp.AgentDeployInfo
type ClientInfo sp.ClientInfo

func main() {
	f := "/tmp/agent.log"
	alogger.InitLogger(f)
	CONNECT := "serverip:12122"
	connTimeout := 2 * time.Second
	for {
		conn, err := net.DialTimeout("tcp", CONNECT, connTimeout)
		if err != nil {
			alogger.Error(err.Error())
			time.Sleep(3 * time.Second)
			continue
		}
		go clientHandleConn(conn)
		sendHeartBeat(conn)
		alogger.Info("reconnect to server")
	}
}

func getHostName() (string, string) {
	hostName, _ := os.Hostname()
	md5str := md5.Sum([]byte(hostName))
	md5HostName := fmt.Sprintf("%x", md5str)
	return md5HostName, hostName
}

func strutcToJson(s *ClientInfo) (t string) {
	v, err := json.Marshal(s)
	if err != nil {
		alogger.Error(err.Error())
	}
	return string(v)
}

func sendHeartBeat(conn net.Conn) {
	defer conn.Close()
	md5HostName, hostName := getHostName()
	for {
		ts := time.Now().Unix()
		ms := &ClientInfo{ClientToken: clientToken, NodeId: md5HostName, HostName: hostName, LastUpdate: ts, Category: "heartbeat", DeployRunningNum: float64(deployRunningNum), AgentDeployInfo: sp.AgentDeployInfo(AgentDeployInfo{})}
		js, _ := json.Marshal(ms)
		_, err := conn.Write([]byte(js))
		if err != nil {
			alogger.Error("disconnect from server " + err.Error())
			break
		}
		//send heartbeat pre 5 sec
		time.Sleep(5 * time.Second)

		alogger.Info("Send to server " + string(js))
		//fmt.Println("Send:", ms)
		continue
	}
}

// deal heartbeat request
func clientHandleConn(conn net.Conn) {
	defer conn.Close()
	readChan := make(chan string, 1)
	writeChan := make(chan string, 1)

	go readFromServer(conn, readChan)
	go writeToServer(conn, writeChan)

	for {
		select {
		case readStr := <-readChan:
			//fmt.Println(readStr)
			go func() {
				deployTask := make(map[string]string)
				err := json.Unmarshal([]byte(readStr), &deployTask)
				if err != nil {
					//fmt.Printf("Unmarshal with error: %+v\n", err)
					alogger.Error(err.Error())
				}
				//fmt.Println(deployTask)
				md5HostName, hostName := getHostName()
				ts := time.Now().Unix()
				taskId := deployTask["taskId"]
				startTime := time.Now().Format("2006-01-02 15:04:05")
				taskStat := 100
				taskInfo := "Deploying"
				deployRunningNum += 1
				ms1 := &ClientInfo{ClientToken: clientToken, NodeId: md5HostName, HostName: hostName, LastUpdate: ts, Category: "deployStats", DeployRunningNum: float64(deployRunningNum), AgentDeployInfo: sp.AgentDeployInfo(AgentDeployInfo{TaskId: taskId, StartTime: startTime, TaskInfo: taskInfo, TaskStats: taskStat})}
				strDate := strutcToJson(ms1)
				writeChan <- strDate
				deployResult, deployTaskId, deployDetails, pushImageName := dealDeployTask(deployTask)
				ts = time.Now().Unix()
				taskStat = deployResult
				taskInfo = deployDetails
				deployRunningNum -= 1
				ms2 := &ClientInfo{ClientToken: clientToken, NodeId: md5HostName, HostName: hostName, LastUpdate: ts, Category: "deployStats", DeployRunningNum: float64(deployRunningNum), AgentDeployInfo: sp.AgentDeployInfo(AgentDeployInfo{TaskId: deployTaskId, StartTime: startTime, TaskInfo: taskInfo, TaskStats: taskStat, PushImageName: pushImageName})}
				strDate = strutcToJson(ms2)
				writeChan <- strDate
			}()
		}
	}
}

func readFromServer(conn net.Conn, readChan chan<- string) {
	for {
		data := make([]byte, 1024)
		total, err := conn.Read(data)
		if err != nil {
			alogger.Error("disconnet from server " + err.Error())
			break
		}

		strData := string(data[:total])
		alogger.Info("received from server " + strData)
		readChan <- strData
	}
}

func writeToServer(conn net.Conn, writeChan <-chan string) {
	for {
		strData := <-writeChan
		_, err := conn.Write([]byte(strData))
		if err != nil {
			//fmt.Println(err)
			alogger.Error(err.Error())
			break
		}
		alogger.Info("send to server " + strData)
	}
}

func dealDeployTask(deployTask map[string]string) (deployResult int, taskId, deployDetails, pushImageName string) {
	taskId = deployTask["taskId"]
	projectName := deployTask["projectName"]
	branchName := deployTask["branchName"]
	branch := deployTask["branch"]
	//gitSshUrl := deployTask["gitSshUrl"]
	githttpUrl := deployTask["gitHttpUrl"]
	//fmt.Println(taskId, projectName, branchName, branch, gitSshUrl, githttpUrl)
	di := new(sp.DeployBasicInfo)
	di.TaskId = taskId
	di.ProjectName = projectName
	di.BranchName = branchName
	di.Branch = branch
	di.GithttpUrl = githttpUrl
	deployStats := dcode.DoDeployTask((*dcode.DeployBasicInfo)(di))
	if deployStats == 1 {
		res, imageName := dimg.DoImageTask((*dimg.DeployBasicInfo)(di))
		if res == 1 {
			deployResult = 1
			deployDetails = "deploy success"
			pushImageName = imageName
			return
		} else {
			deployResult = 0
			deployDetails = "deploy err"
			pushImageName = ""
			return
		}
	} else {
		deployResult = 0
		deployDetails = "deploy err"
		pushImageName = ""
		return
	}
}
