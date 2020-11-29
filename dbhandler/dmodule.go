package dbhandler

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

//Question question structure
type Question struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	Name        string `gorm:"type:varchar(15);" json:"name"`
	Creator     uint   `json:"creator"`
	Difficulty  uint   `json:"difficulty"`
	TotalCount  uint   `json:"totCnt" `
	AcceptCount uint   `json:"acCnt" `
	TimeLimit   uint   `json:"timeLimit" `
	MemoryLimit uint   `json:"memoryLimit" `

	TotalStatus []Status `gorm:"ForeignKey:QuestionID" json:"-"`
	Tags        []Tag    `gorm:"many2many:question_tags;" json:"tags" `

	//for db create transaction
	Success bool       `gorm:"-" json:"-"`
	desc    *[]byte    `gorm:"-" json:"-"`
	code    *[]byte    `gorm:"-" json:"-"`
	data    *[]byte    `gorm:"-" json:"-"`
	db      *DBHandler `gorm:"-" json:"-"`
}

//PrepareForCreate set essential values before creation
func (q *Question) PrepareForCreation(db *DBHandler, desc *[]byte, code *[]byte, data *[]byte) {
	q.db = db
	q.desc = desc
	q.code = code
	q.data = data
}

//AfterCreate Call back after create
func (q *Question) AfterCreate() error {
	fmt.Println("AFTER")
	err := q.db.addQuestionFiles(q.ID, q.desc, q.code, q.data)
	if err != nil {
		q.Success = false
		return err
	}
	q.Success = true
	return nil
}

//User user structure
type User struct {
	gorm.Model
	Name               string      `gorm:"type:varchar(15);"`
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
	State         uint      `json:"state"`
	RunningTime   uint      `json:"time"`
	RunningMemory uint      `json:"memory"`

	TotalContestStatus []ContestStatus `gorm:"ForeignKey:StatusID" json:"-"`
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
	gorm.Model
	Name         string `gorm:"type:varchar(15);"`
	Creator      uint
	Contests     []Contest `gorm:"ForeignKey:UserGroupID"`
	SeenContests []Contest `gorm:"ForeignKey:Visibility"`
}

//Contest contest structure
type Contest struct {
	gorm.Model
	Creator     uint
	UserGroupID uint
	Visiblility uint
	StartTime   time.Time

	TotalContestStatus []ContestStatus `gorm:"ForeignKey:ContestID"`
	Questions          []Question      `gorm:"many2many:contest_questions;"`
}

//ContestStatus status of a contest
type ContestStatus struct {
	ID        uint
	ContestID uint
	StatusID  uint
}
