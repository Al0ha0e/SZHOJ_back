package backserver

import (
	"github.com/Al0ha0e/SZHOJ_back/dbhandler"
	"github.com/gin-contrib/cors"
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

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true                                                                                                 //允许所有域名
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}                                                                      //允许请求的方法
	config.AllowHeaders = []string{"tus-resumable", "upload-length", "upload-metadata", "cache-control", "x-requested-with", "*"} //允许的Header
	bs.server.Use(cors.New(config))

	bs.server.GET("/pgquest", bs.getQuestionsByPage)
	bs.server.GET("/pgstatus", bs.getStatusByPage)
	bs.server.GET("/ministatus", bs.getMiniStatus)
	bs.server.POST("/questions", bs.getOrSetQuestions)
	bs.server.POST("/status", bs.getOrSetStatus)
	return err
}

//Run Start BackServer
func (bs *BackServer) Run() {
	bs.server.Run(":8060")
}
