package scheduler

import (
	"encoding/json"
	"net"
	"sync"

	"github.com/Al0ha0e/SZHOJ_back/backserver"
	"github.com/Al0ha0e/SZHOJ_back/dbhandler"
)

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
	UserCode    string `json:"ucode"`
	TimeLimit   string `json:"timeLimit"`
	MemoryLimit string `json:"memoryLimit"`
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
	CommitChan      chan *dbhandler.Status
	udpChan         chan *workerStatus
	workers         map[string]*workerInfo
	pdQueueLock     sync.Mutex
	pendingTaskHead *taskInfo
	pendingTaskTail *taskInfo
	runningTask     map[string]*taskInfo
	masterServer    *net.UDPConn
	backend         *backserver.BackServer
	//masterClient *net.UDPConn
}

//GetScheduler Get a Scheduler instance
func GetScheduler() *Scheduler {
	return &Scheduler{}
}

//Init init
func (sch *Scheduler) Init(back *backserver.BackServer) (err error) {
	sch.backend = back
	sch.CommitChan = make(chan *dbhandler.Status, 10)
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

func (sch *Scheduler) handleTaskCommit(task *dbhandler.Status) {
	tInfo := &taskInfo{
		status: task,
	}
	sch.pdQueueLock.Lock()
	if sch.pendingTaskHead == nil && sch.pendingTaskTail == nil { //empty
		sch.pendingTaskHead = tInfo
		sch.pendingTaskTail = tInfo
	} else {
		sch.pendingTaskTail.next = tInfo
		sch.pendingTaskTail = tInfo
	}
	sch.pdQueueLock.Unlock()
}

func (sch *Scheduler) handleWorkerStatus(status *workerStatus) {
	workerID := status.worker.String()
	if status.QuestionID == 0 {
		//heartbeat
		_, ok := sch.workers[workerID]
		if ok {
			//update worker
			if sch.pendingTaskHead != nil {
				_, running := sch.runningTask[workerID]
				if running {
					return
				}
				sch.pdQueueLock.Lock()
				task := sch.pendingTaskHead
				sch.pendingTaskHead = sch.pendingTaskHead.next
				sch.pdQueueLock.Unlock()
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
		sch.backend.AddJudgeOutcome(task.status)
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
