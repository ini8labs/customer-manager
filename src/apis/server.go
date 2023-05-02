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
	errIncorrectField    error = errors.New("incorrect field")
	errMinNumbers        error = errors.New("minimum  number that can be selected is 1")
	errMaxNumbers        error = errors.New("maximum numbers that can be selected are  5")
	errDuplicatedNumbers error = errors.New("duplicate numbers not allowed")
	errMinAmount         error = errors.New("minimum amount that can be placed is 1")
	errInvalidAmount     error = errors.New("invalid amount")
	errNumberNotAllowed  error = errors.New("bet numbers should be between 1 and 90")
	errIncorrectPhoneNo  error = errors.New("phone number is entered incorrectly")
	errInvalidEventID    error = errors.New("event ID incorrect")
	errInvalidUserID     error = errors.New("invalid user ID")
	errInvalidBetUID     error = errors.New("Bet UID is incorrect")
	errInvalidID         error = errors.New("ID is incorrect")
	errInvalidPhoneNum   error = errors.New("invalid phone number")
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
	r.GET("/api/v1/bet/:id", server.userBets)
	r.DELETE("/api/v1/bet/:id", server.deleteBets)
	r.GET("/api/v1/bet/history/:eventuid", server.betsHistorybyEvent)
	r.GET("/api/v1/user/:id", server.getUserInfoByID)
	r.POST("/api/v1/user", server.newUserInfo)
	r.PUT("/api/v1/user", server.updateUserInfo) // undesired behaviour
	r.DELETE("/api/v1/user/:id", server.deleteUserInfoByID)
	r.GET("/api/v1/event", server.getAllEvents)
	// r.POST("/api/v1/event/eventdata_bydate", GetEventsByDate

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r.Run(server.Addr)
}
