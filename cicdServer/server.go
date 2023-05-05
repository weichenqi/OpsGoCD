package main

import (
	slogger "cicdServer/log/wlog"
	dlogger "cicdServer/log/zlog"
	doMysql "cicdServer/mypack_mysql"
	doRedis "cicdServer/mypack_redis"
	"cicdServer/todo"
	tools "cicdServer/wcqtools"
	"github.com/go-playground/webhooks/v6/gitlab"
	"net"
	"net/http"
	"os"
	"strings"
)

const (
	taskNum     = 3
	path        = "/webhooks"
	logFile     = "/tmp/server.log"
	dbLogFile   = "/tmp/db.log"
	clientToken = "hhHA6VcuTg8QIjVpsUB7eHhWmc7A0pTi" //客户端和服务端交互的简单认证，两边一样即可，不要纠结token本身。
)

var sqldb = doMysql.InitConn()
var rdb = doRedis.InitRedisConn()

func main() {
	slogger.InitLogger(logFile)
	dlogger.InitLogger(dbLogFile)

	go gitLabWebhook()
	tcpDeployServer()
}

func tcpDeployServer() {
	PORT := ":12122"
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		slogger.Error(err.Error())
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			slogger.Error(err.Error())
			continue
		}
		go handleConn(conn)
	}
}

func gitLabWebhook() {
	hook, _ := gitlab.New(gitlab.Options.Secret("*****"))

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, gitlab.PushEvents)
		if err != nil {
			if err == gitlab.ErrEventNotFound {
			}
		}

		switch payload.(type) {

		case gitlab.PushEventPayload:
			release := payload.(gitlab.PushEventPayload)
			pushDate := release.Commits[0].Timestamp.String()
			projectName := strings.Split(release.Repository.Name, "/")[0]
			branch := strings.Split(release.Ref, "/")[2]
			branchName := release.Ref
			gitSshUrl := release.Repository.GitSSHURL
			gitHttpUrl := release.Repository.GitHTTPURL
			pushUserName := release.UserUsername
			//insert task to redis and mysql
			dtd := new(todo.DeployTaskDetails)
			dtd.TaskId = tools.CreateTaskId()
			dtd.PushDate = tools.UtcDateConvert(pushDate)
			dtd.ProjectName = projectName
			dtd.BranchName = branchName
			dtd.Branch = branch
			dtd.GitSshUrl = gitSshUrl
			dtd.GitHttpUrl = gitHttpUrl
			dtd.PushUserName = pushUserName
			todo.AddDeployTask(dtd, sqldb, rdb)
		}

	})
	http.ListenAndServe(":6060", nil)
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	readChan := make(chan string)
	writeChan := make(chan string)
	stopChan := make(chan bool)

	go readConn(conn, readChan, stopChan)
	go writeConn(conn, writeChan, stopChan)

	for {
		select {
		case readStr := <-readChan:
			ct, cp, num := todo.DealClientPackage(readStr, sqldb)
			if cp == "heartbeat" && ct == clientToken && num < 3 {
				upper := todo.GetDeployTask(rdb)
				if upper != "" {
					writeChan <- upper
				}
			} else if ct != clientToken {
				slogger.Info("Unauthorized client try to connect " + conn.RemoteAddr().String())
			} else if cp == "heartbeat" && ct == clientToken && num >= taskNum {
				slogger.Info("Agent's load task gt 3 " + conn.RemoteAddr().String())
			}
		case stop := <-stopChan:
			if stop {
				break
			}
		}
	}
}

func readConn(conn net.Conn, readChan chan<- string, stopChan chan<- bool) {
	for {
		data := make([]byte, 1024)
		total, err := conn.Read(data)
		if err != nil {
			slogger.Info("Can not received client's data " + err.Error())
			break
		}

		strData := string(data[:total])
		slogger.Info("Received from client " + strData + " " + conn.RemoteAddr().String())
		readChan <- strData
	}
	slogger.Info("Connection has been disconnected" + conn.RemoteAddr().String())
	stopChan <- true
	//fmt.Println("chan closed")
}

func writeConn(conn net.Conn, writeChan <-chan string, stopChan chan<- bool) {
	for {
		strData := <-writeChan
		_, err := conn.Write([]byte(strData))
		if err != nil {
			slogger.Error(err.Error())
			break
		}
		slogger.Info("Send data to client " + strData + " " + conn.RemoteAddr().String())
	}
	slogger.Info("Connection has been disconnected" + conn.RemoteAddr().String())
	stopChan <- true
	//fmt.Println("chan closed")
}
