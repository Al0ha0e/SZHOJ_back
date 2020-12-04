/************
SZHOJ　V１.0.0 后端
由孙梓涵编写
本页面用于处理键值对的存取
************/

package dbhandler

import (
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
)

//添加问题相关的文件
func (hdl *DBHandler) addQuestionFiles(id uint, desc *[]byte, datain *[]byte, dataout *[]byte) error {
	strid := strconv.Itoa(int(id))
	//作为一个事务提交
	batch := new(leveldb.Batch)
	batch.Put([]byte("desc_"+strid), *desc)
	batch.Put([]byte("datain_"+strid), *datain)
	batch.Put([]byte("dataout_"+strid), *dataout)
	err := hdl.kvDB.Write(batch, nil)
	return err
}

func (hdl *DBHandler) addCommitCode(id uint, code *[]byte) error {
	strid := strconv.Itoa(int(id))
	err := hdl.kvDB.Put([]byte("code_"+strid), *code, nil)
	return err
}

//GetCode get code for status
func (hdl *DBHandler) GetCode(id int) ([]byte, error) {
	strid := strconv.Itoa(id)
	return hdl.kvDB.Get([]byte("code_"+strid), nil)
}

//GetJudgeData Get data for judger
func (hdl *DBHandler) GetJudgeData(id uint) (*[]byte, *[]byte, error) {
	strid := strconv.Itoa(int(id))
	datain, err := hdl.kvDB.Get([]byte("datain_"+strid), nil)
	if err != nil {
		return nil, nil, err
	}
	dataout, err := hdl.kvDB.Get([]byte("dataout_"+strid), nil)
	if err != nil {
		return nil, nil, err
	}
	return &datain, &dataout, nil
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
