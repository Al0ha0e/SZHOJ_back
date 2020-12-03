package dbhandler

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql driver
	"github.com/syndtr/goleveldb/leveldb"
)

const dbstr = "ojtest:weakpassword@(127.0.0.1:3307)/oj?charset=utf8&parseTime=True&loc=Local"

//DBHandler handler db operations SQL and leveldb
type DBHandler struct {
	sqlDB *gorm.DB
	kvDB  *leveldb.DB
}

//GetDBHandler Get a instance of DBHandler
func GetDBHandler() *DBHandler {
	return &DBHandler{}
}

//InitDBHandler Init DBHandler
func (hdl *DBHandler) InitDBHandler() (err error) {

	hdl.kvDB, err = leveldb.OpenFile("./db", nil)
	if err != nil {
		return err
	}

	hdl.sqlDB, err = gorm.Open("mysql", dbstr)
	if err != nil {
		return err
	}
	if !hdl.sqlDB.HasTable(&Question{}) {
		hdl.sqlDB.CreateTable(&Question{})
	}
	if !hdl.sqlDB.HasTable(&User{}) {
		hdl.sqlDB.CreateTable(&User{})
	}
	if !hdl.sqlDB.HasTable(&Status{}) {
		hdl.sqlDB.CreateTable(&Status{})
	}
	if !hdl.sqlDB.HasTable(&Tag{}) {
		hdl.sqlDB.CreateTable(&Tag{})
	}
	if !hdl.sqlDB.HasTable(&UserGroup{}) {
		hdl.sqlDB.CreateTable(&UserGroup{})
	}
	if !hdl.sqlDB.HasTable(&Contest{}) {
		hdl.sqlDB.CreateTable(&Contest{})
	}
	if !hdl.sqlDB.HasTable(&ContestStatus{}) {
		hdl.sqlDB.CreateTable(&ContestStatus{})
	}
	return nil
}

//Dispose Release Resources when shut down
func (hdl *DBHandler) Dispose() {
	hdl.sqlDB.Close()
	hdl.kvDB.Close()
}

//AddUser Add a user
func (hdl *DBHandler) AddUser(username, password string) {
	user := &User{
		Name:     username,
		Password: password,
	}
	hdl.sqlDB.Create(user)
}

//AddQuestion AddQuestion
func (hdl *DBHandler) AddQuestion(q *Question) uint {
	hdl.sqlDB.Create(q)
	return q.ID
}

//AddStatus AddStatus
func (hdl *DBHandler) AddStatus(s *Status) {
	hdl.sqlDB.Create(s)
}

//AddUserGroup add
func (hdl *DBHandler) AddUserGroup(g *UserGroup) {
	hdl.sqlDB.Create(g)
}

//AddContestStatus add
func (hdl *DBHandler) AddContestStatus(s *ContestStatus) {
	hdl.sqlDB.Create(s)
}

//AddContest Add a contest
func (hdl *DBHandler) AddContest(c *Contest) {
	hdl.sqlDB.Create(c)
}

//UpdateContestStatus asas
func (hdl *DBHandler) UpdateContestStatus(s *Status) {
	hdl.sqlDB.Save(s)
}

//UpdataStatus update status
func (hdl *DBHandler) UpdataStatus(s *Status) {
	question := hdl.GetQuestionByID(uint64(s.QuestionID))
	if question.ContestID != 0 {
		cStatus := ContestStatus{
			ContestID:     question.ContestID,
			QuestionID:    s.QuestionID,
			UserID:        s.UserID,
			CommitTime:    s.CommitTime,
			State:         s.State,
			RunningTime:   s.RunningTime,
			RunningMemory: s.RunningMemory,
		}
		hdl.sqlDB.Create(cStatus)
		hdl.sqlDB.Delete(s)
	} else {
		hdl.sqlDB.Save(s)
	}

}

//DeleteUserGroup del
func (hdl *DBHandler) DeleteUserGroup(gid uint64) {
	ug := &UserGroup{ID: uint(gid)}
	hdl.sqlDB.Delete(ug)
}

//GetUserByName get user by username
func (hdl *DBHandler) GetUserByName(name string) *User {
	var user User
	hdl.sqlDB.Where("name=?", name).First(&user)
	return &user
}

//GetUserByID get
func (hdl *DBHandler) GetUserByID(uid uint64) *User {
	var user User
	hdl.sqlDB.First(&user, uid)
	return &user
}

//GetQuestionCnt get count
func (hdl *DBHandler) GetQuestionCnt() uint {
	var ret uint
	hdl.sqlDB.Table("questions").Count(&ret)
	return ret
}

//GetQuestionsByPage Get All questions
func (hdl *DBHandler) GetQuestionsByPage(pageNum uint64, itemPerPage uint64) *[]Question {
	ret := make([]Question, 0)
	st := (pageNum - 1) * itemPerPage
	en := pageNum * itemPerPage
	hdl.sqlDB.Preload("Tags").Where("ID > ? AND ID <= ? AND ContestID = ?", st, en, 0).Find(&ret)
	return &ret
}

//GetQuestionByID get by id
func (hdl *DBHandler) GetQuestionByID(id uint64) *Question {
	ret := &Question{ID: uint(id)}
	if hdl.sqlDB.NewRecord(ret) {
		return nil
	}
	hdl.sqlDB.Preload("Tags").First(ret, id)
	return ret
}

//GetQuestions Get All questions fulfill the conditions given by info
func (hdl *DBHandler) GetQuestions(info *Question) *[]Question {
	ret := make([]Question, 0)
	query := hdl.sqlDB.Preload("Tags").Where("ContestID=?", 0)
	if len(info.Name) > 0 {
		query = query.Where("Name=?", info.Name)
	}
	if info.Difficulty > 0 {
		query = query.Where("Difficulty=?", info.Difficulty)
	}
	//Caution!!! Beacuse of the BAD BEHAVIOUR appears when OR operation is performed,
	//DO NOT QUERY OTHER CONDITIONS WITH TAGS!!!
	// if len(info.Tags) > 0 {
	// 	query = query.Joins("JOIN question_tags ON questions.id=question_tags.question_id")
	// 	query = query.Where("tag_name=?", info.Tags[0].Name)
	// 	for i := 1; i < len(info.Tags); i++ {
	// 		query = query.Or("tag_name=?", info.Tags[i].Name)
	// 	}
	// }
	query.Find(&ret)
	return &ret
}

//GetStatusByPage Get All status
func (hdl *DBHandler) GetStatusByPage(pageNum uint64, itemPerPage uint64) *[]Status {
	ret := make([]Status, 0)
	st := (pageNum - 1) * itemPerPage
	en := pageNum * itemPerPage
	hdl.sqlDB.Where("ID > ? AND ID <= ?", st, en).Order("ID DESC").Find(&ret)
	return &ret
}

//GetStatus Get All status fulfill the conditions given by info
func (hdl *DBHandler) GetStatus(info *Status) *[]Status {
	ret := make([]Status, 0)
	query := hdl.sqlDB
	if info.QuestionID > 0 {
		fmt.Println("NOT ZERO")
		query = query.Where("question_id=?", info.QuestionID)
	}
	if info.UserID > 0 {
		query = query.Where("user_id=?", info.UserID)
	}
	query.Order("ID DESC").Find(&ret)
	return &ret
}

//GetMiniStatus Get Mini mum Status For User
func (hdl *DBHandler) GetMiniStatus(userID uint64) *[]MiniStatus {
	ret := make([]MiniStatus, 0)
	hdl.sqlDB.Table("statuses").Where("user_id=?", userID).Select("id, question_id, state").Scan(&ret)
	return &ret
}

//GetStatusByID get status by id
func (hdl *DBHandler) GetStatusByID(sid uint64) *Status {
	ret := &Status{}
	hdl.sqlDB.First(ret, sid)
	return ret
}

//GetContestByPage get contest
func (hdl *DBHandler) GetContestByPage(pageNum uint64, itemPerPage uint64) *[]Contest {
	ret := make([]Contest, 0)
	st := (pageNum - 1) * itemPerPage
	en := pageNum * itemPerPage
	hdl.sqlDB.Where("ID > ? AND ID <= ?", st, en).Order("start_time DESC").Find(&ret)
	return &ret
}

//GetContestByID get contest
func (hdl *DBHandler) GetContestByID(cid uint64) *Contest {
	ret := &Contest{}
	hdl.sqlDB.Preload("Questions").First(ret, cid)
	return ret
}

//GetContestStatus get
func (hdl *DBHandler) GetContestStatus(cid uint64, uid uint64) *[]ContestStatus {
	ret := make([]ContestStatus, 0)
	hdl.sqlDB.Where("ContestID=? AND UserID=?", cid, uid).Find(&ret)
	return &ret
}

//GetTotalContestStatus get
func (hdl *DBHandler) GetTotalContestStatus(cid uint64) *[]ContestStatus {
	ret := make([]ContestStatus, 0)
	hdl.sqlDB.Where("ContestID=?", cid).Find(&ret)
	return &ret
}

//GetUserGroupByID get
func (hdl *DBHandler) GetUserGroupByID(gid uint) *UserGroup {
	ret := &UserGroup{}
	hdl.sqlDB.Preload("Users").First(ret, gid)
	return ret
}

//GetAttendedUserGroups attend
func (hdl *DBHandler) GetAttendedUserGroups(uid uint) *[]UserGroup {
	user := &User{ID: uid}
	ret := make([]UserGroup, 0)
	hdl.sqlDB.Preload("Users").Model(user).Association("AttendedUserGroups").Find(&ret)
	return &ret
}

//GetCreatedUserGroups attend
func (hdl *DBHandler) GetCreatedUserGroups(uid uint) *[]UserGroup {
	user := &User{ID: uid}
	ret := make([]UserGroup, 0)
	hdl.sqlDB.Preload("Users").Model(user).Association("CreatedUserGroups").Find(&ret)
	return &ret
}
