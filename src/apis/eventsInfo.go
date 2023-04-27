package apis

import (
	// "errors"
	// "strconv"
	// "strings"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/ini8labs/lsdb"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

func (s Server) GetAllEvents(c *gin.Context) {

	resp, err := s.Client.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		s.Logger.Errorln(err)
		return
	}

	var eventinfo []EventsInfo

	for i := 0; i < len(resp); i++ {
		eventinfo[i].EventUID = resp[i].EventUID
		eventinfo[i].EventDate = resp[i].EventDate
		eventinfo[i].EventName = resp[i].Name
		eventinfo[i].EventType = resp[i].EventType

	}
	c.JSON(http.StatusOK, eventinfo)
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
