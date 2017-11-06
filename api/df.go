package api

import (
	"strconv"

	"github.com/ggiamarchi/http-check/commands"
	"gopkg.in/gin-gonic/gin.v1"
)

func df(c *gin.Context) {
	entries, err := commands.Df()
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, entries)
}

func dfCheck(c *gin.Context) {

	mountpoint := c.Query("mountpoint")
	field := c.Query("field")
	op := c.Query("op")
	value, e := strconv.Atoi(c.Query("value"))

	if e != nil {
		c.JSON(400, gin.H{"message": "value parameter must be an integer"})
		return
	}

	ok, err := commands.DfCheck(mountpoint, field, op, value)

	if err != nil {
		var code int
		if err.Code == "INTERNAL" {
			code = 500
		} else {
			code = 400
		}
		c.JSON(code, gin.H{"message": err.Error()})
		return
	}

	if ok {
		c.JSON(200, gin.H{"message": "ok"})
	} else {
		c.JSON(409, gin.H{"message": "contraint not satisfied"})
	}

}
