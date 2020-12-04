/************
SZHOJ　V１.0.0 后端
由孙梓涵编写
本页面用于处理任务调度
************/

package scheduler

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/Al0ha0e/SZHOJ_back/dbhandler"
)

const TIME_LIMIT = 100

type taskInfo struct {
	status    *dbhandler.Status
	next      *taskInfo
	startTime time.Time
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
	taskTicker      *time.Ticker
	masterServer    *net.UDPConn
	backendDB       *dbhandler.DBHandler
	timeLimit       int
	//masterClient *net.UDPConn
}

//GetScheduler Get a Scheduler instance
func GetScheduler() *Scheduler {
	return &Scheduler{}
}

//Init init
func (sch *Scheduler) Init(back *dbhandler.DBHandler) (err error) {
	sch.timeLimit = TIME_LIMIT
	sch.taskTicker = time.NewTicker(time.Second)
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

//调度器UDP服务器，死循环等待消息
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

//处理服务端任务提交
func (sch *Scheduler) handleTaskCommit(task *dbhandler.Status) {
	fmt.Println("COMMIT TASK")
	tInfo := &taskInfo{
		status: task,
	}
	sch.pdQueueLock.Lock()
	//向链表中加入任务
	if sch.pendingTaskHead == nil && sch.pendingTaskTail == nil { //empty
		sch.pendingTaskHead = tInfo
		sch.pendingTaskTail = tInfo
	} else {
		sch.pendingTaskTail.next = tInfo
		sch.pendingTaskTail = tInfo
	}
	sch.pdQueueLock.Unlock()
}

//构造任务
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

//当Worker发送消息时被调用，维护Worker状态
func (sch *Scheduler) handleWorkerStatus(status *workerStatus) {
	workerID := status.worker.String()
	if status.QuestionID == 0 { //发送空状态作为心跳信号
		//heartbeat
		//fmt.Println("HEARTBEAT FROM", workerID)
		_, ok := sch.workers[workerID] //获取信号对应Worker
		if ok {
			//update worker
			if sch.pendingTaskHead != nil {
				//评测队列非空，取一个任务给Worker
				fmt.Println("NOT NULL")
				_, running := sch.runningTask[workerID]
				//Worker正在运行就不分配
				if running {
					fmt.Println("RUNNING")
					return
				}

				//对评测队列加锁
				sch.pdQueueLock.Lock()
				//从队列链表中去除任务
				task := sch.pendingTaskHead
				sch.pendingTaskHead = sch.pendingTaskHead.next
				if sch.pendingTaskHead == nil {
					sch.pendingTaskTail = nil
				}
				sch.pdQueueLock.Unlock()

				//绑定Worker与任务
				task.startTime = time.Now()
				sch.runningTask[workerID] = task
				job := sch.constructJob(task.status.QuestionID, task.status.Code)
				jobJSON, _ := json.Marshal(job)
				conn, _ := net.Dial("udp", workerID)
				defer conn.Close()
				conn.Write(jobJSON)
				//sch.masterServer.WriteoUDP(jobjson, fr)
			}
		} else {
			//此前没有Worker，将其加入Worker列表
			sch.workers[workerID] = &workerInfo{ip: workerID}
		}

	} else if status.QuestionID == -1 {
		//connect
		fmt.Println("CONNECT", workerID)
		sch.workers[workerID] = &workerInfo{ip: workerID}
	} else {
		//提交评测结果
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
		//将结果更新到数据库
		sch.backendDB.UpdataStatus(task.status)
	}
}

func (sch *Scheduler) handleTimerDispatch(currentTime time.Time) {
	for worker, task := range sch.runningTask {
		if task != nil {
			fmt.Println(currentTime.Second() - task.startTime.Second())
			if currentTime.Second()-task.startTime.Second() > sch.timeLimit {
				fmt.Println("SERVER TLE", worker)
				delete(sch.runningTask, worker)
				task.status.State = 6
				sch.backendDB.UpdataStatus(task.status)
				break
			}
		}
	}
}

//死循环，等待几种channel获得数据并调用对应的处理函数
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
		case currentTime := <-sch.taskTicker.C:
			sch.handleTimerDispatch(currentTime)
		}
	}
}
