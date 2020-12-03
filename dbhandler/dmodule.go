package dbhandler

import (
	"fmt"
	"time"
)

//Question question structure
type Question struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	Name        string `gorm:"type:varchar(15);" json:"name"`
	Creator     uint   `json:"creator"`
	Difficulty  uint   `json:"difficulty"`
	TotalCount  uint   `json:"totCnt"`
	AcceptCount uint   `json:"acCnt"`
	TimeLimit   uint   `json:"timeLimit"`
	MemoryLimit uint   `json:"memoryLimit"`
	ContestID   uint   `json:"cid"`

	TotalStatus []Status `gorm:"ForeignKey:QuestionID" json:"-"`
	Tags        []Tag    `gorm:"many2many:question_tags;" json:"tags"`

	//for db create transaction
	Success bool       `gorm:"-" json:"-"`
	desc    *[]byte    `gorm:"-" json:"-"`
	datain  *[]byte    `gorm:"-" json:"-"`
	dataout *[]byte    `gorm:"-" json:"-"`
	db      *DBHandler `gorm:"-" json:"-"`
}

//PrepareForCreation set essential values before creation
func (q *Question) PrepareForCreation(db *DBHandler, desc *[]byte, datain *[]byte, dataout *[]byte) {
	q.db = db
	q.desc = desc
	q.datain = datain
	q.dataout = dataout
}

//AfterCreate Call back after create
func (q *Question) AfterCreate() error {
	fmt.Println("AFTER")
	err := q.db.addQuestionFiles(q.ID, q.desc, q.datain, q.dataout)
	if err != nil {
		q.Success = false
		return err
	}
	q.Success = true
	return nil
}

//User user structure
type User struct {
	ID                 uint        `gorm:"primary_key" json:"id"`
	Name               string      `gorm:"type:varchar(15);unique" json:"name"`
	Password           string      `gorm:"type:varchar(15);" json:"password"`
	TotalStatus        []Status    `gorm:"ForeignKey:UserID"`
	CreatedQuestions   []Question  `gorm:"ForeignKey:Creator"`
	AttendedUserGroups []UserGroup `gorm:"many2many:user_usergroups;"`
	CreatedUserGroups  []UserGroup `gorm:"ForeignKey:Creator"`
	CreatedContests    []Contest   `gorm:"ForeignKey:Creator"`
}

//Status status structure
type Status struct {
	ID            uint      `json:"id"`
	QuestionID    uint      `json:"qid"`
	UserID        uint      `json:"uid"`
	CommitTime    time.Time `json:"commitTime"`
	State         int       `json:"state"`
	RunningTime   uint      `json:"time"`
	RunningMemory uint      `json:"memory"`

	Code *[]byte    `gorm:"-" json:"-"`
	db   *DBHandler `gorm:"-" json:"-"`
}

//PrepareForCreation set essential values before creation
func (s *Status) PrepareForCreation(db *DBHandler, code *[]byte) {
	s.db = db
	s.Code = code
}

//AfterCreate Call back after create
func (s *Status) AfterCreate() error {
	fmt.Println("STATUS SQL OK", string(*s.Code))
	err := s.db.addCommitCode(s.ID, s.Code)
	return err
}

//MiniStatus for single user
type MiniStatus struct {
	ID         uint `json:"id"`
	QuestionID uint `json:"qid"`
	State      uint `json:"state"`
}

//Tag tag structure
type Tag struct {
	//	ID   uint   `json:"-"`
	Name string `gorm:"type:varchar(10);primary_key" binding:"required" json:"name"`
}

//UserGroup usergroup structure
type UserGroup struct {
	ID       uint      `gorm:"primary_key" json:"id"`
	Name     string    `gorm:"type:varchar(15);" json:"name"`
	Creator  uint      `json:"creator"`
	Users    []User    `gorm:"many2many:user_usergroups;" json:"users"`
	Contests []Contest `gorm:"ForeignKey:UserGroupID" json:"-"`
}

//Contest contest structure
type Contest struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	Name        string    `json:"name"`
	Creator     uint      `json:"creator"`
	UserGroupID uint      `json:"gid"`
	StartTime   time.Time `json:"start"`
	Duration    uint      `json:"duration"`

	TotalContestStatus []ContestStatus `gorm:"ForeignKey:ContestID" json:"-"`
	Questions          []Question      `gorm:"ForeignKey:ContestID" json:"questions"`
}

//ContestStatus status of a contest
type ContestStatus struct {
	ID            uint      `gorm:"primary_key" json:"id"`
	ContestID     uint      `json:"cid"`
	QuestionID    uint      `json:"qid"`
	UserID        uint      `json:"uid"`
	CommitTime    time.Time `json:"commitTime"`
	State         int       `json:"state"`
	RunningTime   uint      `json:"time"`
	RunningMemory uint      `json:"memory"`
	Code          *[]byte   `gorm:"-" json:"-"`
}
