package lsdb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type LotteryEventInfo struct {
	EventUID      primitive.ObjectID `bson:"_id,omitempty"`
	EventDate     primitive.DateTime `bson:"event_date,omitempty"`
	Name          string             `bson:"name,omitempty"`
	EventType     string             `bson:"event_type,omitempty"`
	WinningNumber int                `bson:"winning_number,omitempty"`
	CreatedAt     primitive.DateTime `bson:"created_at,omitempty"`
	UpdatedAt     primitive.DateTime `bson:"updated_at,omitempty"`
}

const lotteryEventsInfoCollection = "lottery_events_info"

func (c *Client) AddNewEvent(event LotteryEventInfo) error {
	t := time.Now()
	event.CreatedAt = primitive.NewDateTimeFromTime(t)
	event.UpdatedAt = primitive.NewDateTimeFromTime(t)
	_, err := c.Collection(lotteryEventsInfoCollection).InsertOne(context.TODO(), &event)
	if err != nil {
		return fmt.Errorf("error Inserting new event info: %s", err.Error())
	}

	return nil
}

func (c *Client) UpdateEvent(event LotteryEventInfo) error {
	return nil
}

func (c *Client) DeleteEvent(eventID primitive.ObjectID) error {
	filter := bson.M{"_id": eventID}

	_, err := c.Collection(lotteryEventsInfoCollection).DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error deleting the event : %s", err.Error())
	}

	return nil
}

func (c *Client) GetEventsByType(eventType string) ([]LotteryEventInfo, error) {
	filter := bson.M{"event_type": eventType}

	results, err := c.Collection(lotteryEventsInfoCollection).Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("error finding the events: %s", err.Error())
	}

	return decodeResultInSlice(results)
}

func (c *Client) GetEventsByDate(date primitive.DateTime) ([]LotteryEventInfo, error) {
	filter := bson.M{"event_date": date}

	return c.getEventsByFilter(filter)
}

func (c *Client) GetEventByDateRange(startDate, endDate primitive.DateTime) ([]LotteryEventInfo, error) {
	filter := bson.M{
		"event_date": bson.M{
			"$lte": endDate,
			"$gte": startDate,
		},
	}

	return c.getEventsByFilter(filter)
}

func (c *Client) GetAllEvents() ([]LotteryEventInfo, error) {
	filter := bson.M{}
	return c.getEventsByFilter(filter)
}

func (c *Client) getEventsByFilter(filter interface{}) ([]LotteryEventInfo, error) {
	result, err := c.Collection(lotteryEventsInfoCollection).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	return decodeResultInSlice(result)
}
func decodeResultInSlice(results *mongo.Cursor) ([]LotteryEventInfo, error) {
	var lotteryEvents []LotteryEventInfo
	if err := results.All(context.TODO(), &lotteryEvents); err != nil {
		return nil, fmt.Errorf("error decoding the result in LotteryEventInfo slice: %s", err.Error())
	}

	return lotteryEvents, nil
}
