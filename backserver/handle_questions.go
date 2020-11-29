package backserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Al0ha0e/SZHOJ_back/dbhandler"
	"github.com/gin-gonic/gin"
)

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
	if len(*ret) < 1 {
		c.String(http.StatusNotFound, "no such page")
	} else {
		c.JSON(http.StatusOK, *ret)
	}
}

func (bs *BackServer) getQuestions(c *gin.Context) {
	qInfo := dbhandler.Question{}
	if err := c.ShouldBind(&qInfo); err == nil {
		fmt.Println(qInfo.Name, qInfo.ID, qInfo.Tags)
		if qInfo.ID != 0 {
			fmt.Println("QUERY")
			qInfo.ID = 0
			ret := bs.handler.GetQuestions(&qInfo)
			if len(*ret) < 1 {
				c.String(http.StatusNotFound, "question not found")
			} else {
				c.JSON(http.StatusOK, *ret)
			}
		} else {
			//TODO
			fmt.Println("INSERT")
			c.String(http.StatusOK, `insert success`)
		}
	} else {
		c.String(http.StatusNotFound, `invalid form`)
	}
}

func (bs *BackServer) getQuestionDesc(c *gin.Context) {
	qids := c.DefaultQuery("qid", "1")
	qid, err := strconv.ParseUint(qids, 0, 32)
	if err != nil {
		c.String(http.StatusBadRequest, "bad format")
		return
	}
	ret, err := bs.handler.GetQuestionDesc(qid)
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

func (bs *BackServer) uploadQuestion(c *gin.Context) {
	err := c.Request.ParseMultipartForm(16 << 10) //16kb
	if err != nil {
		c.String(http.StatusBadRequest, "form error: too large")
		return
	}
	formdata := c.Request.MultipartForm

	qinfoJSON := formdata.Value["qinfo"][0]
	var qinfo dbhandler.Question
	err = json.Unmarshal([]byte(qinfoJSON), &qinfo)
	if err != nil {
		c.String(http.StatusBadRequest, "form error: bad format")
		return
	}
	qinfo.ID = 0
	//hadle files
	files := formdata.File["file"]
	if len(files) != 3 {
		c.String(http.StatusBadRequest, "file count mismatch")
		return
	}
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
