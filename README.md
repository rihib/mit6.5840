# MIT 6.5840 Distributed Systems Spring 2025

## Setup

```bash
brew install poppler
```

## Lectures

[6.5840 Schedule: Spring 2025](https://pdos.csail.mit.edu/6.824/schedule.html)
[Lab guidance](https://pdos.csail.mit.edu/6.824/labs/guidance.html)

- [ ] Lecture 1: Introduction
  - [x] [Lecture Note](lectures/01/l01.txt)
  - [x] [MapReduce: Simplified Data Processing on Large Clusters](lectures/01/mapreduce.pdf)
  - [x] [Video](https://youtu.be/WtZ7pcRSkOA?si=VU9nhFMlDNbbx08N)
  - [ ] [Lab1: MapReduce](https://pdos.csail.mit.edu/6.824/labs/lab-mr.html)
    - [x] Run & Read `src/main/mrsequential.go`
    - [x] Read `src/mrapps/wc.go`
    - [x] Read & Run `src/main/mrcoordinator.go` & `src/main/mrworker.go` & `src/mr/*.go`
      - How to run:
        1. Uncomment `CallExample()` in `src/mr/worker.go`
        2. Run below:

          ```bash
          go run mrcoordinator.go pg-*.txt

          # Open a new terminal
          go build -buildmode=plugin ../mrapps/wc.go
          go run mrworker.go wc.so
          ```

    - [x] Try to run `src/main/test-mr.sh`
    - [x] Implement `src/mr/*.go`
    - [ ] Run `mrcoordinator.go` & `mrworker.go`
      - How to run:

        ```bash
        cd src/main
        go get github.com/google/uuid
        go build -buildmode=plugin ../mrapps/wc.go
        rm mr-out*
        go run mrcoordinator.go pg-*.txt

        # Open a new terminal
        cd src/main
        go run mrworker.go wc.so
        cat mr-out-* | sort | more
        ```

    - [ ] Run `mrcoordinator.go` & `mrworker.go` with `-race` flag
    - [ ] Run `src/main/test-mr.sh`
    - [ ] Implement backup task
    - [ ] Run with `src/mrapps/crash.go`
    - [ ] Run `src/main/test-mr-many.sh`
    - [ ] No-credit challenge exercises
      - [ ] Implement your own MapReduce application
      - [ ] Get your MapReduce coordinator and workers to run on separate machines
- [x] Lecture 2: RPC and Threads
  - [x] [A Tour of Go Concurrency](https://go.dev/tour/concurrency/1)
  - [x] [Question](lectures/02/question.md) & Read `crawler.go`
  - [x] Read `kv.go`
  - [x] Read `condvar/vote-count-*.go`
    - [x] `go run -race vote-count-1.go`
  - [x] [Lecture Note](lectures/02/l-rpc.txt)
  - [x] [FAQ](lectures/02/tour-faq.txt)
  - [x] [Video](https://youtu.be/oZR76REwSyA?si=ujUaFr8AePOjSzWn)
- [ ] Lecture 3: Primary-Backup Replication
  - [x] [Lecture Note](lectures/03/l-vm-ft.txt)
  - [x] [FAQ](lectures/03/vm-ft-faq.txt)
  - [x] [The Design of a Practical System for Fault-Tolerant Virtual Machines](lectures/03/vm-ft.pdf)
  - [x] [Paper Questions](lectures/03/questions.md)
  - [x] [Video](https://youtu.be/gXiDmq1zDq4?si=vBWLws_WE0pgZZMF)
  - [ ] [Lab 2: Key/Value Server](https://pdos.csail.mit.edu/6.824/labs/lab-kvsrv1.html)
- [x] Lecture 4: Consistency and Linearizability
  - [x] [Lecture Note](lectures/04/l-linearizability.txt)
  - [x] [FAQ](lectures/04/linearizability-faq.txt)
  - [x] [Testing Distributed Systems for Linearizability](https://anishathalye.com/testing-distributed-systems-for-linearizability/)
  - [x] [Paper Questions](lectures/04/questions.md)
