/************
SZHOJ　V１.0.0 后端
由孙梓涵编写
本页面用于处理状态请求
************/
package backserver

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Al0ha0e/SZHOJ_back/dbhandler"
	"github.com/gin-gonic/gin"
)

//按页面获取评测状态
func (bs *BackServer) getStatusByPage(c *gin.Context) {
	pg := c.DefaultQuery("pg", "1")
	cnt := c.DefaultQuery("cnt", "10")
	pgnum, _ := strconv.ParseUint(pg, 0, 32)
	ippg, _ := strconv.ParseUint(cnt, 0, 32)
	ret := bs.handler.GetStatusByPage(pgnum, ippg)
	scnt := bs.handler.GetStatusCnt()
	if len(*ret) < 1 {
		c.String(http.StatusNotFound, "no such status")
	} else {
		c.JSON(http.StatusOK, gin.H{"status": *ret, "count": scnt})
	}
}

//按参数获取状态
func (bs *BackServer) getOrSetStatus(c *gin.Context) {
	sInfo := dbhandler.Status{}
	if err := c.ShouldBind(&sInfo); err == nil {
		if sInfo.ID != 0 {
			sInfo.ID = 0
			ret := bs.handler.GetStatus(&sInfo)
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

//获取最小化状态，仅包含用户id，题目id，结果
func (bs *BackServer) getMiniStatus(c *gin.Context) {
	uid := c.DefaultQuery("uid", "0")
	uidNum, _ := strconv.ParseUint(uid, 0, 32)
	ret := bs.handler.GetMiniStatus(uidNum)
	if len(*ret) < 1 {
		c.String(http.StatusNotFound, "no such status")
	} else {
		c.JSON(http.StatusOK, *ret)
	}
}

//获取状态细节（代码）
func (bs *BackServer) getStatusDetail(c *gin.Context) {
	sids := c.DefaultQuery("sid", "0")
	sid, _ := strconv.ParseUint(sids, 0, 32)
	status := bs.handler.GetStatusByID(sid)
	code, _ := bs.handler.GetCode(int(sid))
	if status.ID == 0 {
		c.String(http.StatusNotFound, "no such status")
	} else {
		c.JSON(http.StatusOK, gin.H{"status": status, "code": string(code)})
	}
}
