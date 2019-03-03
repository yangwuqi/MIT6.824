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
type Entry struct{
	Log_Term int
	Log_Command interface{}
}


type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me int                 // this peer's index into peers[]

	state int//0表示follower，1表示candidate，2表示leader
	term int//当前的任期号
	//leaderID int//当前的leader的下标编号
	last_heartbeat chan int
	leaderID int
	//voted_time int64//记录上次投票的时间，避免同时一票多投，int64格式，单位毫秒
	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.
	applych chan ApplyMsg


	//chan_committed chan bool
	log []Entry//
	//chan_of_commit chan bool
	committed_index int
	last_applied int

	Last_log_index int//大索引
	//Last_log_index_index int//小索引
	Last_log_term int//最后一个log条例属于的term

	append_try_log_index int//用于在append失败的时候和leader协调重试，这个是大索引
	//append_try_index_index int//用于在append失败的时候和leader协调重试，这个是小索引
	append_try_log_term int//

	next_index []int
	match_index []int

	voted []int
	log_added bool
	log_added_content Entry
	Log_Term []int//一个直观和log索引和其term的对应关系数组

	last_term_log_lenth int
	will_commit bool
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
	From int
	Next_Log int//
	Append bool
	Newest_log Entry
	Newest_log_index int
	//Log_Uapdate_small []interface{}//用于小索引
	//Log_Uapdate []Entry//用于小索引后的大索引
	Leader_commited int
	Prev_log_index int
	Prev_log_term int
	//Prev_log_index_index int
	Append_Try bool//这个用于表示是不是在处于双方协商append过程中

	Second_position int
	Second_log []Entry
	Second bool

	Commit_MSG bool
	Commit_Log []Entry

	Last_log_term_lenth int
	Log_Term []int

	Last_log_term int
}

//
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	// Your data here (2A).
	Vote_sent bool//用于表示收到投票请求的机器是否同意投票
	//voted chan bool
	Append_success bool//log是否成功添加
	Try_log_index int//如果没有成功添加，返回不匹配的log索引和term给leader，让leader重试，这个是大索引
	Try_log_term int//这个是相应索引的log的term
	Log_Term []int//直接把机器上的所有log索引和其term的关系回复过去，一次就能找到一致的点
	Term int
	Commit_finished bool
	Last_log_lenth int
	Last_log_term int
	Wrong_leader bool
	Voted bool//用于查看这个机器在这个term有没有投过票
	You_are_true bool//投票时，如果发出投票请求的一方的term小于这个机器的term，但是log更加新，那么说明发出投票请求的一方才是真的更加新，应该把其term更新到最高
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
	rf.mu.Lock()
	/*
		if args.Append==true&&((args.Newest_log.Log_Term<rf.Last_log_term)||(args.Newest_log.Log_Term==rf.Last_log_term&&args.Last_log_term_lenth<rf.Last_log_term)){
			reply.Term=args.Candidate_term+1
			reply.Last_log_term=rf.Last_log_term
			reply.Last_log_lenth=rf.last_term_log_lenth
			reply.Append_success=false
			rf.mu.Unlock()
			return
		}
	*/
	//if args.Second==true{
	//	fmt.Printf("!\n!\n!\n!\n!\n编号为%d的raft实例收到编号为%d的leader的second请求！本机term是%d，leader term是%d，args.Append是%v\n",rf.me,args.From,rf.term,args.Candidate_term,args.Append)
	//}

	if rf.state==2&&((rf.term<args.Candidate_term)||(rf.term==args.Candidate_term&&args.Last_log_term<rf.Last_log_term))&&args.Votemsg==false{
		//fmt.Printf("分区恢复后编号为%d的raft实例的term是%d，发现自己已经不是leader！leader是%d，leader的term是%d\n",rf.me,rf.term,args.From,args.Candidate_term)
		rf.state=0
		rf.leaderID=args.From
	}



	if args.Candidate_term>=rf.term{
		//rf.term=args.Candidate_term
		//if args.Second==true{
		//	fmt.Printf("服务器上的SECOND进入第一个大括号\n")
		//}
		if args.Append == false {
			if args.Votemsg == true && rf.voted[args.Candidate_term] == 0&&((args.Last_log_term>rf.Last_log_term)||(args.Last_log_term==rf.Last_log_term&&args.Last_log_term_lenth>=rf.last_term_log_lenth))  { //合法投票请求
				//fmt.Printf("编号为%d的raft实例对投票请求的回答为true,term统一更新为为%d\n",rf.me,rf.term)

				//rf.term = args.Candidate_term
				rf.voted[args.Candidate_term] = 1
				reply.Vote_sent = true

				//rf.voted_time=time.Now().UnixNano()/1e6

			}else if args.Votemsg==true{ //合法的纯heartbeat
			if rf.voted[args.Candidate_term]==1 {
				reply.Voted = true
			}
			//fmt.Printf("请求方的term是%d，本机的term是%d，来自%d的投票请求被%d拒绝！rf.last_log_term是%d，rf.last_log_lenth是%d，本机的rf.last_log_term是%d，rf.last_log_lenth是%d\n",args.Candidate_term,rf.term,args.From,rf.me,args.Last_log_term,args.Last_log_term_lenth,rf.Last_log_term,rf.last_term_log_lenth)
			}
			reply.Term=rf.term

			//rf.term=args.Candidate_term//!!!!!!!!!!!!!!!!!!!!!!!!!!!
			//if args.Votemsg==true{//!!!!!!!!!!!!!!
			//	rf.term=args.Candidate_term//!!!!!!!!!!!!
			//}//!!!!!!!!!!!!!!!!!

		} else { //这条是关于日志的
			//这个请求是日志同步请求，接收方需要将自己的日志最后一条和leader发过来的声称的进行比较，如果leader的更新且leader的PREV和自己的LAST相同就接受
			//还得找到最后一个一致的日志位置，然后将后面的全部更新为和leader一致的，这意味着中间多次的RPC通信

			/*
			if args.Newest_log.Log_Term<rf.Last_log_term{
				reply.Wrong_leader=true
				reply.Term=rf.term
				reply.Append_success=false
				reply.Last_log_lenth=rf.last_term_log_lenth
				return
			}
*/

if (rf.Last_log_term>args.Last_log_term)||(rf.Last_log_term==args.Last_log_term&&rf.last_term_log_lenth>args.Last_log_term_lenth){
	reply.Append_success=false
	reply.Last_log_term=rf.Last_log_term
	reply.Last_log_lenth=rf.last_term_log_lenth
	rf.mu.Unlock()
	return
}


			rf.term=args.Candidate_term
			if args.Second==true{
				//	fmt.Printf("在服务器端进入second阶段！\n")
				rf.log=rf.log[:args.Second_position]
				rf.log=append(rf.log,args.Second_log...)
				reply.Append_success=true
				rf.Last_log_term=args.Last_log_term
				rf.last_term_log_lenth=args.Last_log_term_lenth
				rf.Last_log_index=len(rf.log)-1
				rf.Log_Term=args.Log_Term
				//fmt.Printf("Second APPend在服务器端成功！现在编号为%d的raft实例的log是%v, last_log_term是%d，term是%d\n",rf.me,rf.log,rf.Last_log_term,rf.term)
			}else{
				if args.Append_Try == false {//try用于表示是否是第一次append失败了现在正在沟通
					rf.append_try_log_index = rf.Last_log_index
					rf.append_try_log_term=rf.Last_log_term
				}
				if args.Prev_log_index != rf.append_try_log_index || args.Prev_log_term != rf.append_try_log_term{
					//fmt.Printf("匹配失败！！！%d号leader发过来的PREV_log_index是%d,本机%d的last_log_index是%d，PREV_term是%d，本机的last_log_term是%d！\n",args.From,args.Prev_log_index,rf.me,rf.append_try_log_index,args.Prev_log_term,rf.append_try_log_term)
					reply.Vote_sent = false//匹配失败后进入双方沟通try
					reply.Append_success = false

					reply.Log_Term=rf.Log_Term

					rf.mu.Unlock()
					return
				} else { //说明没问题。可以更新
					//fmt.Printf("匹配成功！！！%d号是leader，发过来的PREV_log_index是%d,本机的last_log_index是%d，PREV_term是%d，本机的last_log_term是%d，准备更新本机日志！！\n", args.From, args.Prev_log_index, rf.append_try_log_index, args.Prev_log_term, rf.append_try_log_term)
					//rf.Last_log_term = args.Last_log_term
					rf.last_term_log_lenth=args.Last_log_term_lenth
					rf.log = append(rf.log, args.Newest_log)
					rf.Last_log_index += 1
					rf.Log_Term = args.Log_Term
					rf.Last_log_term=args.Newest_log.Log_Term
					reply.Append_success = true
					//fmt.Printf("APPend成功，现在编号为%d的raft实例的log是%v，last_log_term是%d，term是%d\n",rf.me,rf.log,rf.Last_log_term,rf.term)
				}
			}
			rf.log_added_content = args.Newest_log
			rf.last_term_log_lenth=0

			for cc:=len(rf.log)-1;cc>-1;cc--{
				if rf.log[cc].Log_Term!=rf.Last_log_term{
					break
				}
				rf.last_term_log_lenth+=1
			}


		}

		//fmt.Printf("在更新heartbeat之前\n")
		if args.Votemsg==false {//加上个约束条件更严谨，加上了表示是在heartbeat开始之后认同了这个是leader，否则在投票阶段就认同了
//fmt.Printf("rf.last_log_term %d, args.last_log_term %d\n",rf.Last_log_term,args.Last_log_term)
if args.Last_log_term==rf.Last_log_term {//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	if args.Commit_MSG == true {
		//if  len(rf.Log_Term)==len(args.Log_Term)&&rf.Log_Term[len(rf.Log_Term)-1]==args.Log_Term[len(args.Log_Term)-1]{
		//if len(args.Log_Term)==len(rf.Log_Term)&&args.Last_log_term==rf.Last_log_term {
		for cc := rf.committed_index + 1; cc <= rf.Last_log_index; cc++ {
			rf.committed_index = cc
			//fmt.Printf("在follower %d 上进行commit，commit_index是%d，commit的内容是%v，commit的term是%d，last_log_term是%d, rf.log是太长暂时鸽了\n", rf.me, cc, rf.log[cc].Log_Command, rf.log[cc].Log_Term, rf.Last_log_term)
			rf.applych <- ApplyMsg{true, rf.log[rf.committed_index].Log_Command, rf.committed_index}
		}

		reply.Commit_finished = true
		//}else{
		//}
		//}
	}
}//!!!!!!!!!!!!!!!!!!!!!!!!!!!!

			rf.leaderID = args.From
			rf.term = args.Candidate_term
			rf.leaderID=args.From


		}
		reply.Last_log_lenth=rf.last_term_log_lenth
		reply.Last_log_term=rf.Last_log_term

		if args.Votemsg==false {
			if rf.state == 0 {
				rf.last_heartbeat <- 1
			}
		}

	}else{
		//fmt.Printf("term都不符，明显是非法的！\n")
		reply.Vote_sent = false
		reply.Append_success = false
		reply.Term=rf.term
		reply.Last_log_lenth=rf.last_term_log_lenth
		reply.Last_log_term=rf.Last_log_term
		if (args.Last_log_term>rf.Last_log_term)||(args.Last_log_term==rf.Last_log_term&&args.Last_log_term_lenth>=rf.last_term_log_lenth){
			reply.You_are_true=true
		}
	}
	rf.mu.Unlock()
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

func(rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply,signal chan bool) bool {
	//if args.Second==true {
	//fmt.Printf("开始Second调用sendRequestVote!在sendRequestVote可以看到args.Second的内容是%v，接续postion是%d\n",args.Second_log,args.Second_position)
	//}
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	//if ok&&args.Second==true {
	//	fmt.Printf("调用成功Second调用sendRequestVote!\n")
	//}
	//if !ok&&args.Second==true{
	//	fmt.Printf("Second RPC调用失败！\n")
	//}
	if ok{
		signal<-true
	}
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

type Start_leader struct {
	Index int
	Term int
	Leader bool
}
func (rf *Raft)Start_to_leader(command interface{},reply *Start_leader){
	reply.Index,reply.Term,reply.Leader=rf.Start(command)
}

func (rf *Raft) Start(command interface{}) (int, int, bool) {//这个是客户端向服务器请求加log或者说指令的函数
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if rf.leaderID==-1{
		//fmt.Printf("尝试对编号为%d的raft实例调用start，但是当前还没有leader\n",rf.me)
		return -1,-1,false
	}
	isLeader := false
	//xx:=rf.Last_log_index+1
	// Your code here (2B).

	if rf.state==2{
		isLeader=true


		rf.Last_log_term=rf.term

		newlog:=Entry{rf.term,command}
		rf.log=append(rf.log,newlog)
		rf.log_added=true
		rf.log_added_content=newlog
		rf.Last_log_index+=1
		rf.Log_Term=append(rf.Log_Term,rf.term)
		rf.last_term_log_lenth=0
		for cc:=len(rf.log)-1;cc>-1;cc--{
			if rf.log[cc].Log_Term!=rf.Last_log_term{
				break
			}
			rf.last_term_log_lenth+=1
		}

		//fmt.Printf("编号为%d的raft实例是leader，其term是%d, 调用了Start，现在leader的log是太长了先不输出, term是%d，rf.last_log_term %d\n",rf.me,rf.term,rf.term,rf.Last_log_term)
	}else{
		//fmt.Printf("客户端当前请求的编号为%d的raft实例不是leader，将任务转交给leader %d\n",rf.me,rf.leaderID)
		var reply Start_leader
		signal:=make(chan bool,1)
		go rf.send_to_leader(signal,command,&reply)
		select {
		case <-signal:
			//fmt.Printf("已经转交并获得返回！index %d, term %d, leader %v\n",reply.Index,reply.Term,reply.Leader)
			return reply.Index,reply.Term,reply.Leader
		case <-time.After(time.Duration(11)*time.Millisecond):
			//fmt.Printf("转交超时！\n")
		}
	}//这个raft实例是leader，然后log被添加之后，需要通过heartbeat机制复制给其它节点，大多数人都复制了之后commit
	return rf.Last_log_index,rf.term, isLeader
}

func (rf *Raft)send_to_leader(signal chan bool,command interface{},reply *Start_leader){
	if rf.peers[rf.leaderID].Call("Raft.Start_to_leader",command,reply){
		signal<-true
	}
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
	//rf.term=0
	rf.me=me
	rf.state=0//日踏大爷就是因为这个channel没有make 1卡了老子好久，不make一个数的话就在读或者写之后如果在本线程内不能写或者读就会阻塞
	rf.last_heartbeat=make(chan int,1)//这里一定要注意！有缓冲的与有缓冲channel有着重大差别！有缓冲可以是异步的而无缓冲只能是同步的
	rf.voted=make([]int,20000)
	rf.log=[]Entry{}
	rf.Last_log_index=-1
	rf.Last_log_term=-1
	rf.readPersist(persister.ReadRaftState())
	rf.leaderID=-1
	rf.applych=applyCh
	rf.committed_index=-1
	//rf.chan_of_commit=make(chan bool,1)
	go rf.running()
	//go rf.commit(applyCh)
	return rf
}

func(rf *Raft)running(){
	for {
BEGIN:
		//fmt.Printf("这个raft实例的编号是%d，状态是%d，term是%d，实例总数是%d\n",rf.me,rf.state,rf.term,len(rf.peers))
		var args RequestVoteArgs
/*
		rf.mu.Lock()
		args.Candidate_term=rf.term
		args.From=rf.me
		args.Leader_commited=rf.committed_index
		args.Last_log_term_lenth=rf.last_term_log_lenth
		args.Last_log_term=rf.Last_log_term
		args.Log_Term=rf.Log_Term

		xx:=rf.state
		rf.mu.Unlock()
*/
		switch rf.state {

		case 0: //如果状态是follower
			select {
			case <-rf.last_heartbeat:
				//fmt.Printf("编号为%d是raft实例读取了heartbeat\n",rf.me)
			case <-time.After(time.Duration(rand.Intn(666)+300) * time.Millisecond):
				//fmt.Printf("经过等待后超时，编号%d的raft实例成为candidate,其term是%d, log是%v\n",rf.me,rf.term+1,rf.log)
				rf.mu.Lock()
			if rf.state==0 {
				//rf.state = 1 //复习一下select的用法，这里有两种可能，被heartbeat维持或者超时进入candidate
				rf.term += 1
			}
			//----------------------------------------

				votes:=0
				//----------------------------rf.mu.Lock()

				args.Candidate_term=rf.term
				args.From=rf.me
				args.Leader_commited=rf.committed_index
				args.Last_log_term_lenth=rf.last_term_log_lenth
				args.Last_log_term=rf.Last_log_term
				args.Log_Term=rf.Log_Term

				//---------------------------------------if rf.state==1 {
					args.Newest_log_index = rf.Last_log_index
					if len(rf.log) > 0 {
						args.Newest_log = rf.log[len(rf.log)-1]
					} else {

					}
					if rf.voted[rf.term] == 0 {
						rf.voted[rf.term] = 1
						votes++
					}
					//---------------------------rf.mu.Unlock()
					args.Votemsg = true
					for i := 0; i < len(rf.peers); i++ {
						if i != rf.me {
							var reply RequestVoteReply
							//fmt.Printf("编号为%d的raft实例向编号为%d的raft实例请求投票\n",rf.me,i)
							signal := make(chan bool, 1)
							go rf.sendRequestVote(i, &args, &reply, signal)
							select {
							case <-signal:
								if reply.Vote_sent == true {
									votes++
									//fmt.Printf("编号为%d的raft实例收到了来自编号为%d的raft实例的投票，现在有%d票\n",rf.me,i,votes)
								}else{
									//fmt.Printf("本机term是%d，对方term是%d，args.Append %v, args.voteMSG %v，编号为%d的raft实例的投票请求被编号为%d的raft实例拒绝啦！现在有%d票。本机rf.last_log_term是%d，rf.last_term_log_lenth是%d，对方rf.last_log_term是%d，rf.last_term_log_lenth是%d，对方reply.voted==%v\n",rf.term,reply.Term,args.Append,args.Votemsg,rf.me,i,votes,rf.Last_log_term,rf.last_term_log_lenth,reply.Last_log_term,reply.Last_log_lenth,reply.Voted)
									if rf.term<reply.Term&&reply.You_are_true{//!!!!!!!!!!!!!!!
										rf.term=reply.Term+6//!!!!!!!!!!!!!!这里多加点
									}//!!!!!!!!!!!!!!
							}
							case <-time.After(time.Duration(10) * time.Millisecond):
								//fmt.Printf("编号为%d的raft实例对编号为%d的raft实例的RPC调用的返回超时啦！\n",rf.me,i)
							}
						}
					}
					//------------------------------rf.mu.Lock()
					if votes > len(rf.peers)/2 {
						rf.leaderID = rf.me
						//fmt.Printf("len(rf.peers)/2=%d, 下标编号为%d的raft实例收到%d个投票，成为leader, rf.term是%d，rf.last_log_term是%d\n", len(rf.peers)/2, rf.me, votes, rf.term, rf.Last_log_term)
						rf.state = 2
					} else {
						rf.state = 0
						//rf.term-=1
					}
				//--------------------------------------}
				//---------------------------------------rf.mu.Unlock()

			//----------------------------------------
				rf.mu.Unlock()
			}
		//--------------------case 1: //如果状态是candidate

		case 2: //如果状态是leader

				append_ok := 1
				select {
				//case <-rf.last_heartbeat:
				//	rf.mu.Lock()
				//	rf.state=0
				//	rf.mu.Unlock()
				case <-time.After(time.Duration(50) * time.Millisecond):

					rf.mu.Lock()

					args.Candidate_term=rf.term
					args.From=rf.me
					args.Leader_commited=rf.committed_index
					args.Last_log_term_lenth=rf.last_term_log_lenth
					args.Last_log_term=rf.Last_log_term
					args.Log_Term=rf.Log_Term

					if rf.state==2 {

					//------------------------------rf.mu.Lock()
					if rf.will_commit == true {
						args.Commit_MSG = true
						rf.will_commit = false
					}
					if rf.log_added == true {
						args.Append = true
						args.Prev_log_term = rf.Last_log_term
						args.Newest_log = rf.log_added_content

						args.Prev_log_term = -1
						args.Prev_log_index = rf.Last_log_index - 1
						if args.Prev_log_index >= 0 {
							args.Prev_log_term = rf.log[args.Prev_log_index].Log_Term
						} //不然就全是-1，设想第一个log条例是0,0坐标

					}
					for i := 0; i < len(rf.peers); i++ {
						if i != rf.me {
							var reply RequestVoteReply
							var reply1 RequestVoteReply
							//fmt.Printf("下标编号为%d的raft实例向全网发送heartbeat\n",rf.me)
							signal := make(chan bool, 1)
							go rf.sendRequestVote(i, &args, &reply, signal)
							select {
							case <-signal:
								if reply.Append_success == true {
									//fmt.Printf("Append成功！\n")
									append_ok++
								} else if args.Append == true {
									//fmt.Printf("对方的term是%d,这个leader %d的term是%d，append失败！\n", reply.Term, rf.me,rf.term)


										//fmt.Printf("对方的term是%d，这个leader的term是%d，append失败！\n", reply.Term, rf.term)
										//对方是分区回归的follower，直接rf.term=reply.Term，如果本机是分区回归的leader，则本机变为follower，通过日志来判断
										if (reply.Last_log_term > rf.Last_log_term) || (reply.Last_log_term == rf.Last_log_term && reply.Last_log_lenth > rf.last_term_log_lenth) { //这时肯定是本机是分区返回的leader，本机直接变为follower
											rf.state = 0
											//fmt.Printf("本机%d是一个老leader，过时了，变成follower\n",rf.me)
											rf.mu.Unlock()
											goto BEGIN
											//break
										} else {
											if reply.Term > rf.term {
												//fmt.Printf("发现一个超前的follower，本机%d的term增加到%d\n", rf.me, reply.Term) //注意append那边的log相关的term更新全是log里面的，不然增加term之后在同一轮append的log的term都不一样
												rf.term = reply.Term
												//goto BEGIN
												//break
											}
										}

									/*
								if reply.Wrong_leader==true {
									//if reply.Term >= rf.term {
										//rf.term = reply.Term
										rf.state = 0
										break
									//}
								}
								*/
									signal1 := make(chan bool, 1)
									var args1 RequestVoteArgs
									if args.Commit_MSG == true {
										args1.Commit_MSG = true
									}
									args1.Candidate_term = rf.term
									args1.From = rf.me
									args1.Append = true
									args1.Second = true
									args1.Newest_log = rf.log_added_content
									args1.Log_Term = rf.Log_Term
									args1.Last_log_term = rf.Last_log_term
									args1.Last_log_term_lenth = rf.last_term_log_lenth
									//fmt.Printf("leader是%d，编号为%d的raft实例返回的LOG_TERM是%v，leader的LOG_TERM是%v\n",rf.me,i,reply.Log_Term,rf.Log_Term)
									if len(reply.Log_Term) > 0 {
										ii := len(reply.Log_Term) - 1
										for ii > 0 && (ii >= len(rf.log) || (ii < len(rf.log) && reply.Log_Term[ii] != rf.log[ii].Log_Term)) {
											ii--
										}
										args1.Second_position = ii
										for j := ii; j < len(rf.log); j++ {
											args1.Second_log = append(args1.Second_log, rf.log[j])
										}
										//fmt.Printf("服务器短的log_term为%v，leader的log_term,是%v，最后相同的log下标是%d，准备发送%v过去\n",reply.Log_Term,rf.Log_Term,ii,args1.Second_log)
									} else {
										args1.Second_position = 0
										for j := 0; j < len(rf.log); j++ {
											args1.Second_log = append(args1.Second_log, rf.log[j])
										}
										//fmt.Printf("服务器短的log为空，leader的log是%v，准备发送%v过去\n",rf.log,args1.Second_log)
									}
									go rf.sendRequestVote(i, &args1, &reply1, signal1)
									select {
									case <-signal1:
										if reply1.Append_success == true {
											//fmt.Printf("Second Append成功！\n")
											append_ok++
										}
									case <-time.After(time.Duration(12) * time.Millisecond):
										//fmt.Printf("Second Append超时！\n")
									}
								}
								//fmt.Printf("编号为%d的raft实例对编号为%d的raft实例发送heartbeat成功\n",rf.me,i)
							case <-time.After(time.Duration(22) * time.Millisecond):
								//fmt.Printf("编号为%d的raft实例对编号为%d的raft实例发送heartbeat无应答！\n",rf.me,i)
							}
						}
					}

					rf.log_added = false
					if append_ok > (len(rf.peers) / 2) { //说明可以commit了
						for cc := rf.committed_index + 1; cc <= rf.Last_log_index; cc++ {
							rf.committed_index = cc
							//fmt.Printf("在leader%d进行commit！committed_index是%d，commit的内容是%v，commit的term是%d，last_log_term是%d, rf.log是太长暂时鸽了\n", rf.me, cc, rf.log[cc].Log_Command, rf.log[cc].Log_Term, rf.Last_log_term)
							rf.applych <- ApplyMsg{true, rf.log[rf.committed_index].Log_Command, rf.committed_index}
						}
						rf.will_commit = true
					}
					//---------------------------rf.mu.Unlock()
				}
					rf.mu.Unlock()
			}

		}
	}
}
