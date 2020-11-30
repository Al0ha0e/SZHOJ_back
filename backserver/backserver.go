package backserver

import (
	"github.com/Al0ha0e/SZHOJ_back/dbhandler"
	"github.com/Al0ha0e/SZHOJ_back/scheduler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//BackServer backend server
type BackServer struct {
	server       *gin.Engine
	handler      *dbhandler.DBHandler
	jobScheduler *scheduler.Scheduler
}

//GetBackServer Get a backserver instance
func GetBackServer() *BackServer {
	return &BackServer{}
}

//Init init BackServer
func (bs *BackServer) Init() error {
	bs.jobScheduler = scheduler.GetScheduler()

	bs.handler = dbhandler.GetDBHandler()
	err := bs.handler.InitDBHandler()
	if err != nil {
		return err
	}

	err = bs.jobScheduler.Init(bs.handler)
	if err != nil {
		return err
	}
	bs.jobScheduler.Start()

	bs.server = gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true                                                                                                 //允许所有域名
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}                                                                      //允许请求的方法
	config.AllowHeaders = []string{"tus-resumable", "upload-length", "upload-metadata", "cache-control", "x-requested-with", "*"} //允许的Header
	bs.server.Use(cors.New(config))

	bs.server.GET("/question", bs.getQuesionByID)
	bs.server.GET("/pgquest", bs.getQuestionsByPage)
	bs.server.GET("/quedesc", bs.getQuestionDesc)
	bs.server.GET("/pgstatus", bs.getStatusByPage)
	bs.server.GET("/ministatus", bs.getMiniStatus)
	bs.server.POST("/questions", bs.getQuestions)
	bs.server.POST("/status", bs.getOrSetStatus)
	bs.server.POST("/upquestion", bs.uploadQuestion)
	bs.server.POST("/upanswer", bs.commitAnswer)
	return err
}

//Dispose Release Resources when shut down
func (bs *BackServer) Dispose() {
	bs.handler.Dispose()
}

//Run Start BackServer
func (bs *BackServer) Run() {
	bs.server.Run(":8060")
}
