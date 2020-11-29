package dbhandler

import (
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
)

func (hdl *DBHandler) addQuestionFiles(id uint, desc *[]byte, code *[]byte, data *[]byte) error {
	strid := strconv.Itoa(int(id))
	batch := new(leveldb.Batch)
	batch.Put([]byte("desc_"+strid), *desc)
	batch.Put([]byte("code_"+strid), *code)
	batch.Put([]byte("data_"+strid), *data)
	err := hdl.kvDB.Write(batch, nil)
	return err
}

//AddKV Add question/code values
func (hdl *DBHandler) AddKV(key []byte, value []byte) error {
	err := hdl.kvDB.Put(key, value, nil)
	return err
}
