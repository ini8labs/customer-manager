package apis

import (
	// "errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	// "github.com/ini8labs/lsdb"
)

// type EventsInfo struct {
// 	EventUID  string `bson:"_id,omitempty"`
// 	EventDate Date   `bson:"event_date,omitempty"`
// 	EventName string `bson:"name,omitempty"`
// 	EventType string `bson:"event_type,omitempty"`
// }

type EventsInfo struct {
	EventUID  string `bson:"_id,omitempty"`
	EventDate Date   `bson:"event_date,omitempty"`
	EventName string `bson:"name,omitempty"`
	EventType string `bson:"event_type,omitempty"`
}

type Date struct {
	Day   int `bson:"day,omitempty"`
	Month int `bson:"month,omitempty"`
	Year  int `bson:"year,omitempty"`
}

var eventInfo []EventsInfo

func (s Server) getAllEvents(c *gin.Context) {

	resp, err := s.Client.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		s.Logger.Errorln(err)
		return
	}

	fmt.Println(resp[0].EventDate)
	eventInfo = requiredEventInfo(resp)
	fmt.Println("==================")
	fmt.Println(eventInfo[0].EventDate)
	fmt.Println("==================")
	c.JSON(http.StatusOK, eventInfo)
}

func (s Server) getEventInfoByDate(c *gin.Context) {
	date := c.Query("date")

	eventDateInfo := stringToDateStruct(date)

	resp, err := s.Client.GetEventsByDate(convertTimeToPrimitive(eventDateInfo))
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, errInvalidDate.Error())
		return
	}

	result := initializeEventInfo(resp)
	c.JSON(http.StatusOK, result)
}

func (s Server) eventsAvailableToday(c *gin.Context) []EventsInfo {
	currentDate := time.Now().Format("2006-01-02")

	eventDateInfo := stringToDateStruct(currentDate)

	resp, err := s.Client.GetEventsByDate(convertTimeToPrimitive(eventDateInfo))
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, errInvalidDate.Error())
		return nil
	}

	result := initializeEventInfo(resp)
	return result
}
