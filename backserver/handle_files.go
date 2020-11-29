package backserver

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (bs *BackServer) commitAnswer(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "file error")
		return
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		c.String(http.StatusBadRequest, "file error")
		return
	}
	//CAUTION!!! AUTH
	fmt.Println(header.Filename, string(content))
}
