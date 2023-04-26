package apis

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ini8labs/lsdb"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	errIncorrectField   error = errors.New("incorrect field or fields")
	errMinNumbers       error = errors.New("minimum  number that can be selected is 1")
	errMaxNumbers       error = errors.New("maximum numbers that can be selected are  5")
	errMinAmount        error = errors.New("minimum amount that can be placed is 1")
	errNumberNotAllowed error = errors.New("bet numbers should be between 1 and 90")
	errIncorrectPhoneNo error = errors.New("phone number is entered incorrectly")
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

type UserInfoStruct struct {
	UID   primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name,omitempty"`
	Phone int64              `bson:"phone,omitempty"`
	EMail string             `bson:"e_mail,omitempty"`
}

type EventsInfo struct {
	EventUID  primitive.ObjectID `bson:"_id,omitempty"`
	EventDate primitive.DateTime `bson:"event_date,omitempty"`
	EventName string             `bson:"name,omitempty"`
	EventType string             `bson:"event_type,omitempty"`
}

func NewServer(server Server) error {

	r := gin.Default()

	// API end points
	r.POST("/api/v1/bet/new_bet", server.PlaceBets)                // done,
	r.PUT("/api/v1/bet/update", server.UpdateBets)                 // validation
	r.GET("/api/v1/bet/user_bets", server.UserBets)                // required info
	r.DELETE("/api/v1/bet/delete", server.DeleteBets)              // done
	r.GET("/api/v1/bet/history", server.EventHistory)              //done
	r.GET("/api/v1/user_info/userinfo_ID", server.GetUserInfoByID) // done
	r.POST("/api/v1/user_info/new_user", server.NewUserInfo)       // needs to add validation
	//r.PUT("/api/v1/user_info/update_info", server.UpdateUserInfo)
	r.DELETE("/api/v1/user_info/delete", server.DeleteUserInfoByID)
	r.GET("/api/v1/event/eventdata", server.GetAllEvents)
	// r.POST("/api/v1/event/eventdata_bydate", GetEventsByDate

	return r.Run(server.Addr)
}

// ----User Beting info-----------
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

func (s Server) Convert(str string) ([]int, error) {
	split := strings.Split(str, ",")
	strToInt := []int{}

	if len(split) < 1 {
		s.Logger.Error(errMinNumbers)
		return nil, errMinNumbers
	}

	if len(split) > 5 {
		s.Logger.Error(errMaxNumbers)
		return nil, errMaxNumbers
	}

	for _, i := range split {
		j, err := strconv.Atoi(i)
		if err != nil || j < 1 || j > 90 {
			s.Logger.Error(errNumberNotAllowed)
			return nil, errNumberNotAllowed
		}

		strToInt = append(strToInt, j)
	}
	return strToInt, nil
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

func (s Server) EventHistory(c *gin.Context) {
	eventUID, exists1 := c.GetQuery("eventuid")
	if !exists1 {
		s.Logger.Error(errIncorrectField)
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	eventUIDConv, _ := primitive.ObjectIDFromHex(eventUID)

	resp, err := s.Client.GetParticipantsInfoByEventID(eventUIDConv)
	if err != nil || len(resp) < 1 {
		s.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if len(resp) < 1 {
		s.Logger.Error("empty response")
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

// ----User Events info-----
func (s Server) GetAllEvents(c *gin.Context) {

	resp, err := s.Client.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		logrus.Infoln(err)
		return
	}

	var eventinfo []EventsInfo

	for i := 0; i < len(resp); i++ {
		eventinfo[i].EventUID = resp[i].EventUID
		eventinfo[i].EventDate = resp[i].EventDate
		eventinfo[i].EventName = resp[i].Name
		eventinfo[i].EventType = resp[i].EventType

	}
	c.JSON(http.StatusOK, eventinfo)
}

func (s Server) GetEventsByDate(c *gin.Context) {
	var eventdate EventDate
	if err := c.ShouldBind(&eventdate); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	resp, err := s.Client.GetEventsByDate(eventdate.EventDate)
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
		s.Logger.Error(errIncorrectField)
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}
	// convert phone from string to int64
	phoneInt64, err := strconv.ParseInt(phone, 10, 64)
	if err != nil {
		s.Logger.Error(errIncorrectPhoneNo)
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	userInfo := lsdb.UserInfo{
		Name:  name,
		Phone: phoneInt64,
		GovID: govID,
		EMail: eMail,
	}

	err1 := s.Client.AddNewUserInfo(userInfo)
	if err1 != nil {
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
		s.Logger.Error("internal server error")
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		return
	}
	c.JSON(http.StatusOK, "User Updated Successfully")
}

func (s Server) GetUserInfoByID(c *gin.Context) {
	userID, exists1 := c.GetQuery("userid")
	if !exists1 {
		s.Logger.Error(errIncorrectField)
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	userIDConv, _ := primitive.ObjectIDFromHex(userID)

	resp, err := s.Client.GetUserInfoByID(userIDConv)
	if err != nil {
		s.Logger.Error("internal server error")
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		return
	}

	if (lsdb.UserInfo{} == *resp) {
		s.Logger.Error("empty response")
		c.JSON(http.StatusNotFound, "Not found")
		return
	}

	userInfoStruct := UserInfoStruct{
		UID:   resp.UID,
		Name:  resp.Name,
		Phone: resp.Phone,
		EMail: resp.EMail,
	}

	c.JSON(http.StatusOK, userInfoStruct)
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
