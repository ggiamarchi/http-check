package api

import (
	"gopkg.in/gin-gonic/gin.v1"
)

func CreateEngine() *gin.Engine {
	engine := gin.Default()

	engine.GET("/df", df)
	engine.HEAD("/df", df)

	engine.GET("/df/check", dfCheck)
	engine.HEAD("/df/check", dfCheck)

	return engine
}

func RunServer(engine *gin.Engine) {
	engine.Run()
}
