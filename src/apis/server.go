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

type Server struct {
	*logrus.Logger
	*lsdb.Client
	Addr string
}

type EventDate struct {
	EventDate primitive.DateTime `json:"date"`
}

type UserID struct {
	UID primitive.ObjectID `json:"id"`
}

type UserGovtID struct {
	GovtID string `json:"govId"`
}

type UpdateInfoStruct struct {
	UserID primitive.ObjectID `bson:"userid,omitempty"`
	Key    string             `bson:"key,omitempty"`
	Value  string             `bson:"value,omitempty"`
}

func NewServer(server Server) error {

	r := gin.Default()

	// API end points
	r.POST("/api/v1/bet?new_bet", server.PlaceBets)
	r.POST("/api/v1/bet/update", server.UpdateBets)
	r.POST("/api/v1/bet/user_bets", server.UserBets)
	r.DELETE("/api/v1/bet/delete", server.DeleteBets)
	r.POST("/api/v1/bet/history", server.EventHistory)
	r.POST("/api/v1/user_info/userinfo_ID", server.GetUserInfoByID)
	r.POST("/api/v1/user_info/new_user", server.NewUserInfo)
	r.POST("/api/v1/user_info/update_info", server.UpdateUserInfo)
	r.DELETE("/api/v1/user_info/delete", server.DeleteUserInfoByID)
	r.GET("/api/v1/event/eventdata", server.GetAllEvents)
	// r.POST("/api/v1/event/eventdata_bydate", GetEventsByDate

	return r.Run(server.Addr)
}

// ----User Beting info-----------
func (s Server) PlaceBets(c *gin.Context) {
	var eventParticipantInfo lsdb.EventParticipantInfo
	if err := c.ShouldBind(&eventParticipantInfo); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	if err := s.Client.AddUserBet(eventParticipantInfo); err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, "Bets Placed Successfully")
}

func (s Server) UpdateBets(c *gin.Context) {
	var eventParticipantInfo lsdb.EventParticipantInfo
	if err := c.ShouldBind(&eventParticipantInfo); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	if err := s.Client.UpdateUserBet(eventParticipantInfo); err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, "Bet Updated Successfully")
}

func (s Server) UserBets(c *gin.Context) {
	var eventParticipantInfo lsdb.EventParticipantInfo
	if err := c.ShouldBind(&eventParticipantInfo); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	resp, err := s.Client.GetUserBets(eventParticipantInfo.UserID)
	if err != nil {
		panic(err.Error())
	}
	c.JSON(http.StatusOK, resp)
}

func (s Server) DeleteBets(c *gin.Context) {
	var eventParticipantInfo lsdb.EventParticipantInfo
	if err := c.ShouldBind(&eventParticipantInfo); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	if err := s.Client.DeleteUserBet(eventParticipantInfo.BetUID); err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, "Bet Deleted Successfully")
}

func (s Server) EventHistory(c *gin.Context) {
	var eventParticipantInfo lsdb.EventParticipantInfo
	if err := c.ShouldBind(&eventParticipantInfo); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	resp, err := s.Client.GetParticipantsInfoByEventID(eventParticipantInfo.EventUID)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, resp)
}

// ----User Events info-----
func (s Server) GetAllEvents(c *gin.Context) {

	resp, err := s.Client.GetAllEvents()
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
func (s Server) NewUserInfo(c *gin.Context) {
	var userInfo lsdb.UserInfo
	if err := c.ShouldBind(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	err1 := s.Client.AddNewUserInfo(userInfo)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		return
	}
	c.JSON(http.StatusOK, "User Info Added Successfully")
}

func (s Server) UpdateUserInfo(c *gin.Context) {
	var userInfo UpdateInfoStruct
	if err := c.ShouldBind(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	err1 := s.Client.UpdateUserInfo(userInfo.UserID, userInfo.Key, userInfo.Value)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		return
	}
	c.JSON(http.StatusOK, "User Updated Successfully")
}

func (s Server) GetUserInfoByID(c *gin.Context) {
	var userid UserID
	if err := c.ShouldBind(&userid); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	resp, err := s.Client.GetUserInfoByID(userid.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (s Server) DeleteUserInfoByID(c *gin.Context) {
	var userinfo lsdb.UserInfo

	if err := c.ShouldBind(&userinfo); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	err1 := s.Client.DeleteUserInfo(userinfo)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		return
	}
	c.JSON(http.StatusOK, "User info Deleted successfully")
}
