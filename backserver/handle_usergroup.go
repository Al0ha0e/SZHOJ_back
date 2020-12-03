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

func (bs *BackServer) addUserGroup(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("userId")
	err := c.Request.ParseMultipartForm(16 << 10) //16kb
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
	c.String(http.StatusOK, "add success")
}
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
