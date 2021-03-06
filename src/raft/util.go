package raft

import (
	"fmt"
	"math/rand"
	"time"
)


func (rf *Raft) AssertTrue(test bool, format string, a ...interface{}) {
	if !test {
		panic(fmt.Sprintf(fmt.Sprintf("S%d ", rf.me)+format, a...))
	}
}

func getRandomElectionTimeout() time.Duration {
	return time.Duration(rand.Int())%RANDOM_PLUS + ELECTION_TIMEOUT
}


func (rf *Raft) checkTimeOut() bool {
	rf.timerMuLock.Lock()
	defer rf.timerMuLock.Unlock()
	return rf.timeToTimeOut.Before(time.Now())
}

func (rf *Raft) sleepTimeOut() {
	rf.timerMuLock.Lock()
	timeToTimeOut := rf.timeToTimeOut
	rf.timerMuLock.Unlock()	
	now := time.Now()
	if (timeToTimeOut.After(now)){
			time.Sleep(timeToTimeOut.Sub(now))

	}
}
func (rf *Raft) freshTimer() {
	rf.timerMuLock.Lock()
	defer rf.timerMuLock.Unlock()
	rf.timeToTimeOut = time.Now().Add(getRandomElectionTimeout())
}
func (rf *Raft) initTimeOut(){
	rf.timerMuLock.Lock()
	defer rf.timerMuLock.Unlock()
	totalTime := time.Duration(rand.Int())%RANDOM_PLUS + time.Millisecond*100
	rf.timeToTimeOut = time.Now().Add(totalTime)

	// rf.logger.Log(raftlogs.DTimer, 
	// 	"S%d just re-initialized the timeout from the current time %d to the time %d for a total timeout of %d",
	//     rf.me, time.Now().UnixNano()/1e6, rf.timeToTimeOut.UnixNano()/1e6, totalTime)
}


func checkCandidatesLogIsNew(selfTerm int, otherTerm int, selfIdx int, otherIdx int) bool{
	if (selfTerm != otherTerm){
		return otherTerm > selfTerm
	}
	return otherIdx >= selfIdx
}


//RPCs


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
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

func (rf *Raft) sendAppendEntry(server int, args *AppendEntryArgs, reply *AppendEntryReply) bool {
	ok := rf.peers[server].Call("Raft.AppendEntry", args, reply)
	return ok
}

func (rf *Raft) sendInstallSnapshot(server int, args *InstallSnapshotArgs, reply *InstallSnapshotReply) bool {
	ok := rf.peers[server].Call("Raft.InstallSnapshot", args, reply)
	return ok
}