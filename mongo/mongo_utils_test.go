package mongo

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestMongoEnv(t *testing.T, dbName string) (string, func()) {
	t.Helper()
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatalf("failed to connect to mongo: %v", err)
	}

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = client.Database(dbName).Drop(ctx)
		_ = client.Disconnect(ctx)
	}

	return uri, cleanup
}

func seedCollection(t *testing.T, uri string, dbName string, collName string, docs []any) {
	t.Helper()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer func() { _ = client.Disconnect(context.TODO()) }()

	coll := client.Database(dbName).Collection(collName)
	_, err = coll.InsertMany(context.TODO(), docs)
	if err != nil {
		t.Fatalf("failed to seed collection %s: %v", collName, err)
	}
}

func countCollectionDocs(t *testing.T, uri string, dbName string, collName string) int64 {
	t.Helper()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer func() { _ = client.Disconnect(context.TODO()) }()

	count, err := client.Database(dbName).Collection(collName).CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		t.Fatalf("failed to count docs in %s: %v", collName, err)
	}
	return count
}

func TestMongoUtils_Collections(t *testing.T) {
	dbName := "test_utils_db"
	uri, cleanup := setupTestMongoEnv(t, dbName)
	defer cleanup()

	// 1. Seed initial data
	seedCollection(t, uri, dbName, "users", []any{
		bson.M{"name": "Alice", "role": "admin"},
		bson.M{"name": "Bob", "role": "user"},
	})
	seedCollection(t, uri, dbName, "orders", []any{
		bson.M{"item": "book", "price": 10},
	})

	t.Run("ListCollections", func(t *testing.T) {
		cols := ListCollections(uri, dbName)
		sort.Strings(cols)
		if len(cols) != 2 {
			t.Fatalf("expected 2 collections, got %d: %v", len(cols), cols)
		}
		if cols[0] != "orders" || cols[1] != "users" {
			t.Errorf("unexpected collections listed: %v", cols)
		}
	})

	t.Run("CleanCollection", func(t *testing.T) {
		CleanCollection(uri, dbName, "orders")
		if count := countCollectionDocs(t, uri, dbName, "orders"); count != 0 {
			t.Errorf("expected 0 docs in orders after CleanCollection, got %d", count)
		}
		// users should still have 2 docs
		if count := countCollectionDocs(t, uri, dbName, "users"); count != 2 {
			t.Errorf("expected 2 docs in users, got %d", count)
		}
	})

	t.Run("CleanDB", func(t *testing.T) {
		CleanDB(uri, dbName)
		if count := countCollectionDocs(t, uri, dbName, "users"); count != 0 {
			t.Errorf("expected 0 docs in users after CleanDB, got %d", count)
		}
	})

	t.Run("DropCollection", func(t *testing.T) {
		seedCollection(t, uri, dbName, "logs", []any{bson.M{"msg": "test"}})
		DropCollection(uri, dbName, "logs")

		cols := ListCollections(uri, dbName)
		for _, c := range cols {
			if c == "logs" {
				t.Errorf("expected collection 'logs' to be dropped, but it still exists")
			}
		}
	})

	t.Run("DropCollections", func(t *testing.T) {
		seedCollection(t, uri, dbName, "col1", []any{bson.M{"a": 1}})
		seedCollection(t, uri, dbName, "col2", []any{bson.M{"b": 2}})

		DropCollections(uri, dbName)
		cols := ListCollections(uri, dbName)
		if len(cols) != 0 {
			t.Errorf("expected 0 collections after DropCollections, got %d (%v)", len(cols), cols)
		}
	})
}

func TestMongoUtils_ExportImportDumpRestore(t *testing.T) {
	dbName := "test_utils_dump_db"
	uri, cleanup := setupTestMongoEnv(t, dbName)
	defer cleanup()

	tempDir, err := os.MkdirTemp("", "mongo_utils_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	seedCollection(t, uri, dbName, "products", []any{
		bson.M{"name": "Laptop", "price": 1000},
		bson.M{"name": "Phone", "price": 500},
	})
	seedCollection(t, uri, dbName, "categories", []any{
		bson.M{"name": "Electronics"},
	})

	t.Run("ExportCollection and ImportCollection", func(t *testing.T) {
		ExportCollection(uri, dbName, "products", tempDir)

		dumpFilePath := filepath.Join(tempDir, dbName+".products."+DumpFileExt)
		if _, err := os.Stat(dumpFilePath); os.IsNotExist(err) {
			t.Fatalf("expected export file %s to exist", dumpFilePath)
		}

		// Clean products collection
		CleanCollection(uri, dbName, "products")
		if count := countCollectionDocs(t, uri, dbName, "products"); count != 0 {
			t.Fatalf("expected 0 docs after clean, got %d", count)
		}

		// Import collection back
		ImportCollection(uri, dbName, "products", dumpFilePath)
		if count := countCollectionDocs(t, uri, dbName, "products"); count != 2 {
			t.Errorf("expected 2 docs after import, got %d", count)
		}
	})

	t.Run("DumpDB and RestoreDB", func(t *testing.T) {
		dumpDir := filepath.Join(tempDir, "db_dump")
		if err := os.MkdirAll(dumpDir, 0755); err != nil {
			t.Fatalf("failed to create dump dir: %v", err)
		}

		DumpDB(uri, dbName, dumpDir)

		// Verify files exist in dumpDir
		files, err := os.ReadDir(dumpDir)
		if err != nil {
			t.Fatalf("failed to read dumpDir: %v", err)
		}
		if len(files) < 2 {
			t.Fatalf("expected at least 2 dump files, got %d", len(files))
		}

		// Clean entire DB
		CleanDB(uri, dbName)
		if count := countCollectionDocs(t, uri, dbName, "products"); count != 0 {
			t.Fatalf("expected 0 products, got %d", count)
		}
		if count := countCollectionDocs(t, uri, dbName, "categories"); count != 0 {
			t.Fatalf("expected 0 categories, got %d", count)
		}

		// Restore entire DB
		RestoreDB(uri, dbName, dumpDir)
		if count := countCollectionDocs(t, uri, dbName, "products"); count != 2 {
			t.Errorf("expected 2 products after restore, got %d", count)
		}
		if count := countCollectionDocs(t, uri, dbName, "categories"); count != 1 {
			t.Errorf("expected 1 category after restore, got %d", count)
		}
	})
}
