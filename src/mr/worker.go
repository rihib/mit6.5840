package mr

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/rpc"
	"os"
	"sort"
	"strconv"
	"time"
)

// Map functions return a slice of KeyValue.
type KeyValue struct {
	Key   string
	Value string
}

// for sorting by key.
type ByKey []KeyValue

// for sorting by key.
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

// main/mrworker.go calls this function.
// The worker itself does not perform parallel processing;
// parallelism is achieved by running multiple worker processes.
func Worker(
	mapf func(string, string) []KeyValue,
	reducef func(string, []string) string,
) {
workerLoop:
	for {
		// Fetch task or exit worker
		args, reply := TaskAssignArgs{}, TaskAssignReply{}
		if ok := call("Coordinator.Example", &args, &reply); !ok {
			fmt.Printf("call failed!\n")
			// If the worker fails to contact the coordinator,
			// it can assume that the coordinator has exited
			// because the job is done, so the worker can terminate too.
			break
		}
		taskType, taskName, taskNum, nReduce := reply.TaskType, reply.TaskName, reply.TaskNum, reply.nReduce
		switch taskType {
		case "map":
			if err := doMap(mapf, taskName, taskNum, nReduce); err != nil {
				fmt.Printf("worker.map: %v", err)
				break workerLoop
			}
		case "reduce":
			if err := doReduce(reducef, taskName, taskNum); err != nil {
				fmt.Printf("worker.reduce: %v", err)
				break workerLoop
			}
		case "exit":
			break workerLoop
		default:
			fmt.Printf("worker: invalid task name\n")
			if err := handleInvalidTask(taskType, taskName, taskNum); err != nil {
				break workerLoop
			}
		}
		time.Sleep(time.Second)
	}
	fmt.Printf("worker exit\n")
	// uncomment to send the Example RPC to the coordinator.
	// CallExample()
}

func doMap(
	mapf func(string, string) []KeyValue,
	filename string,
	taskNum int,
	nReduce int,
) error {
	// mapf
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("cannot open %v", filename)
		if err := handleInvalidTask("map", filename, taskNum); err != nil {
			return err
		}
		return err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("cannot read %v", filename)
		if err := handleInvalidTask("map", filename, taskNum); err != nil {
			return err
		}
		return err
	}
	kva := mapf(filename, string(content))
	sort.Sort(ByKey(kva))

	// Create nReduce intermediate files
	intermediates := make([][]KeyValue, nReduce)
	for _, kv := range kva {
		reduceTaskNum := ihash(kv.Key) % nReduce
		intermediates[reduceTaskNum] = append(intermediates[reduceTaskNum], kv)
	}
	reduceTasks := make([]string, nReduce)
	prefix := "mr-" + strconv.Itoa(taskNum)
	for reduceTaskNum, intermediate := range intermediates {
		// Intermediate filename should be mr-XY.
		// X is map task number, Y is reduce task number.
		filename := prefix + strconv.Itoa(reduceTaskNum)
		if _, err := os.Stat(filename); err == nil {
			// File already exists, skip writing
			reduceTasks[reduceTaskNum] = filename
			continue
		}
		if err := atomicWriteFile(filename, intermediate, 0666, "json"); err != nil {
			fmt.Printf("Error: %v\n", err)
			if err := handleInvalidTask("map", filename, taskNum); err != nil {
				return err
			}
		} else {
			reduceTasks[reduceTaskNum] = filename
		}
	}
	// Some files may be sent to the coordinator in an empty state,
	// but since those tasks are reported to the coordinator
	// as deterministically unexecutable,
	// the coordinator will skip those tasks.
	args, reply := NewReduceTaskArgs{reduceTasks}, NewReduceTaskReply{}
	if ok := call("Coordinator.NewReduceTask", &args, &reply); !ok {
		return fmt.Errorf("doMap: call failed")
	}
	return nil
}

func doReduce(
	reducef func(string, []string) string,
	taskName string,
	taskNum int,
) error {
	file, err := os.Open(taskName)
	if err != nil {
		return err
	}
	var kva []KeyValue
	dec := json.NewDecoder(file)
	for {
		var kv KeyValue
		if err := dec.Decode(&kv); err != nil {
			break
		}
		kva = append(kva, kv)
	}
	keyToValues := make(map[string][]string)
	for _, kv := range kva {
		keyToValues[kv.Key] = append(keyToValues[kv.Key], kv.Value)
	}

	koa := make([]KeyValue, 0, len(keyToValues))
	for k, v := range keyToValues {
		output := reducef(k, v)
		koa = append(koa, KeyValue{k, output})
	}
	filename := "mr-out-" + strconv.Itoa(taskNum)
	if err := atomicWriteFile(filename, koa, 0666, "text"); err != nil {
		fmt.Printf("Error: %v\n", err)
		if err := handleInvalidTask("reduce", taskName, taskNum); err != nil {
			return err
		}
		return err
	}
	args, reply := ReduceTaskDoneArgs{taskName, taskNum}, ReduceTaskDoneReply{}
	if ok := call("Coordinator.ReduceTaskDone", &args, &reply); !ok {
		return fmt.Errorf("doReduce: call failed")
	}
	return nil
}

// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}

// If a task deterministically fails and cannot be executed,
// notify the coordinator so that it can skip the task.
func handleInvalidTask(taskType, taskName string, taskNum int) error {
	var err error
	args := InvalidTaskArgs{taskType, taskName, taskNum}
	reply := InvalidTaskReply{}
	if ok := call("Coordinator.InvalidTask", &args, &reply); !ok {
		fmt.Printf("call failed!\n")
		err = fmt.Errorf("call failed")
	}
	return err
}

// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

// To prevent a partially written file from being used
// if a crash occurs during writing,
// use os.CreateTemp to create a temporary file,
// write to it, and after the write is complete,
// atomically rename it using os.Rename.
func atomicWriteFile(
	filename string,
	kva []KeyValue,
	perm os.FileMode,
	mode string,
) error {
	tmpFile, err := os.CreateTemp(".", "tmp_"+filename+"_")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer func() {
		if tmpFile != nil {
			tmpFile.Close()
			os.Remove(tmpPath)
		}
	}()
	switch mode {
	case "json":
		// Write key/value pairs as JSON
		enc := json.NewEncoder(tmpFile)
		for _, kv := range kva {
			if err := enc.Encode(&kv); err != nil {
				return fmt.Errorf("failed to encode kv to json: %w", err)
			}
		}
	case "text":
		// Write key/value pairs as "%v %v\n"
		for _, kv := range kva {
			if _, err := fmt.Fprintf(tmpFile, "%v %v\n", kv.Key, kv.Value); err != nil {
				return fmt.Errorf("failed to write kv to text: %w", err)
			}
		}
	default:
		return fmt.Errorf("invalid mode: %v", mode)
	}
	// Ensure the file is flushed and written to disk
	if err := tmpFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}
	if err := os.Chmod(tmpPath, perm); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}
	// Atomically rename the temporary file
	if err := os.Rename(tmpPath, filename); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}
	tmpFile = nil // Prevent deletion by defer on success
	return nil
}

// example function to show how to make an RPC call to the coordinator.
//
// the RPC argument and reply types are defined in rpc.go.
func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	// the "Coordinator.Example" tells the
	// receiving server that we'd like to call
	// the Example() method of struct Coordinator.
	ok := call("Coordinator.Example", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("reply.Y %v\n", reply.Y)
	} else {
		fmt.Printf("call failed!\n")
	}
}
