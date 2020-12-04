/************
SZHOJ　V１.0.0 后端
由孙梓涵编写
本页面用于处理题目提交
************/
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

//向调度器添加任务
func (bs *BackServer) commitToScheduler(status *dbhandler.Status) {
	bs.jobScheduler.CommitChan <- status
}

//提交答案API
func (bs *BackServer) commitAnswer(c *gin.Context) {
	err := c.Request.ParseMultipartForm(16 << 10) //请求大小限制
	if err != nil {
		c.String(http.StatusBadRequest, "form error: too large")
		return
	}
	formdata := c.Request.MultipartForm

	commitJSON := formdata.Value["commit"][0]
	var commitInfo dbhandler.MiniStatus
	err = json.Unmarshal([]byte(commitJSON), &commitInfo) //解析请求所带结构体JSON
	if err != nil {
		c.String(http.StatusBadRequest, "form error: bad format")
		return
	}
	session := sessions.Default(c)
	uid := session.Get("userId").(uint)
	if session.Get("loggedIn") != "true" || uid != commitInfo.ID { //鉴权，提交者必须与登陆者相符
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
	//判断题目是否为比赛题目
	question := bs.handler.GetQuestionByID(uint64(commitInfo.QuestionID))
	cid := question.ContestID
	if cid != 0 {
		contest := bs.handler.GetContestByID(uint64(cid))
		usergroup := bs.handler.GetUserGroupByID(contest.UserGroupID)
		succ := false
		//判断提交者是否在比赛对应用户组中
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

	//异步提交任务
	go bs.commitToScheduler(status)
	c.String(http.StatusOK, "commit success")
}
