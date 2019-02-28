package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import "sync"
import (
	"../labrpc"
	"time"
	"math/rand"
	"fmt"
)

// import "bytes"
// import "labgob"

//先实现选举投票，再实现日志

//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// in Lab 3 you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh; at that point you can add fields to
// ApplyMsg, but set CommandValid to false for these other uses.
//
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int
}

//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me int                 // this peer's index into peers[]

	state int//0表示follower，1表示candidate，2表示leader
	term int//当前的任期号
	leaderID int//当前的leader的下标编号
	last_heartbeat chan int
	voted_time int64//记录上次投票的时间，避免同时一票多投，int64格式，单位毫秒
	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.

}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {
	rf.mu.Lock()
	var term int
	var isleader bool
	term=rf.term
	if rf.state==2{
		isleader=true
	}else{
		isleader=false
	}
	rf.mu.Unlock()
	// Your code here (2A).
	return term, isleader
}


//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here (2C).
	// Example:
	// w := new(bytes.Buffer)
	// e := labgob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// data := w.Bytes()
	// rf.persister.SaveRaftState(data)
}


//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// Your code here (2C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := labgob.NewDecoder(r)
	// var xxx
	// var yyy
	// if d.Decode(&xxx) != nil ||
	//    d.Decode(&yyy) != nil {
	//   error...
	// } else {
	//   rf.xxx = xxx
	//   rf.yyy = yyy
	// }
}




//
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//
type RequestVoteArgs struct {
	Votemsg bool
	// Your data here (2A, 2B).
	Candidate_term int//candidate的term，用于判断请求投票是否有效
}

//
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	// Your data here (2A).
	Vote_sent bool//用于表示收到投票请求的机器是否同意投票
	//voted chan bool
}

//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	//fmt.Printf("成功调用RequestVote！\n")
	// Your code here (2A, 2B).
	//rf.mu.Lock()
	//current_time:=time.Now().UnixNano()/1e6
	//&&current_time-rf.voted_time>800
	if args.Candidate_term>=rf.term&&args.Votemsg==true{
		fmt.Printf("编号为%d的raft实例对投票请求的回答为true,term统一更新为为%d\n",rf.me,rf.term)
		rf.mu.Lock()
		rf.term=args.Candidate_term
		reply.Vote_sent=true
		rf.voted_time=time.Now().UnixNano()/1e6
		rf.mu.Unlock()
	}else {
		reply.Vote_sent = false

	}
//fmt.Printf("在更新heartbeat之前\n")
	rf.last_heartbeat<-1
	//fmt.Printf("编号为%d的raft实例通过RequestVote()收到了heartbeat\n",rf.me)
	//reply.voted<-true
	//rf.mu.Unlock()
}

//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	//fmt.Printf("成功调用sendRequestVote!\n")
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}


//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election. even if the Raft instance has been killed,
// this function should return gracefully.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true

	// Your code here (2B).


	return index, term, isLeader
}

//
// the tester calls Kill() when a Raft instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
func (rf *Raft) Kill() {
	// Your code here, if desired.
}

//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me
	rf.term=0
	rf.me=me
	rf.state=0//日踏大爷就是因为这个channel没有make 1卡了老子好久，不make一个数的话就在读或者写之后如果在本线程内不能写或者读就会阻塞
	rf.last_heartbeat=make(chan int,1)//这里一定要注意！有缓冲的与有缓冲channel有着重大差别！有缓冲可以是异步的而无缓冲只能是同步的

	rf.readPersist(persister.ReadRaftState())

	go rf.running()
	return rf
}

func(rf *Raft)running(){
	for {
		//fmt.Printf("这个raft实例的编号是%d，状态是%d，term是%d，实例总数是%d\n",rf.me,rf.state,rf.term,len(rf.peers))
		var args RequestVoteArgs
		var reply RequestVoteReply
		args.Candidate_term=rf.term
		reply.Vote_sent=false
		switch rf.state {
		case 0: //如果状态是follower
		//rf.mu.Lock()
		//fmt.Printf("进入第一个select\n")
			select {
			case <-rf.last_heartbeat:
				//fmt.Printf("编号为%d是raft实例读取了heartbeat\n",rf.me)
			case <-time.After(time.Duration(rand.Intn(666)+300) * time.Millisecond):
				fmt.Printf("经过等待后超时，编号%d的raft实例成为candidate\n",rf.me)
				rf.mu.Lock()
				rf.state=1//复习一下select的用法，这里有两种可能，被heartbeat维持或者超时进入candidate
				rf.term+=1
				rf.mu.Unlock()
			}
		//rf.mu.Unlock()
		case 1: //如果状态是candidate
			votes:=1
			args.Votemsg=true
			for i:=0;i<len(rf.peers);i++{
				if i!=rf.me{
					//fmt.Printf("编号为%d的raft实例向编号为%d的raft实例请求投票\n",rf.me,i)
					rf.sendRequestVote(i,&args,&reply)
					//<-reply.voted
					if reply.Vote_sent==true{
						fmt.Printf("编号为%d的raft实例收到了来自编号为%d的raft实例的投票\n",rf.me,i)
						votes++
					}
				}
			}
			rf.mu.Lock()
			if votes>len(rf.peers)/2{
				fmt.Printf("len(rf.peers)/2=%d, 下标编号为%d的raft实例收到%d个投票，成为leader\n",len(rf.peers)/2,rf.me,votes)
				rf.state=2
			}else{
				rf.state=0
			}
			rf.mu.Unlock()
		case 2: //如果状态是leader
			select {
			case <-rf.last_heartbeat:
				rf.mu.Lock()
				rf.state=0
				rf.mu.Unlock()
			case <-time.After(time.Duration(80)*time.Millisecond):
				for i:=0;i<len(rf.peers);i++{
					if i!=rf.me{
						//fmt.Printf("下标编号为%d的raft实例向全网发送heartbeat\n",rf.me)
						rf.sendRequestVote(i,&args,&reply)
					}
				}
			}
		}
	}
}