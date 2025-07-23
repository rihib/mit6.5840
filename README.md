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
  - [ ] [Video](https://youtu.be/oZR76REwSyA?si=ujUaFr8AePOjSzWn)
  - [ ] condvar以下のコードを読み、条件変数の使い方を理解する
    - `go run -race vote-count-1.go`
  - [ ] [Question](https://pdos.csail.mit.edu/6.824/questions.html?q=q-gointro&lec=2)
    - 課された課題をやる(クローラーの実装とRPCパッケージの理解)
    - Lecture Noteと一緒にcrawler.goとkv.goを読む
    - Sync.Mutex、sync.Cond、Sync.WaitGroupについて理解する
    - 再度FAQを読んでちゃんと全て理解する
