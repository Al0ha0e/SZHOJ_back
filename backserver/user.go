package backserver

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Al0ha0e/SZHOJ_back/dbhandler"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (bs *BackServer) login(c *gin.Context) {
	fmt.Println("LOGIN")
	user := &dbhandler.User{}
	if err := c.ShouldBind(user); err == nil {
		if len(user.Name) > 15 || len(user.Password) > 15 || len(user.Password) < 8 {
			c.String(http.StatusOK, `invalid form`)
			return
		}
		tUser := bs.handler.GetUserByName(user.Name)
		if tUser.ID == 0 {
			c.String(http.StatusAccepted, `wrong password or unregistered`)
			return
		} else if tUser.Password != user.Password {
			c.String(http.StatusAccepted, `wrong password or unregistered`)
			return
		} else {
			session := sessions.Default(c)
			session.Set("loggedIn", "true")
			session.Set("username", user.Name)
			session.Set("userId", tUser.ID)
			session.Save()
			fmt.Println("EN?", session.Get("loggedIn"), session.Get("username"), session.Get("userId"))
			c.String(http.StatusAccepted, "login success "+strconv.Itoa(int(tUser.ID)))
		}

	} else {
		fmt.Println("SSBSBS", err)
		c.String(http.StatusBadRequest, `invalid form`)
	}
}

func (bs *BackServer) register(c *gin.Context) {
	user := &dbhandler.User{}
	if err := c.ShouldBind(user); err == nil {
		fmt.Println("OOO", user.Name, user.Password)
		if len(user.Name) > 15 || len(user.Password) > 15 || len(user.Password) < 8 {
			c.String(http.StatusOK, `invalid form`)
			return
		}
		tUser := bs.handler.GetUserByName(user.Name)
		if tUser.ID != 0 {
			c.String(http.StatusAccepted, `username has been registered`)
			return
		}
		bs.handler.AddUser(user.Name, user.Password)
		c.String(http.StatusAccepted, "register success")

	} else {
		c.String(http.StatusBadRequest, `invalid form`)
	}
}
