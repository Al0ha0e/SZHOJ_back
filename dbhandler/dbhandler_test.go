package dbhandler

import (
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
	tags := []Tag{Tag{"a"},Tag{"b"}}
	qs := Question{
		Name:        "B+C Problem",
		Creator:     1001,
		Difficulty:  1,
		TotalCount:  10,
		AcceptCount: 9,
		TimeLimit:   1000,
		MemoryLimit: 1024,
		Tags: ,
	}
	hd := GetDBHandler()
	err := hd.InitDBHandler()
	if err != nil {
		t.Error(err)
	}
	hd.AddQuestion(&qs)
}

func TestAddStatus(t *testing.T) {
	st := Status{
		QuestionID:    1001,
		UserID:        50,
		CommitTime:    time.Now(),
		State:         1,
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

func TestGetQuestions(t *testing.T) {
	hd := GetDBHandler()
	err := hd.InitDBHandler()
	if err != nil {
		t.Error(err)
	}
	sq := hd.GetQuestions(&Question{})
	t.Log((*sq)[0].Name)
}

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
