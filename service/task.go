package service

import (
	"time"

	"github.com/airportr/miaospeed/service/taskpoll"
	"github.com/airportr/miaospeed/utils"
)

var SpeedTaskPoll *taskpoll.TPController
var ConnTaskPoll *taskpoll.TPController

func StartTaskServer() {
	SpeedTaskPoll = taskpoll.NewTaskPollController("SpeedPoll", 1, time.Duration(utils.GCFG.PauseSecond)*time.Second, 200*time.Millisecond)
	ConnTaskPoll = taskpoll.NewTaskPollController("ConnPoll", utils.GCFG.ConnTaskTreading, 0, 200*time.Millisecond)

	go SpeedTaskPoll.Start()
	go ConnTaskPoll.Start()
}
