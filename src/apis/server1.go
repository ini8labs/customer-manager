package apis

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Struct for Placing Bets

type Bet struct {
	Number int `json:"number"`
	Amount int `json:"amount"`
}

var bets = make(map[int]int)

func NewServer(addr string, log *logrus.Logger) error {

	r := gin.Default()

	for i := 1; i < 91; i++ {
		bets[i] = 0
	}
	// API end point
	r.POST("/api/v1/bet_number", PlaceBet)

	return r.Run(addr)

}

func PlaceBet(c *gin.Context) {
	var requestData Bet
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	fmt.Println(requestData.Number)
	fmt.Println(requestData.Amount)

	if requestData.Number > 91 {
		c.JSON(http.StatusBadRequest, "Bad request")
		return
	}

	bets[requestData.Number] = requestData.Amount
	fmt.Println(bets)

	c.JSON(http.StatusOK, "Bet Placed")

}
