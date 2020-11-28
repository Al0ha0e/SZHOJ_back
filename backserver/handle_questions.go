package backserver

import (
	"fmt"
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

func (bs *BackServer) getOrSetQuestions(c *gin.Context) {
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
