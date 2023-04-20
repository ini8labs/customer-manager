package apis

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	errMinAmount error = errors.New("minimum amount is moire h")
)

// Struct for Placing betsData
type Bet struct {
	UserInfo string `json:"userinfo"`
	Numbers  []int  `json:"numbers"`
	Amount   int    `json:"amount"`
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

type mapInput struct {
	Numbers []int
	Amount  int
}

var userData = make(map[string]mapInput)

func PlaceBet(c *gin.Context) {
	var (
		requestData Bet
		mapData     mapInput
	)
	if err := c.ShouldBind(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	if err := validateBetPlaceInput(requestData); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	mapData.Numbers, mapData.Amount = requestData.Numbers, requestData.Amount

	userData[requestData.UserInfo] = mapData

	fmt.Println(userData[requestData.UserInfo])

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

func validateBetPlaceInput(requestData Bet) error {
	if len(requestData.Numbers) < 1 {
		return errMinAmount
	}
	if requestData.Amount <= 0 {
		return fmt.Errorf("amount can not be 0 or negative")
	}
	return nil
}
