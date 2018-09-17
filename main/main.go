package main

import (
	"github.com/gin-gonic/gin"

	"www.echoecho.cn/main/app/router"
)

func main(){
	r := gin.Default()

	// admin routers
	router.AdminRouter(r)

	// api routers for mini program
	router.ApiRouter(r)

	r.Run()
}