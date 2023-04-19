package apis

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Struct for Placing betsData
type Bet struct {
	Numbers []int `json:"numbers"`
	Amount  int   `json:"amount"`
}

// Map to store Users DataBase
// var users = map[string]string{
// 	"user1": "password1",
// 	"user2": "password2",
// }

func NewServer(addr string, log *logrus.Logger) error {

	r := gin.Default()

	// API end point
	r.POST("/api/v1/bet_number", PlaceBet)

	return r.Run(addr)

}

var userData = make(map[string]Bet)

func PlaceBet(c *gin.Context) {
	var requestData Bet
	if err := c.ShouldBind(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	fmt.Println(c.GetHeader("cookie_1"))
	userData[c.GetHeader("cookie_1")] = requestData

	BettingDBFunc(requestData)

	// fmt.Println(data["User1"].Numbers)

	c.JSON(http.StatusOK, "Bet Placed")

	fmt.Println(userData)

}

var betsData = make(map[int]int)

func BettingDBFunc(r Bet) {
	for _, val := range r.Numbers {
		temp := betsData[val]
		betsData[val] = r.Amount + temp
	}
	fmt.Println(betsData)
}

// func DelBet(c *gin.Context) {
// 	var requestData DelStr
// 	if err := c.ShouldBind(&requestData); err != nil {
// 		c.JSON(http.StatusBadRequest, "Bad Format")
// 		return
// 	}

// 	betsData[requestData.Number] = 0
// 	c.JSON(http.StatusOK, "Bet Deleted")
// 	fmt.Println(betsData)

// }
