package api

import (
	"fmt"
	"github.com/gin-contrib/cors"
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

	// 创建一个CORS配置
	config := cors.DefaultConfig()

	// 如果你只想允许特定的源进行访问，可以这样设置
	//config.AllowOrigins = []string{"http://example.com"}

	// 允许所有来源进行访问（请注意，这通常只在开发环境中使用，生产环境应更严格）
	config.AllowAllOrigins = true

	// 如果需要的话，你还可以配置允许的HTTP方法、头部等其他选项

	// 应用CORS中间件到所有路由
	service.Use(cors.New(config))

	return &GinService{service}
}

// 注入API路由
func (g *GinService) loadGroupAPI() *gin.RouterGroup {
	group := g.Service.Group("/assets")

	//需要鉴权的接口
	//group.Use(authMiddleware())

	//持有数据
	{

		//group.GET("/lands/:owner",getLands)
		group.POST("/batch_balance",getBalances)

		group.POST("/ord_call",ordCall)

		//代币信息
		group.GET("/token/:kid", getToken)
		//批量获取代币信息
		group.POST("/token/batch",getTokenForBatch)
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
