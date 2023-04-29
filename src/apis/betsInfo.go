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

type PlaceBetsStruct struct {
	UserID     primitive.ObjectID `bson:"user_id,omitempty"`
	EventUID   primitive.ObjectID `bson:"event_id,omitempty"`
	BetNumbers []int              `bson:"bet_numbers,omitempty"`
	Amount     int                `bson:"amount,omitempty"`
}

var getUserResp []UserBetsInfo

func (s Server) PlaceBets(c *gin.Context) {
	var placeBetsStruct PlaceBetsStruct
	if err := c.ShouldBind(&placeBetsStruct); err != nil {
		s.Logger.Error("bad format")
		c.JSON(http.StatusBadRequest, "bad Format")
		return
	}

	userIDValidated, err := validateUserID(placeBetsStruct.UserID.Hex())
	if err != nil {
		s.Logger.Error(errInvalidUserID)
		c.JSON(http.StatusBadRequest, errInvalidUserID)
		return
	}

	eventUIDValidated, err := validateEventID(placeBetsStruct.EventUID.Hex())
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

	eventParticipantInfo := lsdb.EventParticipantInfo{
		EventUID: eventUIDValidated,
		ParticipantInfo: lsdb.ParticipantInfo{
			UserID:     userIDValidated,
			BetNumbers: betNumbersvalidated,
			Amount:     amountValidated,
		},
	}

	if err := s.Client.AddUserBet(eventParticipantInfo); err != nil {
		s.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, "Bets Placed Successfully")
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

func (s Server) UpdateBets(c *gin.Context) {
	var userBetsInfo UserBetsInfo
	if err := c.ShouldBind(&userBetsInfo); err != nil {
		s.Logger.Error("bad format")
		c.JSON(http.StatusBadRequest, "bad Format")
		return
	}

	fmt.Println(userBetsInfo.BetNumbers)
	betUIDValidated, err := validateBetUID(userBetsInfo.BetUID.Hex())
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

/**
 * @Summary Delete a user
 * @Description Delete a user from the system by ID
 * @Param BetUID path int true "BetUID"
 * @Produce json
 * @Success 204 "No content"
 * @Failure 404 {string} string "User not found"
 * @Router /bet/delete/{betuid} [delete]
 */
func (s Server) DeleteBets(c *gin.Context) {
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

// func (s Server) EventHistory(c *gin.Context) {
// 	eventUID, exists1 := c.GetQuery("eventuid")
// 	if !exists1 {
// 		s.Logger.Error(errIncorrectField)
// 		c.JSON(http.StatusBadRequest, "Bad Format")
// 		return
// 	}

// 	eventUIDConv, err := primitive.ObjectIDFromHex(eventUID)
// 	if err != nil {
// 		s.Logger.Errorf("error validateBetnumbersing string to HEX: %s", err.Error())
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
	userID := c.Param("bets")

	userIDConv, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, (err.Error()))
	}

	resp, err := s.Client.GetUserBets(userIDConv)
	if err != nil {
		s.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	getUserResp = requiredInfoUserBets(resp)
	c.JSON(http.StatusOK, getUserResp)
}
