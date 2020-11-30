package scheduler

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"sync"

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
	DataIn      string `json:"datain"`
	DataOut     string `json:"dataout"`
	UserCode    string `json:"ucode"`
	TimeLimit   string `json:"timeLimit"`
	MemoryLimit string `json:"memoryLimit"`
}

type workerStatus struct {
	QuestionID int          `json:"qid"`
	Time       uint         `json:"time"`
	Memory     uint         `json:"memory"`
	State      int          `json:"state"`
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
	backendDB       *dbhandler.DBHandler
	//masterClient *net.UDPConn
}

//GetScheduler Get a Scheduler instance
func GetScheduler() *Scheduler {
	return &Scheduler{}
}

//Init init
func (sch *Scheduler) Init(back *dbhandler.DBHandler) (err error) {
	sch.workers = make(map[string]*workerInfo)
	sch.runningTask = make(map[string]*taskInfo)
	sch.backendDB = back
	sch.CommitChan = make(chan *dbhandler.Status, 10)
	sch.udpChan = make(chan *workerStatus, 10)
	address, err := net.ResolveUDPAddr("udp", ":8040")
	if err != nil {
		return err
	}
	sch.masterServer, err = net.ListenUDP("udp", address)
	if err != nil {
		return err
	}
	return nil
}

//Start start a node scheduler
func (sch *Scheduler) Start() {
	fmt.Println("SCH START")
	go sch.serve()
	go sch.serveUDP()
}

func (sch *Scheduler) informServer(status *workerStatus) {
	sch.udpChan <- status
}

func (sch *Scheduler) serveUDP() {
	fmt.Println("START  SERVERUDP")
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
	fmt.Println("COMMIT TASK")
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

func (sch *Scheduler) constructJob(qid uint, code *[]byte) *Job {
	job := &Job{QuestionID: qid, UserCode: string(*code)}
	datain, dataout, _ := sch.backendDB.GetJudgeData(qid)
	job.DataIn = string(*datain)
	job.DataOut = string(*dataout)
	question := sch.backendDB.GetQuestionByID(uint64(qid))
	job.TimeLimit = strconv.Itoa(int(question.TimeLimit))
	job.MemoryLimit = strconv.Itoa(int(question.MemoryLimit))
	return job
}

func (sch *Scheduler) handleWorkerStatus(status *workerStatus) {
	workerID := status.worker.String()
	if status.QuestionID == 0 {
		//heartbeat
		fmt.Println("HEARTBEAT FROM", workerID)
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
				job := sch.constructJob(task.status.QuestionID, task.status.Code)
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
		fmt.Println("CONNECT", workerID)
		sch.workers[workerID] = &workerInfo{ip: workerID}
	} else {
		task, ok := sch.runningTask[workerID]
		fmt.Println("COMMIT", workerID, status.QuestionID, status.State, status.Time, status.Memory)
		if !ok || task.status.QuestionID != uint(status.QuestionID) {
			fmt.Println("BAD")
			return
		}
		delete(sch.runningTask, workerID)
		task.status.RunningTime = status.Time
		task.status.RunningMemory = status.Memory
		task.status.State = status.State
		sch.backendDB.UpdataStatus(task.status)
	}
}

func (sch *Scheduler) serve() {
	fmt.Println("START SERVE")
	for {
		select {
		case task := <-sch.CommitChan:
			fmt.Println("BACK COMMIT")
			sch.handleTaskCommit(task)
		case status := <-sch.udpChan:
			fmt.Println("UDP COMMIT")
			sch.handleWorkerStatus(status)
		}
	}
}
