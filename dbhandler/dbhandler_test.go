package dbhandler

import (
	"fmt"
	"testing"
	"time"
)

func TestInitDBHandler(t *testing.T) {
	hd := GetDBHandler()
	err := hd.InitDBHandler()
	if err != nil {
		t.Error("TEST INIT FAIL")
		t.Error(err)
	}
}

func TestAddQuestion(t *testing.T) {
	tags := []Tag{Tag{"a"}, Tag{"c"}}
	qs := Question{
		Name:        "中文测试",
		Creator:     1001,
		Difficulty:  10,
		TotalCount:  10,
		AcceptCount: 9,
		TimeLimit:   1000,
		MemoryLimit: 1024,
		Tags:        tags,
	}
	hd := GetDBHandler()
	err := hd.InitDBHandler()
	if err != nil {
		t.Error(err)
	}
	id := hd.AddQuestion(&qs)
	if id == 0 {
		t.Error("BAD ID")
	}
	fmt.Println(id)
}

func TestAddStatus(t *testing.T) {
	st := Status{
		QuestionID:    1,
		UserID:        10,
		CommitTime:    time.Now(),
		State:         2,
		RunningTime:   800,
		RunningMemory: 500,
	}
	hd := GetDBHandler()
	err := hd.InitDBHandler()
	if err != nil {
		t.Error(err)
	}
	hd.AddStatus(&st)
}

func TestGetQuestionByID(t *testing.T) {
	hd := GetDBHandler()
	err := hd.InitDBHandler()
	if err != nil {
		t.Error(err)
	}
	sq := hd.GetQuestionByID(1)
	if sq == nil {
		t.Error("NOT FOUND")
	}
	t.Log(sq.Name)
}

// func TestGetQuestions(t *testing.T) {
// 	hd := GetDBHandler()
// 	err := hd.InitDBHandler()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	cond := &Question{
// 		//Name: "中文测试",
// 		//Difficulty: 10,
// 		Tags: []Tag{Tag{Name: "c"}, Tag{Name: "b"}},
// 	}
// 	sq := hd.GetQuestions(cond)
// 	for _, val := range *sq {
// 		fmt.Println(val.ID)
// 	}
// 	//t.Log((*sq)[0].Name)
// }

func TestGetQuestionsByPage(t *testing.T) {
	hd := GetDBHandler()
	err := hd.InitDBHandler()
	if err != nil {
		t.Error(err)
	}
	sq := hd.GetQuestionsByPage(1, 2)
	t.Log((*sq)[0].ID, (*sq)[1].ID)
	sq = hd.GetQuestionsByPage(2, 2)
	t.Log((*sq)[0].ID)
	sq = hd.GetQuestionsByPage(3, 1)
	t.Log((*sq)[0].Name)
}

func TestGetStatus(t *testing.T) {
	hd := GetDBHandler()
	err := hd.InitDBHandler()
	if err != nil {
		t.Error(err)
	}
	cond := &Status{
		QuestionID: 104,
		UserID:     10,
	}
	sq := hd.GetStatus(cond)
	for _, val := range *sq {
		fmt.Println(val.ID, val.QuestionID, val.UserID)
	}
	//t.Log((*sq)[0].Name)
}

func TestGetMiniStatus(t *testing.T) {
	hd := GetDBHandler()
	err := hd.InitDBHandler()
	if err != nil {
		t.Error(err)
	}
	msq := hd.GetMiniStatus(10)
	for _, val := range *msq {
		fmt.Println(val.ID, val.QuestionID, val.State)
	}
}

func TestAddUserGroup(t *testing.T) {
	hd := GetDBHandler()
	err := hd.InitDBHandler()
	if err != nil {
		t.Error(err)
	}
	g := &UserGroup{
		Name: "test_group",
		Users: []User{
			User{ID: 1, Name: "sa"},
			User{ID: 2, Name: "sb"},
		},
	}
	hd.AddUserGroup(g)
}
