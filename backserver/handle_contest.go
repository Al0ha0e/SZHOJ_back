package backserver

import (
	"github.com/Al0ha0e/SZHOJ_back/dbhandler"
	"github.com/gin-gonic/gin"
)

func (bs *BackServer) addContest(c *gin.Context) {
	cInfo := dbhandler.Contest{}
	if err := c.ShouldBind(&cInfo); err == nil {
		cInfo.ID = 0
		bs.handler.AddContest(&cInfo)
	}
}

func (bs *BackServer) getContestStatus(c *gin.Context) {
	//Auth

}

func (bs *BackServer) getTotalStatus(c *gin.Context) {

}
