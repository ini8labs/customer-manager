package apis

import (
	// "errors"

	"net/http"

	//"strings"

	"github.com/gin-gonic/gin"
	"github.com/ini8labs/lsdb"

	// "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventPartInfo struct {
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

var getUserResp []UserBetsInfo

func (s Server) PlaceBets(c *gin.Context) {
	eventUID, exists1 := c.GetQuery("eventuid")
	userID, exists2 := c.GetQuery("userid")
	betNumbers, exists3 := c.GetQuery("betnumbers")
	Amount, exists4 := c.GetQuery("amount")

	// check for no missing fields
	if !exists1 || !exists2 || !exists3 || !exists4 {
		s.Logger.Error("Field empty")
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	eventUIDConv, err := primitive.ObjectIDFromHex(eventUID)
	if err != nil {
		s.Logger.Error(errInvalidEventID)
		c.JSON(http.StatusBadRequest, errInvalidEventID)
		return
	}
	userIDConv, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		s.Logger.Error(errInvalidUserID)
		c.JSON(http.StatusBadRequest, errInvalidUserID)
		return
	}
	amountConv, err := amountCheck(Amount, c)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	betNumbersConv, err := convert(betNumbers)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	eventParticipantInfo := lsdb.EventParticipantInfo{
		EventUID: eventUIDConv,
		ParticipantInfo: lsdb.ParticipantInfo{
			UserID:     userIDConv,
			BetNumbers: betNumbersConv,
			Amount:     amountConv,
		},
	}

	if err := s.Client.AddUserBet(eventParticipantInfo); err != nil {
		s.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, "Bets Placed Successfully")
}

func (s Server) UpdateBets(c *gin.Context) {
	betUID, exists1 := c.GetQuery("betuid")
	betNumbers, exists2 := c.GetQuery("betnumbers")
	Amount, exists3 := c.GetQuery("amount")

	if !exists1 || !exists2 || !exists3 {
		s.Logger.Error(errIncorrectField)
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	bettUIDConv, _ := primitive.ObjectIDFromHex(betUID)
	betNumbersConv, err := convert(betNumbers)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Bad format")
		return
	}

	amountConv, err := amountCheck(Amount, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	eventParticipantInfo := lsdb.EventParticipantInfo{
		BetUID: bettUIDConv,
		ParticipantInfo: lsdb.ParticipantInfo{
			BetNumbers: betNumbersConv,
			Amount:     amountConv,
		},
	}

	if err := s.Client.UpdateUserBet(eventParticipantInfo); err != nil {
		s.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, "Bet Updated Successfully")
}

func (s Server) DeleteBets(c *gin.Context) {
	betUID, exists1 := c.GetQuery("betuid")
	if !exists1 {
		s.Logger.Error(errIncorrectField)
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

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

// func (s Server) EventHistory(c *gin.Context) {
// 	eventUID, exists1 := c.GetQuery("eventuid")
// 	if !exists1 {
// 		s.Logger.Error(errIncorrectField)
// 		c.JSON(http.StatusBadRequest, "Bad Format")
// 		return
// 	}

// 	eventUIDConv, err := primitive.ObjectIDFromHex(eventUID)
// 	if err != nil {
// 		s.Logger.Errorf("error converting string to HEX: %s", err.Error())
// 		c.JSON(http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	resp, err := s.Client.GetParticipantsInfoByEventID(eventUIDConv)
// 	if err != nil {
// 		s.Logger.Error(err.Error())
// 		c.JSON(http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	respSlice = s.RequiredInfo(resp)
// 	c.JSON(http.StatusOK, respSlice)
// }

func (s Server) UserBets(c *gin.Context) {
	userID, exists := c.GetQuery("userid")
	if !exists {
		s.Logger.Error(errIncorrectField)
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	userIDConv, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		s.Logger.Errorf("error converting string to HEX: %s", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := s.Client.GetUserBets(userIDConv)
	if err != nil {
		s.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	getUserResp = requiredInfo(resp)
	c.JSON(http.StatusOK, getUserResp)
}
