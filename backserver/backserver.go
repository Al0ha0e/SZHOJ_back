package backserver

import (
	"github.com/Al0ha0e/SZHOJ_back/dbhandler"
	"github.com/gin-gonic/gin"
)

//BackServer backend server
type BackServer struct {
	server  *gin.Engine
	handler *dbhandler.DBHandler
}

//GetBackServer Get a backserver instance
func GetBackServer() *BackServer {
	return &BackServer{}
}

//Init init BackServer
func (bs *BackServer) Init() error {
	bs.handler = dbhandler.GetDBHandler()
	err := bs.handler.InitDBHandler()
	bs.server = gin.Default()
	bs.server.GET("/pgquest", bs.getQuestionsByPage)
	bs.server.GET("/pgstatus", bs.getStatusByPage)
	bs.server.POST("/questions", bs.getOrSetQuestions)
	bs.server.POST("/status", bs.getOrSetStatus)
	return err
}

//Run Start BackServer
func (bs *BackServer) Run() {
	bs.server.Run(":8060")
}
