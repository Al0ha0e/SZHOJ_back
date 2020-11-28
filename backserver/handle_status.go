package backserver

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Al0ha0e/SZHOJ_back/dbhandler"
	"github.com/gin-gonic/gin"
)

func (bs *BackServer) getStatusByPage(c *gin.Context) {
	pg := c.DefaultQuery("pg", "1")
	cnt := c.DefaultQuery("cnt", "10")
	pgnum, _ := strconv.ParseUint(pg, 0, 32)
	ippg, _ := strconv.ParseUint(cnt, 0, 32)
	ret := bs.handler.GetStatusByPage(pgnum, ippg)
	if len(*ret) < 1 {
		c.String(http.StatusNotFound, "no such page")
	} else {
		c.JSON(http.StatusOK, *ret)
	}
}

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
