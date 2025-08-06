package mr

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"
)

type taskInfo struct {
	status     string // "unassigned", "in-progress", "done"
	inputFiles []string
	num        int
	workerID   string
	assignedAt time.Time
}

type Coordinator struct {
	mu               sync.Mutex
	workersStatus    map[string]bool
	mapTasksInfos    map[string]*taskInfo
	reduceTasksInfos []*taskInfo
	nReduce          int
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{
		workersStatus:    make(map[string]bool),
		mapTasksInfos:    make(map[string]*taskInfo, len(files)),
		reduceTasksInfos: make([]*taskInfo, nReduce),
		nReduce:          nReduce,
	}
	for taskNum, taskName := range files {
		c.mapTasksInfos[taskName] = &taskInfo{
			status:     "unassigned",
			inputFiles: []string{taskName},
			num:        taskNum,
		}
	}
	for taskNum := range nReduce {
		c.reduceTasksInfos[taskNum] = &taskInfo{
			status: "unassigned",
			num:    taskNum,
		}
	}
	c.server()
	return &c
}

// Your code here -- RPC handlers for the worker to call.
func (c *Coordinator) TaskAssign(
	args *TaskAssignArgs,
	reply *TaskAssignReply,
) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// map
	isMapFinished := true
	c.workersStatus[args.WorkerID] = true
	for taskName, taskInfo := range c.mapTasksInfos {
		switch taskInfo.status {
		case "unassigned":
			reply.TaskType = "map"
			reply.InputFiles = []string{taskName}
			reply.TaskNum = taskInfo.num
			reply.NReduce = c.nReduce
			taskInfo.status = "in-progress"
			taskInfo.workerID = args.WorkerID
			taskInfo.assignedAt = time.Now()
			return nil
		case "in-progress":
			now := time.Now()
			if now.Sub(taskInfo.assignedAt) < 10*time.Second {
				isMapFinished = false
				continue
			}
			currentWorkerID := taskInfo.workerID
			if _, ok := c.workersStatus[currentWorkerID]; !ok {
				log.Fatal("unknown workerID: ", currentWorkerID)
			}
			// Prevent the worker's status from being set to false
			// when assigning the same task to the same worker.
			// This situation can occur if there was a network issue
			// and the worker could not be reached;
			// in that case, treat the worker as recovered and resend the same task.
			if currentWorkerID != args.WorkerID {
				c.workersStatus[currentWorkerID] = false
			}
			reply.TaskType = "map"
			reply.InputFiles = []string{taskName}
			reply.TaskNum = taskInfo.num
			reply.NReduce = c.nReduce
			taskInfo.status = "in-progress"
			taskInfo.workerID = args.WorkerID
			taskInfo.assignedAt = time.Now()
			return nil
		case "done":
			continue
		default:
			log.Fatal("invalid task status: ", taskInfo.status)
		}
	}

	// wait
	// Cannot start reduce tasks until all map tasks are finished
	if !isMapFinished {
		reply.TaskType = "wait"
		return nil
	}

	// reduce
	for taskNum, taskInfo := range c.reduceTasksInfos {
		if taskNum != taskInfo.num {
			log.Fatal("expected taskNum = ", taskNum, ", but got = ", taskInfo.num)
		}
		switch taskInfo.status {
		case "unassigned":
			if len(taskInfo.inputFiles) == 0 {
				continue
			}
			reply.TaskType = "reduce"
			reply.InputFiles = taskInfo.inputFiles
			reply.TaskNum = taskInfo.num
			reply.NReduce = c.nReduce
			taskInfo.status = "in-progress"
			taskInfo.workerID = args.WorkerID
			taskInfo.assignedAt = time.Now()
			return nil
		case "in-progress":
			now := time.Now()
			if now.Sub(taskInfo.assignedAt) < 10*time.Second {
				continue
			}
			currentWorkerID := taskInfo.workerID
			if _, ok := c.workersStatus[currentWorkerID]; !ok {
				log.Fatal("unknown workerID: ", currentWorkerID)
			}
			if currentWorkerID != args.WorkerID {
				c.workersStatus[currentWorkerID] = false
			}
			reply.TaskType = "reduce"
			reply.InputFiles = taskInfo.inputFiles
			reply.TaskNum = taskInfo.num
			reply.NReduce = c.nReduce
			taskInfo.status = "in-progress"
			taskInfo.workerID = args.WorkerID
			taskInfo.assignedAt = time.Now()
			return nil
		case "done":
			continue
		default:
			log.Fatal("invalid task status: ", taskInfo.status)
		}
	}

	// exit
	reply.TaskType = "exit"
	return nil
}

func (c *Coordinator) NewReduceTask(
	args *NewReduceTaskArgs,
	reply *NewReduceTaskReply,
) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// update map task status
	if _, ok := c.workersStatus[args.WorkerID]; !ok {
		log.Fatal("unknown workerID: ", args.WorkerID)
	}
	mapTaskInfo, ok := c.mapTasksInfos[args.MapTaskName]
	if !ok {
		log.Fatal("unknown mapTaskName: ", args.MapTaskName)
	}
	if mapTaskInfo.workerID != args.WorkerID {
		log.Fatal("expected workerID = ", mapTaskInfo.workerID, ", but got = ", args.WorkerID)
	}
	mapTaskInfo.status = "done"

	// update reduce input files
	for reduceTaskNum, inputFile := range args.ReduceInputFiles {
		info := c.reduceTasksInfos[reduceTaskNum]
		info.inputFiles = append(info.inputFiles, inputFile)
	}

	return nil
}

func (c *Coordinator) ReduceTaskDone(
	args *ReduceTaskDoneArgs,
	reply *ReduceTaskDoneReply,
) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.workersStatus[args.WorkerID]; !ok {
		log.Fatal("unknown workerID: ", args.WorkerID)
	}
	if args.TaskNum < 0 || args.TaskNum >= len(c.reduceTasksInfos) {
		log.Fatal("unknown reduceTaskNum: ", args.TaskNum)
	}
	reduceTaskInfo := c.reduceTasksInfos[args.TaskNum]
	if reduceTaskInfo.workerID != args.WorkerID {
		log.Fatal("expected workerID = ", reduceTaskInfo.workerID, ", but got = ", args.WorkerID)
	}
	reduceTaskInfo.status = "done"
	return nil
}

func (c *Coordinator) InvalidTask(
	args *InvalidTaskArgs,
	reply *InvalidTaskReply,
) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.workersStatus[args.WorkerID]; !ok {
		log.Fatal("unknown workerID: ", args.WorkerID)
	}
	switch args.TaskType {
	case "map":
		mapTaskInfo, ok := c.mapTasksInfos[args.TaskName]
		if !ok {
			log.Fatal("unknown mapTaskName: ", args.TaskName)
		}
		if mapTaskInfo.workerID != args.WorkerID {
			log.Fatal("expected workerID = ", mapTaskInfo.workerID, ", but got = ", args.WorkerID)
		}
		mapTaskInfo.status = "done"
	case "reduce":
		if args.TaskNum < 0 || args.TaskNum >= len(c.reduceTasksInfos) {
			log.Fatal("unknown reduceTaskNum: ", args.TaskNum)
		}
		reduceTaskInfo := c.reduceTasksInfos[args.TaskNum]
		if reduceTaskInfo.workerID != args.WorkerID {
			log.Fatal("expected workerID = ", reduceTaskInfo.workerID, ", but got = ", args.WorkerID)
		}
		reduceTaskInfo.status = "done"
	default:
		fmt.Printf("invalid task type, %v", args.TaskType)
		return fmt.Errorf("invalid task type, %v", args.TaskType)
	}
	return nil
}

// start a thread that listens for RPCs from worker.go
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.
	// MapReduce ジョブが完全に終了したときに true を返す

	return ret
}

// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}
