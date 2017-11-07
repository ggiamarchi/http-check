package api

import (
	"os"

	"gopkg.in/gin-gonic/gin.v1"
)

func checkAuthorization(c *gin.Context) {
	secretKeyValues := c.Request.Header["X-Secret-Key"]
	if len(secretKeyValues) != 1 || secretKeyValues[0] != os.Getenv("HTTP_CHECK_SECRET_KEY") {
		c.JSON(403, gin.H{"message": "unauthorized"})
		c.Abort()
	}
}

func CreateEngine() *gin.Engine {
	engine := gin.Default()

	engine.Use(checkAuthorization)

	engine.GET("/df", df)
	engine.HEAD("/df", df)

	engine.GET("/df/check", dfCheck)
	engine.HEAD("/df/check", dfCheck)

	return engine
}

func RunServer(engine *gin.Engine) {
	engine.Run()
}
