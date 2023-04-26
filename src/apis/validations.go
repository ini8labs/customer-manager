package apis

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s Server) Convert(str string) ([]int, error) {
	split := strings.Split(str, ",")
	strToInt := []int{}

	if len(split) < 1 {
		s.Logger.Error(errMinNumbers)
		return nil, errMinNumbers
	}

	if len(split) > 5 {
		s.Logger.Error(errMaxNumbers)
		return nil, errMaxNumbers
	}

	for _, i := range split {
		j, err := strconv.Atoi(i)
		if err != nil || j < 1 || j > 90 {
			s.Logger.Error(errNumberNotAllowed)
			return nil, errNumberNotAllowed
		}

		strToInt = append(strToInt, j)
	}
	return strToInt, nil
}

func (s Server) AmountCheck(amount string, c *gin.Context) (int, error) {
	amountInt, err := strconv.Atoi(amount)
	if err != nil {
		s.Logger.Error("amount invalid")
		c.JSON(http.StatusBadRequest, "give valid amount")
		return 0, err
	}
	if amountInt < 1 {
		s.Logger.Error(errMinAmount)
		c.JSON(http.StatusBadRequest, "Bad Format")
		return 0, errMinAmount
	}
	return amountInt, nil
}
