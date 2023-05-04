package apis

import (
	// "errors"

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

	eventInfo = requiredEventInfo(resp)
	c.JSON(http.StatusOK, eventInfo)
}

func (s Server) getEventInfoByDate(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		s.Logger.Error(errInvalidDate)
		c.JSON(http.StatusBadRequest, errInvalidDate.Error())
		return
	}

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

func (s Server) getEventInfoByDateRange(startDate, endDate string) ([]EventsInfo, error) {
	startRangeDate := stringToDateStruct(startDate)

	endRangeDate := stringToDateStruct(endDate)

	resp, err := s.Client.GetEventByDateRange(convertTimeToPrimitive(startRangeDate), convertTimeToPrimitive(endRangeDate))
	if err != nil {
		return []EventsInfo{}, err
	}

	result := initializeEventInfo(resp)
	return result, nil
}

func (s Server) eventsAvailable() ([]EventsInfo, error) {

	startDate := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	endDate := time.Now().AddDate(3, 0, 0).Format("2006-01-02")

	eventStartDateInfo := stringToDateStruct(startDate)
	eventEndDateInfo := stringToDateStruct(endDate)

	resp, err := s.Client.GetEventByDateRange(convertTimeToPrimitive(eventStartDateInfo), convertTimeToPrimitive(eventEndDateInfo))
	if err != nil {
		return []EventsInfo{}, errInvalidDate
	}

	result := initializeEventInfo(resp)
	return result, nil
}

func (s Server) eventsAvailableForBets(c *gin.Context) {

	startDate := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	endDate := time.Now().AddDate(1, 0, 0).Format("2006-01-02")

	eventStartDateInfo := stringToDateStruct(startDate)
	eventEndDateInfo := stringToDateStruct(endDate)

	resp, err := s.Client.GetEventByDateRange(convertTimeToPrimitive(eventStartDateInfo), convertTimeToPrimitive(eventEndDateInfo))
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	result := initializeEventInfo(resp)
	c.JSON(http.StatusOK, result)
}
