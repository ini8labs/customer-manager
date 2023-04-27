package apis

import (
	"errors"
	//"net/http"
	//"strconv"
	//"strings"

	"github.com/gin-gonic/gin"
	"github.com/ini8labs/lsdb"

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
	errInvalidUserID     error = errors.New("user ID is incorrect")
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

type UserInfoStruct struct {
	Name  string `bson:"name,omitempty"`
	Phone int64  `bson:"phone,omitempty"`
	EMail string `bson:"e_mail,omitempty"`
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
	r.POST("/api/v1/bet/new", server.PlaceBets)
	r.PUT("/api/v1/bet/update", server.UpdateBets)
	r.GET("/api/v1/bet/user/bets", server.UserBets)
	r.DELETE("/api/v1/bet/delete", server.DeleteBets)
	r.GET("/api/v1/bet/history", server.EventHistory)
	r.GET("/api/v1/user/info/userinfo/ID", server.GetUserInfoByID)
	r.POST("/api/v1/user/info/new", server.NewUserInfo)
	r.PUT("/api/v1/user/info/update", server.UpdateUserInfo) // undesired behaviour
	r.DELETE("/api/v1/user/info/delete", server.DeleteUserInfoByID)
	r.GET("/api/v1/event/eventdata", server.GetAllEvents)
	// r.POST("/api/v1/event/eventdata_bydate", GetEventsByDate

	return r.Run(server.Addr)
}
