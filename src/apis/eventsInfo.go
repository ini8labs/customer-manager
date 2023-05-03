package apis

import (
	// "errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ini8labs/lsdb"
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

	eventInfo = requiredEventInfo(resp)
	c.JSON(http.StatusOK, eventInfo)
}

func (s Server) getEventInfoByDate(c *gin.Context) {
	date := c.Query("date")

	eventDate := strings.Split(date, "-")
	intYear, _ := strconv.Atoi(eventDate[0])
	intMonth, _ := strconv.Atoi(eventDate[1])
	intDay, _ := strconv.Atoi(eventDate[2])
	eventdate := Date{
		Year:  intYear,
		Month: intMonth,
		Day:   intDay,
	}

	fmt.Println(eventdate)
	var temp = lsdb.LotteryEventInfo{}
	fmt.Println(convertTimeToPrimitive(eventdate))
	temp.EventDate = convertTimeToPrimitive(eventdate)
	fmt.Println(temp.EventDate)

	resp, err := s.Client.GetEventsByDate(temp.EventDate)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, errInvalidDate.Error())
		return
	}

	fmt.Println(resp)
	result := initializeEventInfo(resp)
	c.JSON(http.StatusOK, result)
}

// func (s Server) getEventInfoByDate1(c *gin.Context) {
// 	date := c.Param("date")
// 	eventDate := strings.Split(date, "-")
// 	intYear, _ := strconv.Atoi(eventDate[0])
// 	intMonth, _ := strconv.Atoi(eventDate[1])
// 	intDay, _ := strconv.Atoi(eventDate[2])

// 	eventdate := Date{
// 		Year:  intYear,
// 		Month: intMonth,
// 		Day:   intDay,
// 	}

// }
