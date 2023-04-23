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

type EventParticipantInfo struct {
	BetUID          primitive.ObjectID `bson:"_id,omitempty"`
	EventUID        primitive.ObjectID `bson:"event_id,omitempty"`
	ParticipantInfo `bson:"participant_info,omitempty"`
	CreatedAt       primitive.DateTime `bson:"created_at,omitempty"`
	UpdatedAt       primitive.DateTime `bson:"updated_at,omitempty"`
}

type ParticipantInfo struct {
	UserID     primitive.ObjectID `bson:"user_id,omitempty"`
	BetNumbers []int              `bson:"bet_numbers,omitempty"`
	Amount     int                `bson:"amount,omitempty"`
}

const eventParticipantInfoCollection = "event_participants_info"

func (c *Client) AddUserBet(participantInfo EventParticipantInfo) error {
	t := time.Now()
	currTime := primitive.NewDateTimeFromTime(t)

	participantInfo.CreatedAt = currTime
	participantInfo.UpdatedAt = currTime

	_, err := c.Collection(eventParticipantInfoCollection).InsertOne(context.TODO(), participantInfo)
	if err != nil {
		return fmt.Errorf("error Inserting the new participant info: %s", err.Error())
	}

	return nil
}

func (c *Client) UpdateUserBet(updatedInfo EventParticipantInfo) error {
	t := time.Now()
	currTime := primitive.NewDateTimeFromTime(t)

	update := bson.M{"$set": bson.M{"amount": updatedInfo.Amount, "updated_at": currTime, "bet_numbers": updatedInfo.BetNumbers}}
	_, err := c.Collection(eventParticipantInfoCollection).UpdateByID(context.TODO(), updatedInfo.BetUID, update)
	if err != nil {
		return fmt.Errorf("error updating the user bets: %s", err.Error())
	}

	return nil
}

func (c *Client) GetUserBets(userID primitive.ObjectID) ([]EventParticipantInfo, error) {
	filter := bson.M{"participant_info.user_id": userID}
	result, err := c.Collection(eventParticipantInfoCollection).Find(context.TODO(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding the user bet details: %s", err.Error())
	}

	var eventInfo []EventParticipantInfo

	if err := result.All(context.TODO(), &eventInfo); err != nil {
		return nil, fmt.Errorf("error decoding the data in participantEventInfo slice: %s", err.Error())
	}

	return eventInfo, nil
}

func (c *Client) DeleteUserBet(betUID primitive.ObjectID) error {
	filter := bson.M{"_id": betUID}

	_, err := c.Collection(eventParticipantInfoCollection).DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error deleting the user bet: %s", err.Error())
	}

	return nil
}

func (c *Client) GetParticipantsInfoByEventID(eventID primitive.ObjectID) ([]EventParticipantInfo, error) {
	filter := bson.M{"event_id": eventID}

	result, err := c.Collection(eventParticipantInfoCollection).Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("error fetching the participants info for the bet: %s", err.Error())
	}

	var eventInfo []EventParticipantInfo

	if err := result.All(context.TODO(), &eventInfo); err != nil {
		return nil, fmt.Errorf("error decoding the data in participantEventInfo slice: %s", err.Error())
	}

	return eventInfo, nil
}
