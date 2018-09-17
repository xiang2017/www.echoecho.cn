package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"www.echoecho.cn/main/app/model"
	"fmt"
)

func Test(c *gin.Context) {
	//var res int
	//var err error
	//if res, err = model.RedisClient.Get("casd").Int(); err != nil {
	//	fmt.Println("error:", err.Error())
	//} else{
	//	fmt.Println(res)
	//}

	if res, err := model.RedisClient.SRandMember("a").Int(); err != nil {
		fmt.Println(err.Error())
	} else{
		fmt.Println(res)
	}

	c.String(http.StatusOK, "s")
}