package todo

import (
	doMysql "cicdServer/mypack_mysql"
	doRedis "cicdServer/mypack_redis"
	sp "cicdServer/struct_pack"
	"database/sql"
	"encoding/json"
	"github.com/redis/go-redis/v9"
)

//type AddDeployTask interface {
//	AddtoRedis() int
//	AddtoMysql() int
//}

type DeployTaskDetails sp.DeployTaskDetails

func (sq *DeployTaskDetails) AddTaskQueue(rdb *redis.Client) {
	v, err := json.Marshal(sq)
	if err != nil {
		//fmt.Println("marshal failed!", err)
		return
	}
	//rdb := doRedis.InitRedisConn()
	doRedis.LpushList(rdb, "deployTaskList", string(v))
}

func (sq *DeployTaskDetails) AddTaskDB(sqldb *sql.DB) {
	//sqldb := doMysql.InitConn()
	doMysql.InsertTask(sqldb, (*doMysql.DeployTaskDetails)(sq))
}

func (sq *DeployTaskDetails) AddTaskInitState(sqldb *sql.DB) {
	//sqldb := doMysql.InitConn()
	doMysql.UpdateTaskInitStats(sqldb, (*doMysql.DeployTaskDetails)(sq))
}

func AddDeployTask(sq *DeployTaskDetails, sqldb *sql.DB, rdb *redis.Client) {
	sq.AddTaskDB(sqldb)
	sq.AddTaskInitState(sqldb)
	sq.AddTaskQueue(rdb)
}

func GetDeployTask(rdb *redis.Client) (taskInfo string) {
	taskInfo = doRedis.RpopList(rdb, "deployTaskList")
	return taskInfo
}
