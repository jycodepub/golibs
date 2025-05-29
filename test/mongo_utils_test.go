package test

import (
	"fmt"
	"github.com/jycodepub/golibs/mongo"
	"testing"
)

const (
	MongoUrl = "mongodb://jysrv01:27017"
	Database = "test"
)

func TestMongoUtils(t *testing.T) {
	collections := mongo.ListCollections(MongoUrl, Database)
	for _, collection := range collections {
		fmt.Println(collection)
	}
	if len(collections) != 2 {
		t.Error("should return 2 collections")
	}

	mongo.DumpDB(MongoUrl, Database, ".")

	mongo.CleanDB(MongoUrl, Database)

	mongo.RestoreDB(MongoUrl, Database, ".")
}
