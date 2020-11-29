package scheduler

import (
	"encoding/json"
	"net"

	"github.com/Al0ha0e/SZHOJ_back/backserver"
	"github.com/Al0ha0e/SZHOJ_back/dbhandler"
)

//TaskToCommit strcture to server
type TaskToCommit struct {
	Status *dbhandler.Status
	Code   *[]byte
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
	QuestionID int          `json:"qid"`
	Time       uint         `json:"time"`
	Memory     uint         `json:"memory"`
	State      uint         `json:"state"`
	worker     *net.UDPAddr `json:"-"`
}

//Scheduler schedule judger nodes
type Scheduler struct {
	CommitChan   chan *TaskToCommit
	udpChan      chan *workerStatus
	workers      map[string]*workerInfo
	pendingTask  *taskInfo
	runningTask  map[string]*taskInfo
	masterServer *net.UDPConn
	backend      *backserver.BackServer
	//masterClient *net.UDPConn
}

//GetScheduler Get a Scheduler instance
func GetScheduler() *Scheduler {
	return &Scheduler{}
}

//Init init
func (sch *Scheduler) Init() (err error) {
	sch.CommitChan = make(chan *TaskToCommit, 10)
	sch.udpChan = make(chan *workerStatus, 10)
	address, _ := net.ResolveUDPAddr("udp", ":8040")
	sch.masterServer, err = net.ListenUDP("udp", address)
	if err != nil {
		return err
	}
	return nil
}

//Start start a node scheduler
func (sch *Scheduler) Start() {
	go sch.serve()
	go sch.serveUDP()
}

func (sch *Scheduler) informServer(status *workerStatus) {
	sch.udpChan <- status
}

func (sch *Scheduler) serveUDP() {

	defer sch.masterServer.Close()
	for {
		data := make([]byte, 4*1024*1024)
		l, fr, _ := sch.masterServer.ReadFromUDP(data)
		denseData := data[:l]
		status := &workerStatus{worker: fr}
		json.Unmarshal(denseData, status)
		go sch.informServer(status)
	}
}

func (sch *Scheduler) handleTaskCommit(task *TaskToCommit) {

}

func (sch *Scheduler) handleWorkerStatus(status *workerStatus) {
	workerID := status.worker.String()
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
				job := Job{QuestionID: task.status.ID}
				jobJSON, _ := json.Marshal(job)
				conn, _ := net.Dial("udp", workerID)
				defer conn.Close()
				conn.Write(jobJSON)
				//sch.masterServer.WriteoUDP(jobjson, fr)
			}
		} else {
			sch.workers[workerID] = &workerInfo{ip: workerID}
		}

	} else if status.QuestionID == -1 {
		//connect
		sch.workers[workerID] = &workerInfo{ip: workerID}
	} else {
		task, ok := sch.runningTask[workerID]
		if !ok {
			return
		}
		delete(sch.runningTask, workerID)
		task.status.RunningTime = status.Time
		task.status.RunningMemory = status.Memory
		task.status.State = status.State
		//WRITE TO DB
	}
}

func (sch *Scheduler) serve() {
	for {
		select {
		case task := <-sch.CommitChan:
			sch.handleTaskCommit(task)
		case status := <-sch.udpChan:
			sch.handleWorkerStatus(status)
		}
	}
}
