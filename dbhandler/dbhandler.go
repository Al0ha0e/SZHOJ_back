package dbhandler

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql driver
)

const dbstr = "ojtest:weakpassword@(127.0.0.1:3307)/oj?charset=utf8&parseTime=True&loc=Local"

//DBHandler handler db operations SQL and leveldb
type DBHandler struct {
	sqlDB *gorm.DB
}

//GetDBHandler Get a instance of DBHandler
func GetDBHandler() *DBHandler {
	return &DBHandler{}
}

//InitDBHandler Init DBHandler
func (this *DBHandler) InitDBHandler() (err error) {
	print("DSDIU")
	this.sqlDB, err = gorm.Open("mysql", dbstr)
	print("DSDIU2")
	if err != nil {
		return err
	}
	if !this.sqlDB.HasTable(&Question{}) {
		this.sqlDB.CreateTable(&Question{})
	}
	if !this.sqlDB.HasTable(&User{}) {
		this.sqlDB.CreateTable(&User{})
	}
	if !this.sqlDB.HasTable(&Status{}) {
		this.sqlDB.CreateTable(&Status{})
	}
	if !this.sqlDB.HasTable(&Tag{}) {
		this.sqlDB.CreateTable(&Tag{})
	}
	if !this.sqlDB.HasTable(&UserGroup{}) {
		this.sqlDB.CreateTable(&UserGroup{})
	}
	if !this.sqlDB.HasTable(&Contest{}) {
		this.sqlDB.CreateTable(&Contest{})
	}
	if !this.sqlDB.HasTable(&ContestStatus{}) {
		this.sqlDB.CreateTable(&ContestStatus{})
	}
	return nil
}

func (this *DBHandler) AddQuestion(q *Question) {
	this.sqlDB.Create(q)
}

func (this *DBHandler) AddStatus(s *Status) {
	this.sqlDB.Create(s)
}

// func (this *DBHandler) GetQuestion(info *Question) error {

// 	return nil
// }

//GetQuestionsByPage Get All questions
func (this *DBHandler) GetQuestionsByPage(pageNum uint64, itemPerPage uint64) *[]Question {
	ret := make([]Question, 0)
	st := (pageNum - 1) * itemPerPage
	en := pageNum * itemPerPage
	this.sqlDB.Where("ID > ? AND ID <= ?", st, en).Find(&ret)
	return &ret
}

//GetQuestions Get All questions fulfill the conditions given by info
func (this *DBHandler) GetQuestions(info *Question) *[]Question {
	ret := make([]Question, 0)
	this.sqlDB.Find(&ret)
	fmt.Println(len(ret))
	return &ret
}

//GetStatusByPage Get All status
func (this *DBHandler) GetStatusByPage(pageNum uint, itemPerPage uint) *[]Status {
	ret := make([]Status, 0)
	st := (pageNum - 1) * itemPerPage
	en := pageNum * itemPerPage
	this.sqlDB.Where("ID > ? AND ID <= ?", st, en).Find(&ret)
	return &ret
}

//GetStatus Get All status fulfill the conditions given by info
func (this *DBHandler) GetStatus(info *Status) *[]Status {
	ret := make([]Status, 0)
	this.sqlDB.Find(&ret)
	fmt.Println(len(ret))
	return &ret
}
