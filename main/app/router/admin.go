package router

import (
	"github.com/gin-gonic/gin"
	"www.echoecho.cn/main/app/middleware"
	"www.echoecho.cn/main/app/controller/bgcontroller"
	"www.echoecho.cn/main/app/controller"
)

// 设置
func AdminRouter(r *gin.Engine){

	r.Use(middleware.Cors())
	r.GET("/ping", controller.Test)

	// 标签管理
	r.GET("/labels", bgcontroller.GetLabelList)
	r.GET("/label/:id", bgcontroller.GetLabel)
	r.POST("/label", bgcontroller.EditLabel)
	r.DELETE("/label/:id", bgcontroller.DeleteLabel)

	// 问题管理
	r.GET("/questions", bgcontroller.QuestionList)
	r.POST("/question", bgcontroller.EditQuestion)
	r.GET("/question", bgcontroller.GetQuestionInfo)
	r.DELETE("/question/:id", bgcontroller.DeleteQuestion)

	r.POST("/knowledge", bgcontroller.EditKnowledge)
	r.GET("/knowledge", bgcontroller.KnowledgeList)
	r.GET("/knowledge/:id", bgcontroller.GetKnowledge)

}