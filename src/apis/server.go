package apis

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ini8labs/lsdb"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	errMinAmount error = errors.New("minimum amount is 1")
)

type UserStaking struct {
	UserID     string    `json:"userid"`
	BetNumbers []int     `json:"betnumbers"`
	Amount     int       `json:"amount"`
	EventName  string    `json:"eventname"`
	EventId    string    `json:"eventid"`
	WinNumber  []int     `json:"winnumber"`
	WinType    string    `json:"wintype"`
	Date       time.Time `json:"date"`
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

func NewServer(addr string, log *logrus.Logger) error {

	r := gin.Default()

	// API end points
	r.GET("/api/v1/eventdata", GetAllEvents)
	//r.POST("/api/v1/userinfo_ID", GetUserInfoByID)
	r.DELETE("/api/v1/delete", DeleteUserInfoByID)
	r.POST("/api/v1/eventdata_bydate", GetEventsByDate)
	r.POST("/api/v1/addnew", addSomething)

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

func addSomething(c *gin.Context) {
	dbClient, err := lsdb.NewClient()
	if err != nil {
		panic(err.Error())
	}

	var newUserInfo lsdb.UserInfo
	if err := c.ShouldBind(&newUserInfo); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	if err := dbClient.OpenConnection(); err != nil {
		panic(err.Error())
	}
	defer dbClient.CloseConnection()

	if err := dbClient.AddNewUserInfo(newUserInfo); err != nil {
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
