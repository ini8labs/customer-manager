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
	"go.mongodb.org/mongo-driver/bson/primitive"
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
)

type Server struct {
	*logrus.Logger
	*lsdb.Client
	Addr string
}

type UpdateInfoStruct struct {
	UserID primitive.ObjectID `bson:"userid,omitempty"`
	Key    string             `bson:"key,omitempty"`
	Value  string             `bson:"value,omitempty"`
}

func NewServer(server Server) error {

	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"

	// API end points
	r.POST("/api/v1/bet/new", server.placeBets)
	r.PUT("/api/v1/bet/update", server.updateBets)
	r.GET("/api/v1/bet/user/:bets", server.userBets)
	r.DELETE("/api/v1/bet/delete/:betuid", server.deleteBets)
	r.GET("/api/v1/bet/history/:eventuid", server.betsHistorybyEvent)
	r.GET("/api/v1/user/info/userinfo/:userid", server.getUserInfoByID)
	r.POST("/api/v1/user/info/new", server.newUserInfo)
	r.PUT("/api/v1/user/info/update", server.UpdateUserInfo) // undesired behaviour
	r.DELETE("/api/v1/user/info/delete/:govid", server.deleteUserInfoByID)
	r.GET("/api/v1/event/eventdata", server.getAllEvents)
	// r.POST("/api/v1/event/eventdata_bydate", GetEventsByDate

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r.Run(server.Addr)
}
