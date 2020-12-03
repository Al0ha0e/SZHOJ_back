package backserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Al0ha0e/SZHOJ_back/dbhandler"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (bs *BackServer) commitToScheduler(status *dbhandler.Status) {
	bs.jobScheduler.CommitChan <- status
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
	session := sessions.Default(c)
	uid := session.Get("userId").(uint)
	if session.Get("loggedIn") != "true" || uid != commitInfo.ID {
		c.String(http.StatusForbidden, "no authority")
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
	question := bs.handler.GetQuestionByID(uint64(commitInfo.QuestionID))
	cid := question.ContestID
	if cid != 0 {
		contest := bs.handler.GetContestByID(uint64(cid))
		usergroup := bs.handler.GetUserGroupByID(contest.UserGroupID)
		succ := false
		for _, usr := range usergroup.Users {
			if uid == usr.ID {
				succ = true
				break
			}
		}
		if !succ {
			c.String(http.StatusUnauthorized, "no authorize")
			return
		}
	}
	status := &dbhandler.Status{
		QuestionID: commitInfo.QuestionID,
		UserID:     commitInfo.ID,
		CommitTime: time.Now(),
	}
	status.PrepareForCreation(bs.handler, &code)
	bs.handler.AddStatus(status)

	go bs.commitToScheduler(status)
	c.String(http.StatusOK, "commit success")
}
