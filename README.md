# MIT 6.5840 Distributed Systems Spring 2025

[6.5840 Schedule: Spring 2025](https://pdos.csail.mit.edu/6.824/schedule.html)

## Lectures

- [ ] Lecture 1: Introduction
  - [x] [Lecture Note](lectures/01/l01.txt)
  - [x] [MapReduce: Simplified Data Processing on Large Clusters](lectures/01/mapreduce.pdf)
  - [x] [Video](https://youtu.be/WtZ7pcRSkOA?si=VU9nhFMlDNbbx08N)
  - [ ] [Lab1: MapReduce](https://pdos.csail.mit.edu/6.824/labs/lab-mr.html)
- [ ] Lecture 2: RPC and Threads
  - [x] [Lecture Note](lectures/02/l-rpc.txt)
  - [x] [FAQ](lectures/02/tour-faq.txt)
  - [x] [Video](https://youtu.be/oZR76REwSyA?si=ujUaFr8AePOjSzWn)
  - [ ] [Question](lectures/02/question.md)
    - condvar以下のコードを読み、条件変数の使い方を理解する＆自分でも書けるようにする
    - `go run -race vote-count-1.go`
    - 課された課題をやる(クローラーの実装とRPCパッケージの理解)
    - Lecture Noteと一緒にcrawler.goとkv.goを読む
    - 再度FAQを読んでちゃんと全て理解する
- [ ] Lecture 3: Primary-Backup Replication
  - [x] [Lecture Note](lectures/03/l-vm-ft.txt)
  - [x] [FAQ](lectures/03/vm-ft-faq.txt)
  - [x] [The Design of a Practical System for Fault-Tolerant Virtual Machines](lectures/03/vm-ft.pdf)
  - [x] [Paper Questions](lectures/03/questions.md)
  - [ ] [Video](https://youtu.be/gXiDmq1zDq4?si=vBWLws_WE0pgZZMF)
  - [ ] [Lab 2: Key/Value Server](https://pdos.csail.mit.edu/6.824/labs/lab-kvsrv1.html)
- [x] Lecture 4: Consistency and Linearizability
  - [x] [Lecture Note](lectures/04/l-linearizability.txt)
  - [x] [FAQ](lectures/04/linearizability-faq.txt)
  - [x] [Testing Distributed Systems for Linearizability](https://anishathalye.com/testing-distributed-systems-for-linearizability/)
  - [x] [Paper Questions](lectures/04/questions.md)

## Setup

```bash
brew install poppler
```
