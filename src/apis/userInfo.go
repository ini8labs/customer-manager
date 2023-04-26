package apis

import (
	// "errors"
	// "strings"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ini8labs/lsdb"

	// "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s Server) NewUserInfo(c *gin.Context) {
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

func (s Server) UpdateUserInfo(c *gin.Context) {
	var userInfo UpdateInfoStruct
	if err := c.ShouldBind(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	err1 := s.Client.UpdateUserInfo(userInfo.UserID, userInfo.Key, userInfo.Value)
	if err1 != nil {
		s.Logger.Error("internal server error")
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		return
	}
	c.JSON(http.StatusOK, "User Updated Successfully")
}

func (s Server) GetUserInfoByID(c *gin.Context) {
	userID, exists1 := c.GetQuery("userid")
	if !exists1 {
		s.Logger.Error(errIncorrectField)
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	userIDConv, err := primitive.ObjectIDFromHex(userID)
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

	if (lsdb.UserInfo{} == *resp) {
		s.Logger.Error("empty response")
		c.JSON(http.StatusNotFound, "Not found")
		return
	}

	userInfoStruct := UserInfoStruct{
		UID:   resp.UID,
		Name:  resp.Name,
		Phone: resp.Phone,
		EMail: resp.EMail,
	}

	c.JSON(http.StatusOK, userInfoStruct)
}

func (s Server) DeleteUserInfoByID(c *gin.Context) {
	var userinfo lsdb.UserInfo

	if err := c.ShouldBind(&userinfo); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Format")
		return
	}

	err1 := s.Client.DeleteUserInfo(userinfo)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		return
	}
	c.JSON(http.StatusOK, "User info Deleted successfully")
}
