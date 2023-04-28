package tools

import (
	slogger "cicdServer/log/wlog"
	"github.com/google/uuid"
	"time"
)

func UtcDateConvert(utcTime string) (d string) {
	t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", utcTime)
	d = t.Local().Format("2006-01-02 15:04:05")
	return
}

func CreateTaskId() string {
	id2, err := uuid.NewRandom()
	if err != nil {
		//fmt.Printf("%v\n", err)
		slogger.Error(err.Error())
		return ""
	}
	return id2.String()
}
