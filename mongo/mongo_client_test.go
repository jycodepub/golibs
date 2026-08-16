package mongo

import (
	"context"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func getTestMongoURI() string {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	return uri
}

type SampleUser struct {
	Name  string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	Age   int    `bson:"age" json:"age"`
}

func TestMongoClient(t *testing.T) {
	uri := getTestMongoURI()
	client := NewClient(uri)
	if client == nil {
		t.Fatal("expected non-nil Client from NewClient")
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbName := "test_db"
	collName := "test_users"

	// Cleanup test collection before running tests
	_ = client.GetCollection(dbName, collName).Drop(ctx)

	t.Run("GetCollection", func(t *testing.T) {
		coll := client.GetCollection(dbName, collName)
		if coll == nil {
			t.Fatal("expected non-nil Collection")
		}
		if coll.Name() != collName {
			t.Errorf("expected collection name %s, got %s", collName, coll.Name())
		}
	})

	var insertedID string
	t.Run("Insert", func(t *testing.T) {
		user := SampleUser{
			Name:  "Alice",
			Email: "alice@example.com",
			Age:   30,
		}

		id, err := client.Insert(ctx, dbName, collName, user)
		if err != nil {
			t.Fatalf("failed to insert document: %v", err)
		}
		if id == "" {
			t.Fatal("expected non-empty hex string ObjectID")
		}
		insertedID = id
	})

	t.Run("QueryForStruct", func(t *testing.T) {
		var result SampleUser
		err := client.QueryForStruct(ctx, dbName, collName, bson.M{"email": "alice@example.com"}, &result)
		if err != nil {
			t.Fatalf("QueryForStruct failed: %v", err)
		}
		if result.Name != "Alice" || result.Age != 30 {
			t.Errorf("unexpected query result: %+v", result)
		}
	})

	t.Run("QueryForStruct_NotFound", func(t *testing.T) {
		var result SampleUser
		err := client.QueryForStruct(ctx, dbName, collName, bson.M{"email": "nonexistent@example.com"}, &result)
		if err == nil {
			t.Fatal("expected DataNotFound error, got nil")
		}
		dnf, ok := err.(*DataNotFound)
		if !ok {
			t.Fatalf("expected *DataNotFound error type, got %T: %v", err, err)
		}
		if dnf.Error() != "Data not found" {
			t.Errorf("unexpected error message: %s", dnf.Error())
		}
	})

	t.Run("Query", func(t *testing.T) {
		cursor, err := client.Query(ctx, dbName, collName, bson.M{"name": "Alice"})
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}
		defer cursor.Close(ctx)

		var users []SampleUser
		if err := cursor.All(ctx, &users); err != nil {
			t.Fatalf("cursor.All failed: %v", err)
		}
		if len(users) != 1 || users[0].Name != "Alice" {
			t.Errorf("unexpected users from cursor: %+v", users)
		}
	})

	t.Run("InsertMany", func(t *testing.T) {
		docs := []any{
			SampleUser{Name: "Bob", Email: "bob@example.com", Age: 25},
			SampleUser{Name: "Charlie", Email: "charlie@example.com", Age: 35},
		}

		count, err := client.InsertMany(ctx, dbName, collName, docs)
		if err != nil {
			t.Fatalf("InsertMany failed: %v", err)
		}
		if count != 2 {
			t.Errorf("expected 2 inserted documents, got %d", count)
		}
	})

	t.Run("Update", func(t *testing.T) {
		filter := bson.M{"email": "bob@example.com"}
		update := bson.M{"$set": bson.M{"age": 26}}

		err := client.Update(ctx, dbName, collName, filter, update)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		var updated SampleUser
		err = client.QueryForStruct(ctx, dbName, collName, filter, &updated)
		if err != nil {
			t.Fatalf("QueryForStruct after update failed: %v", err)
		}
		if updated.Age != 26 {
			t.Errorf("expected updated age 26, got %d", updated.Age)
		}
	})

	t.Run("DeleteById", func(t *testing.T) {
		deletedCount, err := client.DeleteById(ctx, dbName, collName, insertedID)
		if err != nil {
			t.Fatalf("DeleteById failed: %v", err)
		}
		if deletedCount != 1 {
			t.Errorf("expected 1 deleted document, got %d", deletedCount)
		}
	})

	t.Run("DeleteMany", func(t *testing.T) {
		deletedCount, err := client.DeleteMany(ctx, dbName, collName, bson.M{})
		if err != nil {
			t.Fatalf("DeleteMany failed: %v", err)
		}
		if deletedCount != 2 {
			t.Errorf("expected 2 deleted documents, got %d", deletedCount)
		}
	})
}

func TestDataNotFound(t *testing.T) {
	err := &DataNotFound{message: "custom message"}
	if err.Error() != "custom message" {
		t.Errorf("expected 'custom message', got '%s'", err.Error())
	}
}
