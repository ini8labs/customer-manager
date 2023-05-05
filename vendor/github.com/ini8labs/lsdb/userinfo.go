package lsdb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserInfo struct {
	UID       primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name,omitempty"`
	Phone     int64              `bson:"phone,omitempty"`
	GovID     string             `bson:"gov_id,omitempty"`
	EMail     string             `bson:"e_mail,omitempty"`
	CreatedAT string             `bson:"created_at,omitempty"`
	UpdatedAT string             `bson:"updated_at,omitempty"`
}

const userInfoCollection = "user_info"

func (c *Client) GetUserInfoByPhone(phoneNumber int64) (*UserInfo, error) {
	phoneFilter := bson.M{"phone": phoneNumber}

	return c.getUserInfoUsingFilter(phoneFilter)
}

func (c *Client) GetUserInfoByID(id primitive.ObjectID) (*UserInfo, error) {
	idFilter := bson.M{"_id": id}

	return c.getUserInfoUsingFilter(idFilter)
}

func (c *Client) GetUserInfoByGovID(govID string) (*UserInfo, error) {
	idFilter := bson.M{"gov_id": govID}

	return c.getUserInfoUsingFilter(idFilter)
}

func (c *Client) getUserInfoUsingFilter(filter bson.M) (*UserInfo, error) {
	userInfo := &UserInfo{}

	if err := c.Collection(userInfoCollection).FindOne(context.TODO(), filter).Decode(userInfo); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return userInfo, nil
}

func (c *Client) AddNewUserInfo(userInfo UserInfo) error {
	t := time.Now()
	createTime := t.Format(time.RFC3339)
	userInfo.CreatedAT = createTime
	userInfo.UpdatedAT = createTime

	_, err := c.Collection(userInfoCollection).InsertOne(context.TODO(), &userInfo)
	if err != nil {
		return fmt.Errorf("error adding new user info: %s", err.Error())
	}

	return nil
}

func (c *Client) DeleteUserInfo(userInfo UserInfo) error {
	govIdFilter := bson.M{"gov_id": userInfo.GovID}
	_, err := c.Collection(userInfoCollection).DeleteOne(context.TODO(), govIdFilter)
	if err != nil {
		return fmt.Errorf("error deleting the user info: %s", err.Error())
	}

	return nil
}

func (c *Client) UpdateUserInfo(uid primitive.ObjectID, key, value string) error {
	update := bson.M{"$set": bson.M{key: value}}
	_, err := c.Collection(userInfoCollection).UpdateByID(context.TODO(), uid, update)
	if err != nil {
		return fmt.Errorf("error updating the User info %s", err.Error())
	}

	return nil
}

func GetAllUserInfo() ([]UserInfo, error) {
	return nil, nil
}
