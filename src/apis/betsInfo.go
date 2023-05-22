package apis

import (
	"encoding/json"
	"fmt"

	"net/http"
	"os"
	"strconv"

	//"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"

	// "github.com/sirupsen/logrus"
	// "go.mongodb.org/mongo-driver/bson/primitive"

	//"github.com/ini8labs/sns/src/apis"
	"github.com/ini8labs/lsdb"
)

type UserBetsInfo struct {
	BetUID     string `bson:"bet_id,omitempty"`
	BetNumbers []int  `bson:"bet_numbers,omitempty"`
	Amount     int    `bson:"amount,omitempty"`
}

type UserBetsInfoByEvent struct {
	BetNumbers []int `bson:"bet_numbers,omitempty"`
	Amount     int   `bson:"amount,omitempty"`
}

type NewBetsFomat struct {
	UserID     string `bson:"user_id,omitempty"`
	EventUID   string `bson:"event_id,omitempty"`
	BetNumbers []int  `bson:"bet_numbers,omitempty"`
	Amount     int    `bson:"amount,omitempty"`
}

var respConv []UserBetsInfo

var eventsList = []string{"MS", "LT", "MW", "FT", "FB", "NW"}

var eventParticipantInfo lsdb.EventParticipantInfo

var respSlice []UserBetsInfoByEvent

func (s Server) placeBets(c *gin.Context) {
	var NewBetsFomat NewBetsFomat
	if err := c.ShouldBind(&NewBetsFomat); err != nil {
		s.Logger.Error("bad format")
		c.JSON(http.StatusBadRequest, "bad Format")
		return
	}

	userIDValidated, err := validateID(string(NewBetsFomat.UserID))
	if err != nil {
		s.Logger.Error(errInvalidUserID)
		c.JSON(http.StatusBadRequest, errInvalidUserID)
		return
	}

	resp, err := s.eventsAvailable()
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println(resp)

	eventUIDValidated, err := validateEventID(NewBetsFomat.EventUID, resp)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	amountValidated, err := validateAmount(NewBetsFomat.Amount)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	betNumbersvalidated, err := validateBetnumbers(NewBetsFomat.BetNumbers)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// eventParticipantInfo := lsdb.EventParticipantInfo{
	// 	EventUID: eventUIDValidated,
	// 	ParticipantInfo: lsdb.ParticipantInfo{
	// 		UserID:     userIDValidated,
	// 		BetNumbers: betNumbersvalidated,
	// 		Amount:     amountValidated,
	// 	},
	// }

	eventParticipantInfo.UserID = userIDValidated
	eventParticipantInfo.EventUID = eventUIDValidated
	eventParticipantInfo.Amount = amountValidated
	eventParticipantInfo.BetNumbers = betNumbersvalidated
	if err := s.Client.AddUserBet(eventParticipantInfo); err != nil {
		s.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, "Bets Placed Successfully")

	var userInfo *lsdb.UserInfo
	userInfo, _ = s.Client.GetUserInfoByID(eventParticipantInfo.UserID)
	tempFrom := "+12707477263"
	tempName := userInfo.Name
	tempPhone := ("+91" + strconv.Itoa(int(userInfo.Phone)))
	tempMessage := ("Dear" + " " + tempName + " " + "your bet was placed successfully.")
	SMS(tempFrom, tempPhone, tempMessage)

}

var Client *twilio.RestClient

func SMS(from, to, message string) {
	if err := godotenv.Load(); err != nil {
		panic(err.Error())
	}

	accountSid := os.Getenv("A_SID")
	authToken := os.Getenv("AUTH_TOKEN")

	twilioClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(from)
	params.SetBody(message)

	resp, err := twilioClient.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Error sending SMS message: " + err.Error())
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println("Response: " + string(response))
	}
}

func (s Server) updateBets(c *gin.Context) {
	var userBetsInfo UserBetsInfo
	userID := c.Query("userid")

	if err := c.ShouldBind(&userBetsInfo); err != nil {
		s.Logger.Error("bad format")
		c.JSON(http.StatusBadRequest, "bad Format")
		return
	}

	resp, _ := s.userBetsResp(userID)
	betUIDValidated, err := validateBetUID(userBetsInfo.BetUID, resp)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, errInvalidBetUID.Error())
		return
	}

	betNumbersValidated, err := validateBetnumbers(userBetsInfo.BetNumbers)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, errNumberNotAllowed.Error())
		return
	}

	amountValidated, err := validateAmount(userBetsInfo.Amount)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, errInvalidAmount.Error())
		return
	}

	eventParticipantInfo.BetUID = betUIDValidated
	eventParticipantInfo.Amount = amountValidated
	eventParticipantInfo.BetNumbers = betNumbersValidated

	if err := s.Client.UpdateUserBet(eventParticipantInfo); err != nil {
		s.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusCreated, "Bet Updated Successfully")
}

func (s Server) deleteBets(c *gin.Context) {
	betUID := c.Param("betuid")
	userID := c.Query("userid")

	resp, _ := s.userBetsResp(userID)
	bettUIDConv, err := validateBetUID(betUID, resp)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := s.Client.DeleteUserBet(bettUIDConv); err != nil {
		s.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, "Bet Deleted Successfully")
}

func (s Server) betsHistoryByEventType(c *gin.Context) {
	eventType := c.Query("type")
	userID := c.Query("userid")

	if err := validateEventType(eventType); err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp1, err := s.Client.GetEventsByType(eventType)
	if err != nil {
		s.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	for i, _ := range resp1 {
		EventUID := resp1[i].EventUID
		resp2, err := s.Client.GetParticipantsInfoByEventID(EventUID)
		if err != nil {
			s.Logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		requiredBetsByEventType(resp2, userID)
	}

	if respSlice == nil {
		s.Logger.Error(errNoRecords)
		c.JSON(http.StatusNotFound, errNoRecords.Error())
		return
	}

	c.JSON(http.StatusOK, respSlice)
	respSlice = []UserBetsInfoByEvent{}
}

func (s Server) userBets(c *gin.Context) {
	userID := c.Param("userid")

	userIDConv, err := validateID(userID)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, errInvalidUserID.Error())
		return
	}

	resp, err := s.Client.GetUserBets(userIDConv)
	if err != nil {
		s.Logger.Error(err.Error())
		c.JSON(http.StatusNotFound, errInvalidUserID.Error())
		return
	}

	if resp == nil {
		s.Logger.Error(errEmptyResp)
		c.JSON(http.StatusNotFound, errEmptyResp.Error())
		return
	}

	respConv = requiredInfoUserBets(resp)
	c.JSON(http.StatusOK, respConv)
}

func (s Server) userBetsResp(str string) ([]UserBetsInfo, error) {

	userIDConv, _ := validateID(str)

	resp, err := s.Client.GetUserBets(userIDConv)
	if err != nil {
		return []UserBetsInfo{}, err
	}

	respConv = requiredInfoUserBets(resp)
	return respConv, nil
}

// func errHandle(err error) {
// 	var s Server
// 	var c *gin.Context
// 	if err != nil {
// 		s.Logger.Error(err.Error())
// 		c.JSON(http.StatusBadRequest, err.Error())
// 		return
// 	}
// }
