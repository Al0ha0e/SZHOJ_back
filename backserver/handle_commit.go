package backserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Al0ha0e/SZHOJ_back/dbhandler"
	"github.com/Al0ha0e/SZHOJ_back/scheduler"
	"github.com/gin-gonic/gin"
)

func (bs *BackServer) commitToScheduler(code *[]byte, status *dbhandler.Status) {
	bs.jobScheduler.CommitChan <- &scheduler.TaskToCommit{Code: code, Status: status}
}

func (bs *BackServer) commitAnswer(c *gin.Context) {
	err := c.Request.ParseMultipartForm(16 << 10) //16kb
	if err != nil {
		c.String(http.StatusBadRequest, "form error: too large")
		return
	}
	formdata := c.Request.MultipartForm

	commitJSON := formdata.Value["commit"][0]
	var commitInfo dbhandler.MiniStatus
	err = json.Unmarshal([]byte(commitJSON), &commitInfo)
	if err != nil {
		c.String(http.StatusBadRequest, "form error: bad format")
		return
	}
	//hadle files
	files := formdata.File["file"]
	if len(files) != 1 {
		c.String(http.StatusBadRequest, "file count mismatch")
		return
	}
	file, err := files[0].Open()
	defer file.Close()
	if err != nil {
		c.String(http.StatusBadRequest, "file error: cannot open file")
		return
	}
	code, err := ioutil.ReadAll(file)
	if err != nil {
		c.String(http.StatusBadRequest, "file error: cannot read file")
		return
	}
	status := &dbhandler.Status{
		QuestionID: commitInfo.QuestionID,
		UserID:     commitInfo.ID,
		CommitTime: time.Now(),
	}
	bs.handler.AddStatus(status)
	go bs.commitToScheduler(&code, status)
	c.String(http.StatusOK, "commit success")
}
