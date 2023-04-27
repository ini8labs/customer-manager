package apis

import (
	"strconv"
	"strings"

	//"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ini8labs/lsdb"
)

func convert(str string) ([]int, error) {
	split := strings.Split(str, ",")
	strToInt := []int{}

	if len(split) < 1 {
		return nil, errMinNumbers
	}

	if len(split) > 5 {
		return nil, errMaxNumbers
	}

	for _, i := range split {
		j, err := strconv.Atoi(i)
		if err != nil || j < 1 || j > 90 {
			return nil, errNumberNotAllowed
		}

		strToInt = append(strToInt, j)
	}
	return strToInt, nil
}

func amountCheck(amount string, c *gin.Context) (int, error) {
	amountInt, err := strconv.Atoi(amount)
	if err != nil {
		return 0, errInvalidAmount
	}
	if amountInt < 1 {
		return 0, errMinAmount
	}
	return amountInt, nil
}

func requiredInfo(resp []lsdb.EventParticipantInfo) []UserBetsInfo {
	var respSlice []UserBetsInfo
	for i := 0; i < len(resp); i++ {
		var tempResp UserBetsInfo
		tempResp.BetUID = resp[i].BetUID
		tempResp.Amount = resp[i].Amount
		tempResp.BetNumbers = resp[i].BetNumbers

		respSlice = append(respSlice, tempResp)
	}
	return respSlice
}

// func (s Server)badRequestError(err error, c *gin.Context) {
// 	if err != nil {
// 		s.Logger.Error(err)
// 		c.JSON(http.StatusBadRequest, err)
// 		return
// 	}
// }
