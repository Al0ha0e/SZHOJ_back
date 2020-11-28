package backserver

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Al0ha0e/SZHOJ_back/dbhandler"
	"github.com/gin-gonic/gin"
)

func (bs *BackServer) getQuestionsByPage(c *gin.Context) {
	pg := c.DefaultQuery("pg", "1")
	cnt := c.DefaultQuery("cnt", "10")
	pgnum, _ := strconv.ParseUint(pg, 0, 32)
	ippg, _ := strconv.ParseUint(cnt, 0, 32)
	// if err != nil{
	// 	c.
	// }
	ret := bs.handler.GetQuestionsByPage(pgnum, ippg)
	c.JSON(200, *ret)
}

func (bs *BackServer) getOrSetQuestions(c *gin.Context) {
	qInfo := dbhandler.Question{}
	if err := c.ShouldBind(&qInfo); err == nil {
		fmt.Println(qInfo.Name, qInfo.ID, qInfo.Tags)
	} else {
		fmt.Println("NONONO", err)
	}
	c.String(http.StatusOK, `HAHAHA`)
}
