package apis

import (
	"errors"
	"fmt"
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
	r.GET("/api/v1/eventdata", GetAllEvents)
	//r.POST("/api/v1/userinfo_ID", GetUserInfoByID)
	r.DELETE("/api/v1/delete", DeleteUserInfoByID)
	r.POST("/api/v1/eventdata_bydate", GetEventsByDate)
	r.POST("/api/v1/addnew", PlaceBets)

	r.POST("/api/v1/userinfo_ID", GetUserInfoByID)
	return r.Run(addr)
}

// function for validation
// func validateBetPlaceInput(requestData Bet) error {
// 	if len(requestData.Numbers) < 1 {
// 		return errMinAmount
// 	}
// 	if requestData.Amount <= 0 {
// 		return fmt.Errorf("amount can not be 0 or negative")
// 	}
// 	return nil
// }

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
	fmt.Println("User Added Successfully")
}

// ----Events-----
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

// need to delete user info from userinfo as well as user stack
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
