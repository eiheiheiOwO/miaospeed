package taskpoll

type TaskPollItem interface {
	ID() string
	TaskName() string
	Weight() uint
	Count() int

	Yield(i int, tpc *TPController)
	OnExit(exitCode TPExitCode)
	Init() TaskPollItem
}
