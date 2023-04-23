package lsdb

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create MONGO_DB_CONN_STRING as env variable

const mongoDBConnString = "MONGO_DB_CONN_STRING"
const lotterySystemDataBase = "lottery_system_db"

type Client struct {
	ConnectionURI string
	*mongo.Client
	*mongo.Database
}

func NewClient() (*Client, error) {
	connURI, err := getConnString()
	if err != nil {
		return nil, err
	}

	return &Client{
		ConnectionURI: connURI,
	}, nil

}

func getConnString() (string, error) {
	val, ok := os.LookupEnv(mongoDBConnString)
	if !ok {
		return "", fmt.Errorf("%s env variable is not set", mongoDBConnString)
	}

	return val, nil
}

func (c *Client) OpenConnection() error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(c.ConnectionURI))
	if err != nil {
		return fmt.Errorf("error connecting with Mongo: %s", err.Error())
	}

	c.Client = client
	c.Database = client.Database(lotterySystemDataBase)
	return nil
}

func (c *Client) CloseConnection() error {
	if err := c.Disconnect(context.TODO()); err != nil {
		return fmt.Errorf("error closing the mongo client : %s", err.Error())
	}

	return nil
}
