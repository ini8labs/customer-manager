package apis

import (
	// "errors"
	// "strings"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ini8labs/lsdb"
	// "github.com/sirupsen/logrus"
)

type UserInformation struct {
	Name  string `bson:"name,omitempty"`
	Phone int64  `bson:"phone,omitempty"`
	EMail string `bson:"e_mail,omitempty"`
}

type NewUserInfoFormat struct {
	GovID string `bson:"gov_id,omitempty"`
	UserInformation
}

type UpdateInfoStruct struct {
	UID   string `bson:"_id,omitempty"`
	Key   string `bson:"key,omitempty"`
	Value string `bson:"value,omitempty"`
}

var respUserInfo UserInformation
var userInfo lsdb.UserInfo

func (s Server) newUserInfo(c *gin.Context) {
	var newUserInfoFormat NewUserInfoFormat
	if err := c.ShouldBind(&newUserInfoFormat); err != nil {
		s.Logger.Error("bad format")
		c.JSON(http.StatusBadRequest, "bad Format")
		return
	}

	userInfo.Name = newUserInfoFormat.Name
	userInfo.Phone = newUserInfoFormat.Phone
	if err := validatePhoneNumberInt(userInfo.Phone); err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userInfo.GovID = newUserInfoFormat.GovID
	userInfo.EMail = newUserInfoFormat.EMail

	err := s.Client.AddNewUserInfo(userInfo)
	if err != nil {
		s.Logger.Error("internal server error")
		c.JSON(http.StatusBadRequest, "something went wrong with the server")
		return
	}

	c.JSON(http.StatusOK, "User Info Added Successfully")

}

// not running as expecetd
func (s Server) updateUserInfo(c *gin.Context) {
	var userInfo UpdateInfoStruct
	if err := c.ShouldBind(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	userIDConv, err := validateID(userInfo.UID)
	if err != nil {
		s.Logger.Error(errInvalidUserID)
		c.JSON(http.StatusBadRequest, errInvalidUserID.Error())
	}

	if userInfo.Key == "phone" {
		if err := validatePhoneNumberString(userInfo.Value); err != nil {
			s.Logger.Error(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}

	err1 := s.Client.UpdateUserInfo(userIDConv, userInfo.Key, userInfo.Value)
	if err1 != nil {
		s.Logger.Error("internal server error")
		c.JSON(http.StatusBadRequest, "something went wrong with the server")
		return
	}
	c.JSON(http.StatusCreated, "User Updated Successfully")
}

func (s Server) getUserInfoByID(c *gin.Context) {
	userID := c.Param("id")

	userIDConv, err := validateID(userID)
	if err != nil {
		s.Logger.Errorf("error converting string to HEX: %s", err.Error())
		c.JSON(http.StatusBadRequest, "invalid User ID")
		return
	}

	resp, err := s.Client.GetUserInfoByID(userIDConv)
	if err != nil {
		s.Logger.Error("internal server error")
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		return
	}

	respUserInfo = requiredUserInfo(resp)
	c.JSON(http.StatusOK, respUserInfo)

	fmt.Println(respUserInfo)
}

// not running as expecetd
func (s Server) deleteUserInfoByID(c *gin.Context) {
	userInfo.GovID = c.Param("id")

	err := s.Client.DeleteUserInfo(userInfo)
	if err != nil {
		s.Logger.Error("internal server error")
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		return
	}
	c.JSON(http.StatusOK, "User info Deleted successfully")
}
