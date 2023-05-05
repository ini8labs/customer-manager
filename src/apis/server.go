package apis

import (
	docs "customer-manager/docs"
	"errors"

	//"net/http"
	//"strconv"
	//"strings"

	"github.com/gin-gonic/gin"
	"github.com/ini8labs/lsdb"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/sirupsen/logrus"
)

var (
	errIncorrectField     error = errors.New("incorrect field")
	errMinNumbers         error = errors.New("minimum  number that can be selected is 1")
	errMaxNumbers         error = errors.New("maximum numbers that can be selected are  5")
	errDuplicatedNumbers  error = errors.New("duplicate numbers not allowed")
	errMinAmount          error = errors.New("minimum amount that can be placed is 1")
	errInvalidAmount      error = errors.New("invalid amount")
	errNumberNotAllowed   error = errors.New("bet numbers should be between 1 and 90")
	errIncorrectPhoneNo   error = errors.New("phone number is entered incorrectly")
	errInvalidEventID     error = errors.New("event ID incorrect")
	errInvalidUserID      error = errors.New("invalid user ID")
	errInvalidBetUID      error = errors.New("Bet UID is incorrect")
	errInvalidID          error = errors.New("id is incorrect")
	errInvalidPhoneNum    error = errors.New("invalid phone number")
	errIncorrectEventType error = errors.New("invalid event type")
	errInvalidDateFormat  error = errors.New("invalid date format")
	errInvalidDate        error = errors.New("invalid date")
	errInvaildGovID       error = errors.New("invalid government ID")
	errNoRecords          error = errors.New("no bets placed in this event type")
	errEmptyResp          error = errors.New("no records found")
	errInvalidName        error = errors.New("not a valid name")
	errInvalidKey         error = errors.New("this key is not allowed")
)

type Server struct {
	*logrus.Logger
	*lsdb.Client
	Addr string
}

func NewServer(server Server) error {

	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"

	// API end points
	r.POST("/api/v1/bet", server.placeBets)
	r.PUT("/api/v1/bet", server.updateBets)
	r.GET("/api/v1/bet/:userid", server.userBets)
	r.DELETE("/api/v1/bet/:betuid", server.deleteBets)
	r.GET("/api/v1/bet/history", server.betsHistoryByEventType)
	r.GET("/api/v1/user/:id", server.getUserInfoByID)
	r.POST("/api/v1/user", server.newUserInfo)
	r.PUT("/api/v1/user", server.updateUserInfo)
	r.DELETE("/api/v1/user/:id", server.deleteUserInfoByID)
	// r.GET("/api/v1/event", server.getAllEvents)
	// r.GET("/api/v1/event", server.getEventInfoByDate)
	r.GET("/api/v1/event", server.getEvents)
	r.GET("/api/v1/event/available", server.eventsAvailableForBets)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r.Run(server.Addr)
}
