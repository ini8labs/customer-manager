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

type WinnerInfo struct {
	EventID   primitive.ObjectID `bson:"event_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id,omitempty"`
	WinType   string             `bson:"win_type,omitempty"`
	AmountWon int                `bson:"amount_won,omitempty"`
	CreatedAt primitive.DateTime `bson:"created_at,omitempty"`
}

const winnersInfoCollection = "winners_info"

func (c *Client) AddNewWinner(winnerInfo WinnerInfo) error {
	t := time.Now()
	currTime := primitive.NewDateTimeFromTime(t)

	winnerInfo.CreatedAt = currTime
	_, err := c.Collection(winnersInfoCollection).InsertOne(context.TODO(), winnerInfo)
	if err != nil {
		return fmt.Errorf("error adding winner info: %s", err.Error())
	}

	return nil
}

func (c *Client) GetEventWinners(eventID primitive.ObjectID) ([]WinnerInfo, error) {
	filter := bson.M{"event_id": eventID}

	result, err := c.Collection(winnersInfoCollection).Find(context.TODO(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching the winners for the event: %s", err.Error())
	}

	var winnersInfo []WinnerInfo

	if err := result.All(context.TODO(), &winnersInfo); err != nil {
		return nil, fmt.Errorf("error decoding the winners info in winnerInfo slice: %s", err.Error())
	}

	return winnersInfo, nil
}
