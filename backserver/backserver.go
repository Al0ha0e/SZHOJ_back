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
func (this *BackServer) Init() error {
	this.handler = dbhandler.GetDBHandler()
	err := this.handler.InitDBHandler()
	this.server = gin.Default()
	this.server.GET("/pgquest", this.getQuestionsByPage)
	this.server.POST("/questions", this.getOrSetQuestions)
	return err
}

//Run Start BackServer
func (this *BackServer) Run() {
	this.server.Run(":8060")
}
