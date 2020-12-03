package backserver

import (
	"fmt"
	"net/http"

	"github.com/Al0ha0e/SZHOJ_back/dbhandler"
	"github.com/Al0ha0e/SZHOJ_back/scheduler"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

//BackServer backend server
type BackServer struct {
	server       *gin.Engine
	handler      *dbhandler.DBHandler
	jobScheduler *scheduler.Scheduler
	store        cookie.Store
}

//GetBackServer Get a backserver instance
func GetBackServer() *BackServer {
	return &BackServer{}
}

func auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		fmt.Println("SESSION", session.Get("loggedIn"), session.Get("username"), session.Get("userId"))
		if session.Get("loggedIn") != "true" {
			if c.FullPath() != "/login" && c.FullPath() != "/register" {
				c.String(http.StatusForbidden, "please login first")
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
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

	bs.server.Use(cors())

	bs.store = cookie.NewStore([]byte("secret"))
	bs.server.Use(sessions.Sessions("mysession", bs.store))

	bs.server.Use(auth())

	bs.server.GET("/question", bs.getQuesionByID)
	bs.server.GET("/contest", bs.getContestByID)
	bs.server.GET("/pgquest", bs.getQuestionsByPage)
	bs.server.GET("/quedesc", bs.getQuestionDesc)
	bs.server.GET("/pgstatus", bs.getStatusByPage)
	bs.server.GET("/pgcontest", bs.getContestByPage)
	bs.server.GET("/ministatus", bs.getMiniStatus)
	bs.server.GET("/singlestatus", bs.getStatusDetail)
	bs.server.GET("/usergroup", bs.getUserGroup)
	bs.server.POST("/login", bs.login)
	bs.server.POST("/register", bs.register)
	bs.server.POST("/questions", bs.getQuestions)
	bs.server.POST("/status", bs.getOrSetStatus)
	bs.server.POST("/upquestion", bs.uploadQuestion)
	bs.server.POST("/upanswer", bs.commitAnswer)
	bs.server.POST("/addgroup", bs.addUserGroup)
	bs.server.POST("/delgroup", bs.deleteUserGroup)
	bs.server.POST("/addcontest", bs.addContest)
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
