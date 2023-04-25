package apis

import (
	"errors"
	"net/http"
	"strconv"

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
	r.GET("/api/v1/bet/new_bet", server.PlaceBets1)
	r.POST("/api/v1/bet/update", server.UpdateBets)
	r.POST("/api/v1/bet/user_bets", server.UserBets)
	r.DELETE("/api/v1/bet/delete", server.DeleteBets)
	r.POST("/api/v1/bet/history", server.EventHistory)
	r.POST("/api/v1/user_info/userinfo_ID", server.GetUserInfoByID)
	r.GET("/api/v1/user_info/new_user", server.NewUserInfo)
	r.POST("/api/v1/user_info/update_info", server.UpdateUserInfo)
	r.DELETE("/api/v1/user_info/delete", server.DeleteUserInfoByID)
	r.GET("/api/v1/event/eventdata", server.GetAllEvents)
	// r.POST("/api/v1/event/eventdata_bydate", GetEventsByDate

	return r.Run(server.Addr)
}

// ----User Beting info-----------
func (s Server) PlaceBets1(c *gin.Context) {
	var eventParticipantInfo lsdb.EventParticipantInfo
	if err := c.ShouldBind(&eventParticipantInfo); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	if err := s.Client.AddUserBet(eventParticipantInfo); err != nil {
		s.Logger.Error(err.Error())
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
	name, exists1 := c.GetQuery("name")
	phone, exists2 := c.GetQuery("phone")
	govID, exists3 := c.GetQuery("govid")
	eMail, exists4 := c.GetQuery("email")

	// check for no missing fields
	if !exists1 || !exists2 || !exists3 || !exists4 {
		s.Logger.Error("Field empty")
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	phone_int64, _ := strconv.ParseInt(phone, 10, 64)

	UserInfo1 := lsdb.UserInfo{
		Name:  name,
		Phone: phone_int64,
		GovID: govID,
		EMail: eMail,
	}

	s.Logger.Println("after creating a struct ")
	s.Logger.Println(UserInfo1.Phone)
	err := s.Client.AddNewUserInfo(UserInfo1)
	if err != nil {
		s.Logger.Error("Internal server Error")
		c.JSON(http.StatusInternalServerError, "Something is wrong with the server")
		return
	}
	c.JSON(http.StatusOK, "User Info Added Successfully")
	s.Logger.Info("Create operation performed")
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
