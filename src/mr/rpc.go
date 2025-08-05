package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import (
	"os"
	"strconv"
)

// Add your RPC definitions here.
type TaskAssignArgs struct{}

type TaskAssignReply struct {
	TaskType   string
	inputFiles []string
	TaskNum    int
	nReduce    int
}

type NewReduceTaskArgs struct {
	TaskNames []string
}

type NewReduceTaskReply struct{}

type ReduceTaskDoneArgs struct {
	TaskNum int
}

type ReduceTaskDoneReply struct{}

type InvalidTaskArgs struct {
	TaskType string
	TaskName string
	taskNum  int
}

type InvalidTaskReply struct{}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/5840-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}
