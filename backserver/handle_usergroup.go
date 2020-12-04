/************
SZHOJ　V１.0.0 后端
由孙梓涵编写
本页面用于处理用户组
************/
package backserver

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/Al0ha0e/SZHOJ_back/dbhandler"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//按用户id获取用户组，参数为1表示获取参加的，否则为获取创建的
func (bs *BackServer) getUserGroup(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("userId")
	attended := c.DefaultQuery("attend", "1")
	fmt.Println("GET UG", attended)
	var ug *[]dbhandler.UserGroup
	if attended == "1" {
		ug = bs.handler.GetAttendedUserGroups(uid.(uint))
	} else {
		ug = bs.handler.GetCreatedUserGroups(uid.(uint))
	}

	c.JSON(http.StatusOK, *ug)
}

//从文件中读取用户组
func (bs *BackServer) parseUserGroup(content string) (*[]dbhandler.User, error) {
	ret := make([]dbhandler.User, 0)
	users := strings.Split(content, "|")
	for _, usr := range users {
		ufield := strings.Fields(usr)
		if len(ufield) != 1 {
			return nil, fmt.Errorf("field mismatch")
		}
		uid, err := strconv.ParseUint(ufield[0], 0, 32)
		if err != nil {
			return nil, err
		}
		ret = append(ret, *bs.handler.GetUserByID(uint64(uid)))
	}
	return &ret, nil
}

//添加用户组
func (bs *BackServer) addUserGroup(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("userId")
	err := c.Request.ParseMultipartForm(16 << 10) //数据大小限制
	if err != nil {
		c.String(http.StatusBadRequest, "form error: too large")
		return
	}

	formdata := c.Request.MultipartForm
	usergroup := &dbhandler.UserGroup{Creator: uid.(uint), Name: formdata.Value["name"][0]}
	files := formdata.File["file"]
	file, err := files[0].Open()
	defer file.Close()
	if err != nil {
		c.String(http.StatusBadRequest, "file error: cannot open file")
		return
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		c.String(http.StatusBadRequest, "file error: cannot read file")
		return
	}

	users, err := bs.parseUserGroup(string(content))
	if err != nil {
		c.String(http.StatusBadRequest, "file error: bad format")
		return
	}
	usergroup.Users = *users
	bs.handler.AddUserGroup(usergroup)
	if usergroup.ID == 0 {
		c.String(http.StatusInternalServerError, "add failed")
	} else {
		c.String(http.StatusOK, "add success")
	}
}

//删除用户组
func (bs *BackServer) deleteUserGroup(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("userId")
	gids := c.DefaultQuery("gid", "0")
	gid, err := strconv.ParseUint(gids, 0, 32)
	if err != nil {
		c.String(http.StatusBadRequest, "form error: bad format")
		return
	}
	group := bs.handler.GetUserGroupByID(uint(gid))
	if group.ID == 0 || group.Creator != uid {
		c.String(http.StatusBadRequest, "cannot delete")
		return
	}
	bs.handler.DeleteUserGroup(gid)
	c.String(http.StatusOK, "delete success")
}
