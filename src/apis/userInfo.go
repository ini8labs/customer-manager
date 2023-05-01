package apis

import (
	// "errors"
	// "strings"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ini8labs/lsdb"
	// "github.com/sirupsen/logrus"
)

type UserInformation struct {
	Name  string `bson:"name,omitempty"`
	Phone int64  `bson:"phone,omitempty"`
	EMail string `bson:"email,omitempty"`
}

var respUserInfo UserInformation
var userInfo lsdb.UserInfo

// working as expecetd )
func (s Server) newUserInfo(c *gin.Context) {
	name, exists1 := c.GetQuery("name")
	phone, exists2 := c.GetQuery("phone")
	govID, exists3 := c.GetQuery("govid")
	eMail, exists4 := c.GetQuery("email")

	// check for no missing fields
	if !exists1 || !exists2 || !exists3 || !exists4 {
		s.Logger.Error(errIncorrectField)
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}
	// convert phone from string to int64
	phoneInt64, err := strconv.ParseInt(phone, 10, 64)
	if err != nil {
		s.Logger.Error(errIncorrectPhoneNo)
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	userInfo := lsdb.UserInfo{
		Name:  name,
		Phone: phoneInt64,
		GovID: govID,
		EMail: eMail,
	}

	err1 := s.Client.AddNewUserInfo(userInfo)
	if err1 != nil {
		s.Logger.Error("Internal server Error")
		c.JSON(http.StatusInternalServerError, "Something is wrong with the server")
		return
	}
	c.JSON(http.StatusOK, "User Info Added Successfully")
	s.Logger.Info("Create operation performed")
}

// not running as expecetd
func (s Server) UpdateUserInfo(c *gin.Context) {
	var userInfo UpdateInfoStruct
	if err := c.BindJSON(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	err1 := s.Client.UpdateUserInfo(userInfo.UserID, userInfo.Key, userInfo.Value)
	if err1 != nil {
		s.Logger.Error("internal server error")
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		return
	}
	c.JSON(http.StatusCreated, "User Updated Successfully")
}

func (s Server) getUserInfoByID(c *gin.Context) {
	userID := c.Param("userid")

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
	userInfo.GovID = c.Param("govid")

	err := s.Client.DeleteUserInfo(userInfo)
	if err != nil {
		s.Logger.Error("internal server error")
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		return
	}
	c.JSON(http.StatusOK, "User info Deleted successfully")
}
