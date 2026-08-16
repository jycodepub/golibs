// Package mongo provides a MongoDB client & MongoDB admin utilities
package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (c *Client) Insert(ctx context.Context, database string, collection string, document any) (string, error) {
	result, err := c.GetCollection(database, collection).InsertOne(ctx, document)
	if err != nil {
		log.Printf("Failed to insert document: %v", err)
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (c *Client) InsertMany(ctx context.Context, database string, collection string, documents []any) (int, error) {
	rst, err := c.GetCollection(database, collection).InsertMany(ctx, documents)
	if err != nil {
		log.Printf("Failed to insert documents: %v", err)
		return 0, err
	}
	return len(rst.InsertedIDs), nil
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

func (c *Client) Update(ctx context.Context, database string, collection string, filter any, update any) error {
	_, err := c.GetCollection(database, collection).UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		log.Printf("Failed to upsert document: %v", err)
	}
	return err
}

func (c *Client) DeleteById(ctx context.Context, database string, collection string, id string) (int64, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Failed to convert id to object id: %v", err)
		return 0, err
	}
	rst, err := c.GetCollection(database, collection).DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		log.Printf("Failed to delete document: %v", err)
		return 0, err
	}
	return rst.DeletedCount, nil
}

func (c *Client) DeleteMany(ctx context.Context, database string, collection string, filter any) (int64, error) {
	rst, err := c.GetCollection(database, collection).DeleteMany(ctx, filter)
	if err != nil {
		log.Printf("Failed to delete many documents: %v", err)
		return 0, err
	}
	return rst.DeletedCount, nil
}