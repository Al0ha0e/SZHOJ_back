package scheduler

import (
	"encoding/json"
	"net"

	"github.com/Al0ha0e/SZHOJ_back/dbhandler"
)

type TaskToCommit struct {
	Status *dbhandler.Status
	Code *[]byte
}

type taskInfo struct {
	status *dbhandler.Status
	next   *taskInfo
}

type workerInfo struct {
	ip string
}

//Job job sent to worker
type Job struct {
	QuestionID  uint   `json:"qid"`
	Data        string `json:"data"`
	UserCode    string `json"ucode"`
	TimeLimit   string `timeLimit`
	MemoryLimit string `memoryLimit`
}

type workerStatus struct {
	QuestionID int  `json:"qid"`
	Time       uint `json:"time"`
	Memory     uint `json:"memory"`
	State      uint `json:"state"`
}

//Scheduler schedule judger nodes
type Scheduler struct {
	CommitChan chan TaskToCommit
	workers      map[string]*workerInfo
	pendingTask  *taskInfo
	runningTask  map[string]*taskInfo
	masterServer *net.UDPConn
}

//GetScheduler Get a Scheduler instance
func GetScheduler() *Scheduler {
	return &Scheduler{}
}

//Init init
func (sch *Scheduler) Init() (err error) {
	sch.CommitChan = make(chan TaskToCommit,10)
	address, _ := net.ResolveUDPAddr("udp", ":8040")
	sch.masterServer, err = net.ListenUDP("udp", address)
	if err != nil {
		return err
	}
}

//Start start a node scheduler
func (sch *Scheduler) Start() {
	go sch.serve()
}

func (sch *Scheduler) serve() {
	defer sch.masterServer.Close()
	for {
		data := make([]byte, 4*1024*1024)
		l, fr, _ := sch.masterServer.ReadFromUDP(data)
		denseData := data[:l]
		var status workerStatus
		json.Unmarshal(denseData, &status)
		workerID := fr.String()
		if status.QuestionID == 0 {
			//heartbeat
			_, ok := sch.workers[workerID]
			if ok {
				//update worker
				if sch.pendingTask != nil {
					_, running := sch.runningTask[workerID]
					if running {
						return
					}
					task := sch.pendingTask
					sch.pendingTask = sch.pendingTask.next
					sch.runningTask[workerID] = task
					job := Job{QuestionID: task.status.ID,Data: }
					jobjson, _ := json.Marshal(job)
					sch.masterServer.WriteToUDP(jobjson, fr)
				}
			} else {
				sch.workers[workerID] = &workerInfo{ip: workerID}
			}

		} else if status.QuestionID == -1 {
			//connect
			sch.workers[workerID] = &workerInfo{ip: workerID}
		} else {

		}
	}
}
