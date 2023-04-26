package apis

import (
	// "errors"
	"net/http"
	"strconv"

	//"strings"

	"github.com/gin-gonic/gin"
	"github.com/ini8labs/lsdb"

	// "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

	eventUIDConv, _ := primitive.ObjectIDFromHex(eventUID)
	userIDConv, _ := primitive.ObjectIDFromHex(userID)
	amount, _ := strconv.Atoi(Amount)
	if amount < 1 {
		s.Logger.Error(errMinAmount)
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	betNumbersConv, err := s.Convert(betNumbers)
	if err != nil {
		s.Logger.Error("Bet numbers no proper")
		c.JSON(http.StatusBadRequest, "Bad format")
		return
	}

	eventParticipantInfo := lsdb.EventParticipantInfo{
		EventUID: eventUIDConv,
		ParticipantInfo: lsdb.ParticipantInfo{
			UserID:     userIDConv,
			BetNumbers: betNumbersConv,
			Amount:     amount,
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
	betNumbersConv, err := s.Convert(betNumbers)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Bad format")
		return
	}
	amount, _ := strconv.Atoi(Amount)
	if amount < 1 {
		s.Logger.Error(errMinAmount)
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	eventParticipantInfo := lsdb.EventParticipantInfo{
		BetUID: bettUIDConv,
		ParticipantInfo: lsdb.ParticipantInfo{
			BetNumbers: betNumbersConv,
			Amount:     amount,
		},
	}

	if err := s.Client.UpdateUserBet(eventParticipantInfo); err != nil {
		s.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, "Bet Updated Successfully")
}

func (s Server) UserBets(c *gin.Context) {
	userID, exists1 := c.GetQuery("userid")
	if !exists1 {
		s.Logger.Error(errIncorrectField)
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	userIDConv, _ := primitive.ObjectIDFromHex(userID)
	resp, err := s.Client.GetUserBets(userIDConv)
	if err != nil {
		s.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if len(resp) < 1 {
		s.Logger.Error("empty response")
		c.JSON(http.StatusNotFound, "Not found")
		return
	}
	c.JSON(http.StatusOK, resp)
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

// make changes
func (s Server) EventHistory(c *gin.Context) {
	eventUID, exists1 := c.GetQuery("eventuid")
	if !exists1 {
		s.Logger.Error(errIncorrectField)
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	eventUIDConv, err := primitive.ObjectIDFromHex(eventUID)
	if err != nil {
		s.Logger.Errorf("error converting string to HEX: %s", err.Error())
		c.JSON(http.StatusBadRequest, "invalid event UID")
		return
	}

	resp, err := s.Client.GetParticipantsInfoByEventID(eventUIDConv)
	if err != nil {
		s.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if len(resp) < 1 {
		s.Logger.Error("empty response")
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	// var respSlice []lsdb.EventParticipantInfo

	// for _, item := range resp {
	// 	item.CreatedAt="",
	// 	item.UpdatedAt="",
	// }

	c.JSON(http.StatusOK, resp)
}
