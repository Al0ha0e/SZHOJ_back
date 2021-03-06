/************
SZHOJ　V１.0.0 后端
由孙梓涵编写
本页面用于处理问题上传及查询请求
************/
package backserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/Al0ha0e/SZHOJ_back/dbhandler"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//按ID获取问题
func (bs *BackServer) getQuesionByID(c *gin.Context) {
	qids := c.DefaultQuery("qid", "0")
	qid, err := strconv.ParseUint(qids, 0, 32)
	if err != nil {
		c.String(http.StatusBadRequest, "bad format")
	}
	ret := bs.handler.GetQuestionByID(qid)
	if ret == nil {
		c.String(http.StatusNotFound, "no such page")
	} else {
		c.JSON(http.StatusOK, *ret)
	}
}

//按页面获取问题
func (bs *BackServer) getQuestionsByPage(c *gin.Context) {
	fmt.Println("HEAR")
	pg := c.DefaultQuery("pg", "1")
	cnt := c.DefaultQuery("cnt", "10")
	pgnum, err := strconv.ParseUint(pg, 0, 32)
	if err != nil {
		c.String(http.StatusBadRequest, "bad format")
	}
	ippg, err := strconv.ParseUint(cnt, 0, 32)
	if err != nil {
		c.String(http.StatusBadRequest, "bad format")
	}
	ret := bs.handler.GetQuestionsByPage(pgnum, ippg)
	qcnt := bs.handler.GetQuestionCnt()
	if len(*ret) < 1 {
		c.String(http.StatusNotFound, "no such page or no questions")
	} else {
		c.JSON(http.StatusOK, gin.H{"questions": *ret, "count": qcnt})
	}
}

// func (bs *BackServer) getQuestions(c *gin.Context) {
// 	qInfo := dbhandler.Question{}
// 	if err := c.ShouldBind(&qInfo); err == nil {
// 		fmt.Println(qInfo.Name, qInfo.ID, qInfo.Tags)
// 		qInfo.ID = 0
// 		ret := bs.handler.GetQuestions(&qInfo)
// 		if len(*ret) < 1 {
// 			c.String(http.StatusNotFound, "question not found")
// 		} else {
// 			c.JSON(http.StatusOK, *ret)
// 		}

// 	} else {
// 		c.String(http.StatusNotFound, `invalid form`)
// 	}
// }

//获取问题描述
func (bs *BackServer) getQuestionDesc(c *gin.Context) {
	qids := c.DefaultQuery("qid", "1")
	qid, err := strconv.ParseUint(qids, 0, 32)
	if err != nil {
		c.String(http.StatusBadRequest, "bad format")
		return
	}
	ret, err := bs.handler.GetQuestionDesc(qid) //在键值数据库中获取题目描述
	if err != nil {
		fmt.Printf(err.Error())
		if err.Error() == "leveldb: not found" {
			c.String(http.StatusNotFound, "fail to get desc")
		} else {
			c.String(http.StatusInternalServerError, "fail to get desc")
		}
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", ret)
}

//上传问题
func (bs *BackServer) uploadQuestion(c *gin.Context) {
	fmt.Println("UPLOAD")
	err := c.Request.ParseMultipartForm(16 << 10) //限制上传数据量
	if err != nil {
		c.String(http.StatusBadRequest, "form error: too large")
		return
	}
	formdata := c.Request.MultipartForm

	qinfoJSON := formdata.Value["qinfo"][0]
	var qinfo dbhandler.Question
	err = json.Unmarshal([]byte(qinfoJSON), &qinfo) //解析题目信息
	if err != nil {
		c.String(http.StatusBadRequest, "form error: bad format")
		return
	}

	fmt.Println(qinfo.Name, qinfo.Creator)
	nameL := utf8.RuneCountInString(qinfo.Name) //限制题目名称长度
	if nameL > 15 || nameL < 1 {
		fmt.Println("BAD", qinfo.Name, len(qinfo.Name))
		c.String(http.StatusOK, "bad name format")
		return
	}

	qinfo.ID = 0
	session := sessions.Default(c)
	if session.Get("loggedIn") != "true" || session.Get("userId") != qinfo.Creator {
		//鉴权，要求登录用户与提交者一致
		c.String(http.StatusForbidden, "no authority")
		return
	}

	//hadle files
	files := formdata.File["file"]
	if len(files) != 3 {
		c.String(http.StatusBadRequest, "file count mismatch")
		return
	}
	//读取数据
	qfiles := make([]*[]byte, 3)
	for i, fh := range files {
		file, err := fh.Open()
		if err != nil {
			c.String(http.StatusBadRequest, "file error: cannot open file")
			return
		}
		defer file.Close()
		content, err := ioutil.ReadAll(file)
		if err != nil {
			c.String(http.StatusBadRequest, "file error: cannot read file")
			return
		}
		qfiles[i] = &content
	}
	qinfo.PrepareForCreation(bs.handler, qfiles[0], qfiles[1], qfiles[2])
	//it is a transaction
	fmt.Println("TRANS ", bs.handler.AddQuestion(&qinfo))
	if !qinfo.Success {
		c.String(http.StatusInternalServerError, "upload failed")
		return
	}
	c.String(http.StatusOK, "upload success")
}
