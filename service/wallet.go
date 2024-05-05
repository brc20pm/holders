package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"holders/db"
	"holders/models"
	"net/http"
)

func getToken(c *gin.Context) {
	var result models.Result
	kid := c.Param("kid")
	if kid == "" {
		handleError(c, errors.New("invalid params"))
		return
	}

	token, err := db.FindToken(kid)
	if err != nil {
		handleError(c, err)
		return
	}

	result.Code = http.StatusOK
	result.Msg = "success"
	result.Data = token
	c.JSON(http.StatusOK, result)
}

func getWalletHolds(c *gin.Context) {
	var result models.Result
	owner := c.Param("owner")
	if owner == "" {
		handleError(c, errors.New("invalid params"))
		return
	}
	holds, err := db.FindWalletHold(owner)
	if err != nil {
		handleError(c, err)
		return
	}
	result.Code = http.StatusOK
	result.Msg = "success"
	result.Data = holds
	c.JSON(http.StatusOK, result)
}

func getTokenIds(c *gin.Context) {
	var result models.Result

	kid, _ := c.GetQuery("kid")
	owner, _ := c.GetQuery("owner")
	if kid == "" || owner == "" {
		handleError(c, errors.New("invalid params"))
		return
	}
	tokenIds, err := db.FindTokenIds(kid, owner)
	if err != nil {
		handleError(c, err)
		return
	}
	result.Code = http.StatusOK
	result.Msg = "success"
	result.Data = tokenIds
	c.JSON(http.StatusOK, result)
}

func getDist20(c *gin.Context) {
	var result models.Result
	kid := c.Param("kid")
	if kid == "" {
		handleError(c, errors.New("invalid params"))
		return
	}
	dist, err := db.FindDist(kid, true)
	if err != nil {
		handleError(c, err)
		return
	}
	result.Code = http.StatusOK
	result.Msg = "success"
	result.Data = dist
	c.JSON(http.StatusOK, result)
}

func getDist721(c *gin.Context) {
	var result models.Result
	kid := c.Param("kid")
	if kid == "" {
		handleError(c, errors.New("invalid params"))
		return
	}
	dist, err := db.FindDist(kid, false)
	if err != nil {
		handleError(c, err)
		return
	}
	result.Code = http.StatusOK
	result.Msg = "success"
	result.Data = dist
	c.JSON(http.StatusOK, result)
}
