# Lecture 5

The Raft paper describes a consensus algorithm, including many details that are needed to build replication-state machine applications. The paper is also the topic of several of the 6.5840 labs. The important sections are 2, 5, 7, and 8.

The paper positions itself as a better Paxos, but another way to look at Raft is that it solves a bigger problem than Paxos. To build a real-world replicated service, the replicas need to agree on an indefinite sequence of values (the client commands), and they need ways to efficiently recover when servers crash and restart or miss messages. People have built such systems with Paxos as the starting point (e.g., Google's Chubby and Paxos Made Live papers, and ZooKeeper/ZAB). There is also a protocol called Viewstamped Replication; it's a good design, and similar to Raft, but the paper about it is hard to understand.

These real-world protocols are complex, and (before Raft) there was not a good introductory paper describing how they work. The Raft paper, in contrast, is relatively easy to read and fairly detailed.

Question: Suppose we have the scenario shown in the Raft paper's Figure 7: a cluster of seven servers, with the log contents shown. The first server crashes (the one at the top of the figure), and cannot be contacted. A leader election ensues. For each of the servers marked (a), (d), and (f), could that server be elected? If yes, which servers would vote for it? If no, what specific Raft mechanism(s) would prevent it from being elected?
