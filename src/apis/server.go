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
	errIncorrectField   error = errors.New("incorrect field or fields")
	errMinNumbers       error = errors.New("minimum  number that can be selected is 1")
	errMaxNumbers       error = errors.New("maximum numbers that can be selected are  5")
	errMinAmount        error = errors.New("minimum amount that can be placed is 1")
	errNumberNotAllowed error = errors.New("bet numbers should be between 1 and 90")
	errIncorrectPhoneNo error = errors.New("phone number is entered incorrectly")
	errInvalidEventID   error = errors.New("event ID incorrect")
	errInvalidUserID    error = errors.New("User ID is incorrect")
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
	r.POST("/api/v1/bet/new_bet", server.PlaceBets)
	r.PUT("/api/v1/bet/update", server.UpdateBets)
	r.GET("/api/v1/bet/user_bets", server.UserBets)
	r.DELETE("/api/v1/bet/delete", server.DeleteBets)
	r.GET("/api/v1/bet/history", server.EventHistory)
	r.GET("/api/v1/user_info/userinfo_ID", server.GetUserInfoByID)
	r.POST("/api/v1/user_info/new_user", server.NewUserInfo)
	r.PUT("/api/v1/user_info/update_info", server.UpdateUserInfo)
	r.DELETE("/api/v1/user_info/delete", server.DeleteUserInfoByID)
	r.GET("/api/v1/event/eventdata", server.GetAllEvents)
	// r.POST("/api/v1/event/eventdata_bydate", GetEventsByDate

	return r.Run(server.Addr)
}
