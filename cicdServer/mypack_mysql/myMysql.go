package myMysql

import (
	dlogger "cicdServer/log/zlog"
	sp "cicdServer/struct_pack"
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"os"
	"strconv"
)

var ctx = context.Background()

type DeployTaskState sp.DeployTaskState
type DeployTaskDetails sp.DeployTaskDetails
type DeployNodeState sp.DeployNodeState

func InitConn() *sql.DB {
	// Set up the DSN
	dsn := "cicd_user:******@tcp(127.0.0.1:3306)/wcicd"
	// Open a connection to the database
	sqldb, err := sql.Open("mysql", dsn)
	if err != nil {
		dlogger.Error(err.Error())
		os.Exit(1)
	}

	// Set the maximum number of open connections to the database
	sqldb.SetMaxOpenConns(50)

	// Set the maximum number of idle connections in the connection pool
	sqldb.SetMaxIdleConns(100)

	err = sqldb.Ping()
	if err != nil {
		dlogger.Error(err.Error())
		os.Exit(1)
	}
	return sqldb
}

func InsertTask(sqldb *sql.DB, deployTaskDetails *DeployTaskDetails) {
	db := bun.NewDB(sqldb, mysqldialect.New())
	taskId := deployTaskDetails.TaskId
	pushDate := deployTaskDetails.PushDate
	projectName := deployTaskDetails.ProjectName
	branchName := deployTaskDetails.BranchName
	branch := deployTaskDetails.Branch
	gitSshUrl := deployTaskDetails.GitSshUrl
	gitHttpUrl := deployTaskDetails.GitHttpUrl
	pushName := deployTaskDetails.PushUserName
	task := &DeployTaskDetails{TaskId: taskId, PushDate: pushDate, ProjectName: projectName, BranchName: branchName, Branch: branch, GitSshUrl: gitSshUrl, GitHttpUrl: gitHttpUrl, PushUserName: pushName}
	res, err := db.NewInsert().Model(task).Exec(ctx)
	if err != nil {
		dlogger.Error(err.Error())
	} else {
		l, _ := res.RowsAffected()
		dlogger.Info("insert task successful " + strconv.FormatInt(l, 10))
	}
}

// deploy task init state
func UpdateTaskInitState(sqldb *sql.DB, deployTaskDetails *DeployTaskDetails) {
	db := bun.NewDB(sqldb, mysqldialect.New())
	taskId := deployTaskDetails.TaskId
	startTime := "0000-00-00 00:00:00"
	deployState := 200
	deployDetails := "waiting for deploy"
	task := DeployTaskState{TaskId: taskId, StartTime: startTime, DeployState: float64(deployState), DeployDetails: deployDetails}
	res, err := db.NewInsert().Model(&task).Exec(ctx)
	if err != nil {
		//fmt.Println(err)
		dlogger.Error(err.Error())
	} else {
		//fmt.Println(res.RowsAffected())
		l, _ := res.RowsAffected()
		dlogger.Info("update task initial state successful " + strconv.FormatInt(l, 10))
	}
}

// deploy task state
func UpdateTaskState(sqldb *sql.DB, deployTaskState *DeployTaskState) {
	db := bun.NewDB(sqldb, mysqldialect.New())
	taskId := deployTaskState.TaskId
	startTime := deployTaskState.StartTime
	deployNode := deployTaskState.NodeId
	deployState := deployTaskState.DeployState
	deployDetails := deployTaskState.DeployDetails
	pushImageName := deployTaskState.PushImageName
	task := DeployTaskState{TaskId: taskId, StartTime: startTime, NodeId: deployNode, DeployState: deployState, DeployDetails: deployDetails, PushImageName: pushImageName}
	res, err := db.NewInsert().
		Model(&task).
		On("DUPLICATE KEY UPDATE").
		Exec(ctx)
	if err != nil {
		dlogger.Error(err.Error())
	} else {
		l, _ := res.RowsAffected()
		dlogger.Info("update task  state successful " + strconv.FormatInt(l, 10))
	}
}

// agent state from heartbeat
func UpdateNodeState(sqldb *sql.DB, deployNodeState *DeployNodeState) {
	db := bun.NewDB(sqldb, mysqldialect.New())
	nodeId := deployNodeState.NodeId
	hostName := deployNodeState.HostName
	loadNum := deployNodeState.LoadNum
	nodeUpdateTs := deployNodeState.NodeUpdateTs
	task := DeployNodeState{NodeId: nodeId, HostName: hostName, LoadNum: loadNum, NodeUpdateTs: nodeUpdateTs}
	res, err := db.NewInsert().
		Model(&task).
		On("DUPLICATE KEY UPDATE").
		Exec(ctx)
	if err != nil {
		dlogger.Error(err.Error())
	} else {
		l, _ := res.RowsAffected()
		dlogger.Info("update agent state success " + strconv.FormatInt(l, 10))
	}
}
