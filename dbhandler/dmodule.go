package dbhandler

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Question question structure
type Question struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	Name        string `gorm:"type:varchar(15);" json:"name"`
	Creator     uint   `json:"creator"`
	Difficulty  uint   `json:"difficulty" `
	TotalCount  uint   `json:"totCnt" `
	AcceptCount uint   `json:"acCnt" `
	TimeLimit   uint   `json:"timeLimit" `
	MemoryLimit uint   `json:"memoryLimit" `

	TotalStatus []Status `gorm:"ForeignKey:QuestionID" json:"-"`
	Tags        []Tag    `gorm:"many2many:question_tags;" json:"tags" `
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
	ID            uint
	QuestionID    uint
	UserID        uint
	CommitTime    time.Time
	State         uint
	RunningTime   uint
	RunningMemory uint

	TotalContestStatus []ContestStatus `gorm:"ForeignKey:StatusID"`
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
