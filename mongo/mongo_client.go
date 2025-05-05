package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Client struct {
	client *mongo.Client
}

func NewClient(connectString string) *Client {
	c, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectString))
	if err != nil {
		panic(err)
	}
	return &Client{client: c}
}

type DataNotFound struct {
	message string
}

func (e *DataNotFound) Error() string {
	return e.message
}

func (c *Client) Close() {
	err := c.client.Disconnect(context.TODO())
	if err != nil {
		log.Printf("Failed to disconnect client: %v", err)
	}
	log.Printf("Disconnected mongodb client")
}

func (c *Client) Insert(ctx context.Context, database string, collection string, document interface{}) error {
	_, err := c.GetCollection(database, collection).InsertOne(ctx, document)
	if err != nil {
		log.Printf("Failed to insert document: %v", err)
	}
	return err
}

func (c *Client) Query(ctx context.Context, database string, collection string, filter interface{}) (*mongo.Cursor, error) {
	return c.GetCollection(database, collection).Find(ctx, filter)
}

func (c *Client) GetCollection(database string, collection string) *mongo.Collection {
	return c.client.Database(database).Collection(collection)
}

func (c *Client) QueryForStruct(ctx context.Context, database string, collection string, filter interface{}, result interface{}) error {
	cur, err := c.Query(ctx, database, collection, filter)
	if err != nil {
		return err
	}
	if cur.Next(ctx) {
		err := cur.Decode(result)
		if err != nil {
			return err
		}
	} else {
		return &DataNotFound{
			message: "Data not found",
		}
	}
	return nil
}
