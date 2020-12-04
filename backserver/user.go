/************
SZHOJ　V１.0.0 后端
由孙梓涵编写
本页面用于处理用户相关的请求
************/
package backserver

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Al0ha0e/SZHOJ_back/dbhandler"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//登录接口
func (bs *BackServer) login(c *gin.Context) {
	fmt.Println("LOGIN")
	user := &dbhandler.User{}
	if err := c.ShouldBind(user); err == nil {
		//检查用户名及密码是否符合格式
		if len(user.Name) > 15 || len(user.Password) > 15 || len(user.Password) < 8 {
			c.String(http.StatusOK, `invalid form`)
			return
		}

		tUser := bs.handler.GetUserByName(user.Name)
		if tUser.ID == 0 {
			//用户不存在
			c.String(http.StatusAccepted, `wrong password or unregistered`)
			return
		} else if tUser.Password != user.Password {
			//密码错误
			c.String(http.StatusAccepted, `wrong password or unregistered`)
			return
		} else {
			//登陆成功，后端保存session信息
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

//用户注册
func (bs *BackServer) register(c *gin.Context) {
	user := &dbhandler.User{}
	if err := c.ShouldBind(user); err == nil {
		//检查用户名及密码是否符合格式
		if len(user.Name) > 15 || len(user.Password) > 15 || len(user.Password) < 8 {
			c.String(http.StatusOK, `invalid form`)
			return
		}
		tUser := bs.handler.GetUserByName(user.Name)
		//判断是否注册过
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
