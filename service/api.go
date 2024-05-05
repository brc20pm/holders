package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"holders/models"
	"net/http"
)

type GinService struct {
	Service *gin.Engine
}

// 创建一个Gin服务
func NewGinService() *GinService {
	//gin.SetMode(gin.ReleaseMode)
	service := gin.Default()
	return &GinService{service}
}

// 注入API路由
func (g *GinService) loadGroupAPI() *gin.RouterGroup {
	group := g.Service.Group("/api")

	//需要鉴权的接口
	//group.Use(authMiddleware())

	//持有数据
	{
		//代币信息
		group.GET("/token/:kid", getToken)
		//获取钱包持有数据
		group.GET("/wallet/:owner", getWalletHolds)
		//获取持有的TokenId 列表
		group.GET("/tokenIds", getTokenIds)
		//获取代币持有分布
		group.GET("/dist/20/:kid", getDist20)
		//获取NFT持有分布
		group.GET("/dist/721/:kid", getDist721)
	}


	return group
}

//用户验证
//func authMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		if false {
//			c.AbortWithStatus(http.StatusUnauthorized)
//		} else {
//			c.Next()
//		}
//	}
//}

func handleError(c *gin.Context, err any) {
	var result models.Result
	result.Code = http.StatusInternalServerError
	result.Msg = fmt.Sprint(err)
	c.JSON(result.Code, result)
	return
}

/*
* 功能介绍: 启动接口服务
* @receiver g
* @param port
 */
func (g *GinService) Run(port string) {
	g.loadGroupAPI()
	//启动接口服务
	g.Service.Run(port)
}
