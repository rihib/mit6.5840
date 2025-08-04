package mr

import (
	"fmt"
	"hash/fnv"
	"log"
	"net/rpc"
)

// Map functions return a slice of KeyValue.
type KeyValue struct {
	Key   string
	Value string
}

// main/mrworker.go calls this function.
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	// Your worker implementation here.
	// コーディネータにタスクを要求する RPC を送信

	// Map
	// coordinatorからファイル名を返されたら、それを入力としてMap関数を呼び出してnReduce個の中間ファイルを生成
	// workerは中間Map出力を現在のディレクトリ内のファイルに格納する必要がある
	// 中間ファイルの適切な命名規則はmr-XY（XはMapタスク番号、Yはreduceタスク番号）
	// Reduceタスク番号はihash(key) % nReduceで求められる
	// reduceタスクで中間ファイルの内容を読み取りやすくするために、encoding/jsonパッケージを使用してJSON形式でkey/valueを書き込む

	// mapタスクが全て完了するまでコーディネータに定期的にタスク要求を送信
	// mapタスクが全て完了するまでreduceタスクを開始できない
	// 各リクエストの間にtime.Sleep()でスリープする

	// Reduce
	// Reduce関数を呼び出して結果を生成
	// workerは中間ファイルをReduceタスクへの入力として読み取る
	// X番目のreduceタスクの出力をmr-out-Xファイルに書き込む
	// "%v %v\n"の形式でkeyとvalueを書き込む
	// 書き込み途中でクラッシュした際は書き込み途中のファイルが使われるのを防ぐために、os.CreateTempを使って一時ファイルを作成して書き込みを行い、書き込みが完了したらos.Renameを使用してアトミックに名前を変更する

	// ジョブが完全に終了したらworkerプロセスを終了させる
	// call()の戻り値を使って簡単に実装できる
	// ワーカーがコーディネータへの接続に失敗した場合、ジョブが完了したためコーディネータが終了したとみなし、ワーカーも終了できる
	// 設計によっては、コーディネータがワーカーに「終了してください」という疑似タスクを与えることも有効

	// uncomment to send the Example RPC to the coordinator.
	// CallExample()
}

// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
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
