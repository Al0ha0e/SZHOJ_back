package dbhandler

import (
	"testing"
)

func TestInitDBHandler(t *testing.T) {
	hd := GetDBHandler()
	err := hd.InitDBHandler()
	if err != nil {
		t.Error("TEST INIT FAIL")
		t.Error(err)
	}
}

func TestAddUser(t *testing.T) {
	hd := GetDBHandler()
	err := hd.InitDBHandler()
	if err != nil {
		t.Error("TEST INIT FAIL")
		t.Error(err)
	}
	hd.AddUser("test1", "sdsdsad")
}

func TestAddQuestion(t *testing.T) {
	hd := GetDBHandler()
	err := hd.InitDBHandler()
	if err != nil {
		t.Error("TEST INIT FAIL")
		t.Error(err)
	}
	data := []byte("asdad")
	ques := &Question{
		Name:    "afdf",
		desc:    &data,
		datain:  &data,
		dataout: &data,
		db:      hd,
	}
	hd.AddQuestion(ques)
}

// func TestAddStatus(t *testing.T) {
// 	hd := GetDBHandler()
// 	err := hd.InitDBHandler()
// 	if err != nil {
// 		t.Error("TEST INIT FAIL")
// 		t.Error(err)
// 	}

// }
