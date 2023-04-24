package apis

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ini8labs/lsdb"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	errMinAmount error = errors.New("minimum amount is 1")
)

type EventDate struct {
	EventDate primitive.DateTime `json:"date"`
}

type UserID struct {
	UID primitive.ObjectID `json:"id"`
}

type UserGovtID struct {
	GovtID string `json:"govId"`
}

func NewServer(addr string, log *logrus.Logger) error {

	r := gin.Default()

	// API end points
	r.GET("/api/v1/event/eventdata", GetAllEvents)

	r.POST("/api/v1/bet/new_bet", PlaceBets)
	r.POST("/api/v1/bet/update", UpdateBets)
	r.POST("/api/v1/bet/user_bets", UserBets)
	r.POST("/api/v1/bet/history", EventHistory)
	r.POST("/api/v1/event/eventdata_bydate", GetEventsByDate)
	r.POST("/api/v1/user_info/userinfo_ID", GetUserInfoByID)
	r.POST("/api/v1/user_info/new_user", NewUserInfo)
	r.POST("/api/v1/user_info/update_info", UpdateUserInfo)

	r.DELETE("/api/v1/bet/delete", DeleteBets)
	r.DELETE("/api/v1/user_info/delete", DeleteUserInfoByID)

	return r.Run(addr)
}

// ----User Beting info-----------
func PlaceBets(c *gin.Context) {
	dbClient, err := lsdb.NewClient()
	if err != nil {
		panic(err.Error())
	}

	var eventParticipantInfo lsdb.EventParticipantInfo
	if err := c.ShouldBind(&eventParticipantInfo); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	if err := dbClient.OpenConnection(); err != nil {
		panic(err.Error())
	}
	defer dbClient.CloseConnection()

	if err := dbClient.AddUserBet(eventParticipantInfo); err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, "Bets Placed Successfully")
}

func UpdateBets(c *gin.Context) {
	dbClient, err := lsdb.NewClient()
	if err != nil {
		panic(err.Error())
	}

	var eventParticipantInfo lsdb.EventParticipantInfo
	if err := c.ShouldBind(&eventParticipantInfo); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	if err := dbClient.OpenConnection(); err != nil {
		panic(err.Error())
	}
	defer dbClient.CloseConnection()

	if err := dbClient.UpdateUserBet(eventParticipantInfo); err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, "Bet Updated Successfully")
}

func UserBets(c *gin.Context) {
	dbClient, err := lsdb.NewClient()
	if err != nil {
		panic(err.Error())
	}

	var eventParticipantInfo lsdb.EventParticipantInfo
	if err := c.ShouldBind(&eventParticipantInfo); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	if err := dbClient.OpenConnection(); err != nil {
		panic(err.Error())
	}
	defer dbClient.CloseConnection()

	resp, err := dbClient.GetUserBets(eventParticipantInfo.UserID)
	if err != nil {
		panic(err.Error())
	}
	c.JSON(http.StatusOK, resp)
}

func DeleteBets(c *gin.Context) {
	dbClient, err := lsdb.NewClient()
	if err != nil {
		panic(err.Error())
	}

	var eventParticipantInfo lsdb.EventParticipantInfo
	if err := c.ShouldBind(&eventParticipantInfo); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	if err := dbClient.OpenConnection(); err != nil {
		panic(err.Error())
	}
	defer dbClient.CloseConnection()

	if err := dbClient.DeleteUserBet(eventParticipantInfo.BetUID); err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, "Bet Deleted Successfully")
}

func EventHistory(c *gin.Context) {
	dbClient, err := lsdb.NewClient()
	if err != nil {
		panic(err.Error())
	}

	var eventParticipantInfo lsdb.EventParticipantInfo
	if err := c.ShouldBind(&eventParticipantInfo); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	if err := dbClient.OpenConnection(); err != nil {
		panic(err.Error())
	}
	defer dbClient.CloseConnection()

	resp, err := dbClient.GetParticipantsInfoByEventID(eventParticipantInfo.EventUID)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, resp)
}

// ----User Events info-----
func GetAllEvents(c *gin.Context) {
	dbClient, err := lsdb.NewClient()
	if err != nil {
		panic(err.Error())
	}

	if err := dbClient.OpenConnection(); err != nil {
		panic(err.Error())
	}
	defer dbClient.CloseConnection()

	resp, err := dbClient.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func GetEventsByDate(c *gin.Context) {
	var eventdate EventDate
	dbClient, err := lsdb.NewClient()
	if err != nil {
		panic(err.Error())
	}

	if err := dbClient.OpenConnection(); err != nil {
		panic(err.Error())
	}
	defer dbClient.CloseConnection()

	if err := c.ShouldBind(&eventdate); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	resp, err := dbClient.GetEventsByDate(eventdate.EventDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		return
	}
	c.JSON(http.StatusOK, resp)
}

// ----User personal info----------
func NewUserInfo(c *gin.Context) {
	var userInfo lsdb.UserInfo
	dbClient, err := lsdb.NewClient()
	if err != nil {
		panic(err.Error())
	}

	if err := dbClient.OpenConnection(); err != nil {
		panic(err.Error())
	}
	defer dbClient.CloseConnection()

	if err := c.ShouldBind(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	err1 := dbClient.AddNewUserInfo(userInfo)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		return
	}
	c.JSON(http.StatusOK, "User Info Added Successfully")
}

type UpdateInfoStruct struct {
	UserID primitive.ObjectID `bson:"userid,omitempty"`
	Key    string             `bson:"key,omitempty"`
	Value  string             `bson:"value,omitempty"`
}

func UpdateUserInfo(c *gin.Context) {
	var userInfo UpdateInfoStruct
	dbClient, err := lsdb.NewClient()
	if err != nil {
		panic(err.Error())
	}

	if err := dbClient.OpenConnection(); err != nil {
		panic(err.Error())
	}
	defer dbClient.CloseConnection()

	if err := c.ShouldBind(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	err1 := dbClient.UpdateUserInfo(userInfo.UserID, userInfo.Key, userInfo.Value)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		return
	}
	c.JSON(http.StatusOK, "User Updated Successfully")
}

func GetUserInfoByID(c *gin.Context) {
	var userid UserID
	dbClient, err := lsdb.NewClient()
	if err != nil {
		panic(err.Error())
	}

	if err := dbClient.OpenConnection(); err != nil {
		panic(err.Error())
	}
	defer dbClient.CloseConnection()

	if err := c.ShouldBind(&userid); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	resp, err := dbClient.GetUserInfoByID(userid.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func DeleteUserInfoByID(c *gin.Context) {
	var userinfo lsdb.UserInfo
	dbClient, err := lsdb.NewClient()
	if err != nil {
		panic(err.Error())
	}

	if err := dbClient.OpenConnection(); err != nil {
		panic(err.Error())
	}
	defer dbClient.CloseConnection()

	if err := c.ShouldBind(&userinfo); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	err1 := dbClient.DeleteUserInfo(userinfo)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		return
	}
	c.JSON(http.StatusOK, "User info Deleted successfully")
}
