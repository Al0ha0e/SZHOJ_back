package dbhandler

import (
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
)

func (hdl *DBHandler) addQuestionFiles(id uint, desc *[]byte, data *[]byte) error {
	strid := strconv.Itoa(int(id))
	batch := new(leveldb.Batch)
	batch.Put([]byte("desc_"+strid), *desc)
	batch.Put([]byte("data_"+strid), *data)
	err := hdl.kvDB.Write(batch, nil)
	return err
}

func (hdl *DBHandler) addCommitCode(id uint, code *[]byte) error {
	strid := strconv.Itoa(int(id))
	err := hdl.kvDB.Put([]byte("code_"+strid), *code, nil)
	return err
}

//GetJudgeData Get data for judger
func (hdl *DBHandler) GetJudgeData(id uint) (*[]byte, error) {
	strid := strconv.Itoa(int(id))
	data, err := hdl.kvDB.Get([]byte("data_"+strid), nil)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

//GetQuestionDesc get description
func (hdl *DBHandler) GetQuestionDesc(id uint64) ([]byte, error) {
	strid := strconv.Itoa(int(id))
	desc, err := hdl.kvDB.Get([]byte("desc_"+strid), nil)
	if err != nil {
		return nil, err
	}
	return desc, nil
}

// //AddKV Add question/code values
// func (hdl *DBHandler) AddKV(key []byte, value []byte) error {
// 	err := hdl.kvDB.Put(key, value, nil)
// 	return err
// }
