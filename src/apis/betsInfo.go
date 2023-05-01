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
	BetUID     primitive.ObjectID `bson:"_id,omitempty"`
	BetNumbers []int              `bson:"bet_numbers,omitempty"`
	Amount     int                `bson:"amount,omitempty"`
}

type PlaceBetsStruct struct {
	UserID     string `bson:"user_id,omitempty"`
	EventUID   string `bson:"event_id,omitempty"`
	BetNumbers []int  `bson:"bet_numbers,omitempty"`
	Amount     int    `bson:"amount,omitempty"`
}

var respConv []UserBetsInfo

var eventParticipantInfo lsdb.EventParticipantInfo

func (s Server) placeBets(c *gin.Context) {
	var placeBetsStruct PlaceBetsStruct
	if err := c.ShouldBind(&placeBetsStruct); err != nil {
		s.Logger.Error("bad format")
		c.JSON(http.StatusBadRequest, "bad Format")
		return
	}

	userIDValidated, err := validateID(string(placeBetsStruct.UserID))
	if err != nil {
		s.Logger.Error(errInvalidUserID)
		c.JSON(http.StatusBadRequest, errInvalidUserID)
		return
	}

	eventUIDValidated, err := validateID(placeBetsStruct.EventUID)
	if err != nil {
		s.Logger.Error(errInvalidEventID)
		c.JSON(http.StatusBadRequest, errInvalidEventID)
		return
	}

	amountValidated, err := validateAmount(placeBetsStruct.Amount)
	errHandle(err)

	betNumbersvalidated, err := validateBetnumbers(placeBetsStruct.BetNumbers)
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
	if err := c.ShouldBind(&userBetsInfo); err != nil {
		s.Logger.Error("bad format")
		c.JSON(http.StatusBadRequest, "bad Format")
		return
	}

	fmt.Println(userBetsInfo.BetNumbers)
	betUIDValidated, err := validateID(userBetsInfo.BetUID.Hex())
	errHandle(err)

	betNumbersValidated, err := validateBetnumbers(userBetsInfo.BetNumbers)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	amountValidated, err := validateAmount(userBetsInfo.Amount)
	fmt.Println(amountValidated)
	errHandle(err)

	eventParticipantInfo := lsdb.EventParticipantInfo{
		BetUID: betUIDValidated,
		ParticipantInfo: lsdb.ParticipantInfo{
			BetNumbers: betNumbersValidated,
			Amount:     amountValidated,
		},
	}

	if err := s.Client.UpdateUserBet(eventParticipantInfo); err != nil {
		s.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusCreated, "Bet Updated Successfully")
}

func (s Server) deleteBets(c *gin.Context) {
	betUID := c.Param("betuid")

	bettUIDConv, err := primitive.ObjectIDFromHex(betUID)
	if err != nil {
		s.Logger.Error("Bad BetUID")
		c.JSON(http.StatusBadRequest, "Bad format")
		return
	}

	if err := s.Client.DeleteUserBet(bettUIDConv); err != nil {
		s.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, "Bet Deleted Successfully")
}

func (s Server) betsHistorybyEvent(c *gin.Context) {
	eventUID := c.Param("eventuid")

	eventUIDConv, err := validateID(eventUID)
	if err != nil {
		s.Logger.Errorf("error validateBetnumbersing string to HEX: %s", err.Error())
		c.JSON(http.StatusBadRequest, errInvalidEventID.Error())
		return
	}

	resp, err := s.Client.GetParticipantsInfoByEventID(eventUIDConv)
	if err != nil {
		s.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	respConv = requiredInfoUserBets(resp)
	c.JSON(http.StatusOK, respConv)
}

func (s Server) userBets(c *gin.Context) {
	userID := c.Param("bets")

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

func errHandle(err error) {
	var s Server
	var c *gin.Context
	if err != nil {
		s.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
}
