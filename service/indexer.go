package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"holders/conf"
	"holders/jsonrpc"
	"math/rand"
	"net/http"
)

type land struct {
	Index int `json:"index"`
	SeedIndex int `json:"seedIndex"`
	Unix int64 `json:"unix"`
}

type kList struct {
	Owner string `json:"owner"`
	KIDS []string `json:"kids"`
}

func getLands(c *gin.Context)  {
	owner := c.Param("owner")
	if owner == "" {
		handleError(c, errors.New("invalid params"))
		return
	}

	param := jsonrpc.CallParam{
		KID:    conf.FarmScript,
		Method: "$getLands",
		Params: []interface{}{owner},
	}

	result, err := call(param)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK,result)
}

func getBalances(c *gin.Context)  {
	var kids kList
	if err := c.ShouldBind(&kids); err != nil {
		handleError(c, err)
		return
	}
	if kids.KIDS == nil || kids.Owner == ""{
		handleError(c, errors.New("invalid params"))
		return
	}

	bMap := make(map[string]interface{})
	//for _, kid := range kids.KIDS {
	//	param := jsonrpc.CallParam{
	//		KID:    "",
	//		Method: "",
	//		Params: []interface{}{kids.Owner},
	//	}
	//	result, err := call(param)
	//	if err != nil {
	//		bMap[kid] = 0
	//		continue
	//	}
	//	bMap[kid] = result
	//}

	for _, kid := range kids.KIDS {
		// 生成100到10000之间的随机数
		randomNum := rand.Intn(10000-100+1) + 100
		bMap[kid] = randomNum
	}

	c.JSON(http.StatusOK,bMap)

}

func ordCall(c *gin.Context)  {
	var param jsonrpc.CallParam
	if err := c.ShouldBind(&param); err != nil {
		handleError(c, err)
		return
	}

	result, err := call(param)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK,result)
}

func call(param jsonrpc.CallParam) (any, error) {
	cli, err := jsonrpc.NewClient(conf.NodeUrl)
	if err != nil {
		return nil, err
	}
	return cli.CallContract(param)
}

