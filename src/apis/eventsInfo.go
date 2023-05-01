package apis

import (
	// "errors"
	// "strconv"
	// "strings"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/ini8labs/lsdb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventsInfo struct {
	EventUID  primitive.ObjectID `bson:"_id,omitempty"`
	EventDate primitive.DateTime `bson:"event_date,omitempty"`
	EventName string             `bson:"name,omitempty"`
	EventType string             `bson:"event_type,omitempty"`
}

var eventInfo []EventsInfo

func (s Server) getAllEvents(c *gin.Context) {

	resp, err := s.Client.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		s.Logger.Errorln(err)
		return
	}

	eventInfo = requiredEventInfo(resp)
	c.JSON(http.StatusOK, eventInfo)
}

// func (s Server) GetEventsByDate(c *gin.Context) {
// 	var eventdate EventDate
// 	if err := c.ShouldBind(&eventdate); err != nil {
// 		c.JSON(http.StatusBadRequest, "Bad Format")
// 		return
// 	}

// 	resp, err := s.Client.GetEventsByDate(eventdate.EventDate)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
// 		return
// 	}
// 	c.JSON(http.StatusOK, resp)
// }
