package app

import (
	"net/http"
	"strconv"
	"vote_cj/config"

	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()
	r.StaticFS("/static", http.Dir("./static"))
	r.StaticFile("/favicon.ico", "./static/favicon.ico")

	r.Use(Cors())

	//登录
	r.POST("/login", Login)

	//获取参选人信息（完整）
	r.GET("/data.js", GetCompetitors)
	//获取参选人信息（生产）
	r.GET("/data", GetComp)

	//获取参选人信息（完整）
	r.GET("/cands_complete.js", GetCandsComplete)
	//获取参选人信息（生产）
	r.GET("/cands", GetCands)
	//投票
	r.POST("/cand/:group_id/:cand_id", MiddleWare(), Vote)
	r.Run(":"+strconv.Itoa(config.System.Port))
}
