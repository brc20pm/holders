package api

import (
	"github.com/gin-gonic/gin"
	"holders/conf"
	"holders/jsonrpc"
	"net/http"
)


type kList struct {
	Owner string `json:"owner"`
	KIDS []string `json:"kids"`
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

