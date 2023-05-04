package apis

import (
	"fmt"
	"net/http"

	//"strings"

	"github.com/gin-gonic/gin"
	"github.com/ini8labs/lsdb"

	// "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventParticipantInfo struct {
	BetUID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID     primitive.ObjectID `bson:"user_id,omitempty"`
	BetNumbers []int              `bson:"bet_numbers,omitempty"`
	Amount     int                `bson:"amount,omitempty"`
}

type UserBetsInfo struct {
	BetUID     string `bson:"_id,omitempty"`
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
	errHandle(err)

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
}

func (s Server) updateBets(c *gin.Context) {
	var userBetsInfo UserBetsInfo
	userID := c.Query("userid")

	if err := c.ShouldBind(&userBetsInfo); err != nil {
		s.Logger.Error("bad format")
		c.JSON(http.StatusBadRequest, "bad Format")
		return
	}

	resp, err := s.userBetsResp(userID)
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
	betUID := c.Param("id")
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

	if len(respSlice) < 1 {
		s.Logger.Error(errNoRecords)
		c.JSON(http.StatusNotFound, errNoRecords.Error())
		return
	}

	// getParticipantsInfoByEventIDLoop(resp1)
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
	s.Logger.Errorln(resp)
	if err != nil {
		s.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, errInvalidUserID)
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

func errHandle(err error) {
	var s Server
	var c *gin.Context
	if err != nil {
		s.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
}

// /user/{userid}/bets
