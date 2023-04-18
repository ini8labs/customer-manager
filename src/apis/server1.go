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
		bets[i] = 12
	}
	// API end point
	r.POST("/api/v1/bet_number", PlaceBet)
	r.GET("/api/v1/bet_number", Welcome)

	return r.Run(addr)

}

func PlaceBet(c *gin.Context) {
	var requestData Bet
	if err := c.ShouldBind(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	fmt.Println(requestData.Number)
	fmt.Println(requestData.Amount)

	c.JSON(http.StatusOK, "all good")

}

var supportedLanguages = []string{"eng", "twi"}

func Welcome(c *gin.Context) {
	languageSupported := false
	lang := c.Param("language")

	for _, language := range supportedLanguages {
		if language == lang {
			languageSupported = true
			break
		}
	}

	if !languageSupported {
		errMessage, _ := fmt.Printf("%s is not a supported language", lang)
		c.JSON(http.StatusNotAcceptable, errMessage)
		return
	}

	resp := "okok"

	c.JSON(http.StatusOK, resp)
}
