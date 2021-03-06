package raftlogs

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type LogTopic int

const (
	RAFT_IGNORE    = true
	COMMIT_IGNORE  = false //apply,commit
	TIMER_IGNORE   = true  //timer
	LEADER_IGNORE  = true  //leader
	APPEND_IGNORE  = true  //append
	ROLE_IGNORE    = false //term vote
	PERSIST_IGNORE = false //persist log

	//raft
	Start LogTopic = iota
	Commit
	Drop
	Error
	Info
	Leader
	LogModify
	Persist
	Snap
	Term
	Test
	Timer
	Trace
	Vote
	Warn
	Log2
	RaftShutdown
	Role
	Apply
	Append
	SnapSize

	//kv server
	ServerReq = 100
	ServerSnap
	ServerApply
	ServerSnapSize
	ServerShutdown
	ServerStart
	ServerConfig
	ServerMove

	//kv clerk
	Clerk

	//test config
	Cfg

	//ctrler
	CtrlerStart = 200
	CtrlerQuery
	CtrlerJoin
	CtrlerLeave
	CtrlerMove
	CtrlerBalance
	CtrlerApply
	CtrlerSnap
	CtrlerReq
	CtrlerSnapSize

	//shardkv
	ShardKVReq = 500
	ShardKVStart
	ShardKVApply
	ShardKVConfig
	ShardKVMigration
	ShardKVSnap
	ShardKVSnapSize
	ShardKVShutDown

	Debug_Level = 2
)

var debugStart time.Time
var debugVerbosity int
var Print_Map map[LogTopic]bool

func getVerbosity() int {
	v := os.Getenv("VERBOSE")
	level := 0
	if v != "" {
		var err error
		level, err = strconv.Atoi(v)
		if err != nil {
			log.Fatalf("Invalid verbosity %v", v)
		}
	}
	return level
}
func init() {
	debugVerbosity = getVerbosity()
	debugStart = time.Now()

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	Print_Map = make(map[LogTopic]bool)
	print_list := []LogTopic{
		// Commit,
		// Drop,
		// Error,
		// Info,
		// Leader,
		// LogModify,
		// Persist,
		// Snap,
		// Term,
		// Test,
		// Timer,
		// Trace,
		// Vote,
		// Warn,
		// Log2,
		// RaftShutdown,
		// Role,
		// Apply,
		// Append,
		// SnapSize,
		// ServerSnap,
		// ServerApply,
		// ServerSnapSize,
		// ServerShutdown,
		// ServerStart,
		// ServerConfig,
		// ServerMove,
		// Log2,
		// Clerk,
		// ShardKVApply,
		// ShardKVConfig,
		// ShardKVMigration,
		// ShardKVReq,
		// ShardKVStart,
		//CtrlerApply,
	}

	for _, v := range print_list {
		Print_Map[v] = true
	}
}

type TopicLogger struct {
	Me    int
	Group int
}

func (tp *TopicLogger) L(topic LogTopic, format string, a ...interface{}) {
	if Print_Map[topic] {
		time := time.Since(debugStart).Milliseconds()
		time_seconds := time / 1000
		time = time % 1000
		prefix := ""
		if topic >= ShardKVReq {
			prefix = fmt.Sprintf("%03d'%03dms G%03d-S%d ", time_seconds, time, tp.Group, tp.Me)
		} else {
			prefix = fmt.Sprintf("%03d'%03dms S%d ", time_seconds, time, tp.Me)
		}
		format = prefix + format
		log.Printf(format, a...)
	}
}
