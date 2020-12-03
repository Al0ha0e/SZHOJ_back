package backserver

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Al0ha0e/SZHOJ_back/dbhandler"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (bs *BackServer) addContest(c *gin.Context) {
	cInfo := dbhandler.Contest{}
	if err := c.ShouldBind(&cInfo); err == nil {
		cInfo.ID = 0
		bs.handler.AddContest(&cInfo)
		c.String(http.StatusOK, "contest add success")
	} else {
		fmt.Println(err)
		c.String(http.StatusBadRequest, "invalid form")
	}
}

func (bs *BackServer) getContestByPage(c *gin.Context) {
	pg := c.DefaultQuery("pg", "1")
	cnt := c.DefaultQuery("cnt", "10")
	pgnum, _ := strconv.ParseUint(pg, 0, 32)
	ippg, _ := strconv.ParseUint(cnt, 0, 32)
	ret := bs.handler.GetContestByPage(pgnum, ippg)
	if len(*ret) < 1 {
		c.String(http.StatusNotFound, "no such status")
	} else {
		c.JSON(http.StatusOK, *ret)
	}
}

func (bs *BackServer) getContestByID(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("userId").(uint)
	cids := c.DefaultQuery("cid", "0")
	cid, err := strconv.ParseUint(cids, 0, 32)
	if err != nil {
		c.String(http.StatusBadRequest, "form error: bad format")
		return
	}
	contest := bs.handler.GetContestByID(cid)
	gid := contest.UserGroupID
	usergroup := bs.handler.GetUserGroupByID(gid)
	succ := false
	for _, usr := range usergroup.Users {
		if uid == usr.ID {
			succ = true
			break
		}
	}
	if succ || uid == contest.Creator {
		c.JSON(http.StatusOK, contest)
	} else {
		c.String(http.StatusForbidden, "no authorize")
	}
}

func (bs *BackServer) getContestStatus(c *gin.Context) {
	//Auth

}

func (bs *BackServer) getTotalStatus(c *gin.Context) {

}
