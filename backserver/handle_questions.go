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

func (bs *BackServer) getQuestionsByPage(c *gin.Context) {
	fmt.Println("HEAR")
	pg := c.DefaultQuery("pg", "1")
	cnt := c.DefaultQuery("cnt", "10")
	pgnum, _ := strconv.ParseUint(pg, 0, 32)
	ippg, _ := strconv.ParseUint(cnt, 0, 32)
	// if err != nil{
	// 	c.
	// }
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

	//hadle files
	files := formdata.File["file"]
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
	bs.handler.AddQuestion(&qinfo)
	if !qinfo.Success {
		c.String(http.StatusInternalServerError, "upload failed")
	}
	c.String(http.StatusOK, "upload success")
}
