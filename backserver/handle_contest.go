/************
SZHOJ　V１.0.0 后端
由孙梓涵编写
本页面用于处理比赛相关请求
************/
package backserver

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Al0ha0e/SZHOJ_back/dbhandler"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//添加比赛的API
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

//按页面获取比赛
func (bs *BackServer) getContestByPage(c *gin.Context) {
	pg := c.DefaultQuery("pg", "1")
	cnt := c.DefaultQuery("cnt", "10")
	pgnum, _ := strconv.ParseUint(pg, 0, 32)
	ippg, _ := strconv.ParseUint(cnt, 0, 32)
	ret := bs.handler.GetContestByPage(pgnum, ippg)
	ccnt := bs.handler.GetContestCnt()
	if len(*ret) < 1 {
		c.String(http.StatusNotFound, "no such status")
	} else {
		c.JSON(http.StatusOK, gin.H{"contest": *ret, "count": ccnt})
	}
}

//按ID获取比赛
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
	//检查用户是否在比赛用户组中
	for _, usr := range usergroup.Users {
		if uid == usr.ID {
			succ = true
			break
		}
	}
	if succ || uid == contest.Creator { //需要用户在比赛用户组中或为比赛创建者
		c.JSON(http.StatusOK, contest)
	} else {
		c.String(http.StatusUnauthorized, "no authorize")
	}
}

//获取用户自身比赛状态
func (bs *BackServer) getContestStatus(c *gin.Context) {
	//Auth
	cids := c.DefaultQuery("cid", "0")
	cid, _ := strconv.ParseUint(cids, 0, 32)
	uids := c.DefaultQuery("uid", "0")
	uid, _ := strconv.ParseUint(uids, 0, 32)
	status := bs.handler.GetContestStatus(cid, uid)
	if len(*status) == 0 {
		c.String(http.StatusNotFound, "no such status")
	} else {
		c.JSON(http.StatusOK, *status)
	}
}

//获取榜单
func (bs *BackServer) getTotalStatus(c *gin.Context) {
	cids := c.DefaultQuery("cid", "0")
	cid, _ := strconv.ParseUint(cids, 0, 32)
	status := bs.handler.GetTotalContestStatus(cid)
	if len(*status) == 0 {
		c.String(http.StatusNotFound, "no such status")
	} else {
		c.JSON(http.StatusOK, *status)
	}
}
