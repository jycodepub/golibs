package mongo

import (
	"bufio"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const ImportBatchSize = 1000

func ListCollections(connectionUri string, database string) []string {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionUri))
	if err != nil {
		panic(err)
	}
	defer close(client)
	db := client.Database(database)
	return listCollections(db)
}

func CleanDB(connectionUri string, database string) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionUri))
	if err != nil {
		panic(err)
	}
	defer close(client)
	db := client.Database(database)
	names := listCollections(db)
	for _, name := range names {
		count := cleanCollection(db, name)
		log.Printf("Deleted %d document(s) in collection %s", count, name)
	}
}

func DumpDB(connectionUri string, database string, outputDir string) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionUri))
	if err != nil {
		panic(err)
	}
	defer close(client)
	db := client.Database(database)
	names := listCollections(db)
	for _, name := range names {
		saveCollectionToFile(db, name, outputDir)
	}
}

func RestoreDB(connectionUri string, database string, outputDir string) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionUri))
	if err != nil {
		panic(err)
	}
	defer close(client)
	db := client.Database(database)
	files, err := os.ReadDir(outputDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		filename := f.Name()
		collection := strings.Split(filename, ".")[1]
		inputFile := filepath.Join(outputDir, filename)
		count := doImport(db.Collection(collection), inputFile)
		fmt.Printf("Imported %d document(s) to %s\n", count, collection)
	}
}

func ExportCollection(connectionUri string, database string, collectionName string, outputDir string) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionUri))
	if err != nil {
		panic(err)
	}
	defer close(client)
	db := client.Database(database)
	saveCollectionToFile(db, collectionName, outputDir)
}

func ImportCollection(connectionUri string, database string, collectionName string, inputFile string) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionUri))
	if err != nil {
		panic(err)
	}
	defer close(client)
	db := client.Database(database)
	collection := db.Collection(collectionName)
	count := doImport(collection, inputFile)
	log.Printf("Imported %d document(s) in collection: %s\n", count, collectionName)
}

func doImport(collection *mongo.Collection, inputFile string) int {
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	totalCount := 0
	batchCount := 0
	ctx := context.TODO()
	scanner := bufio.NewScanner(file)
	docs := make([]interface{}, ImportBatchSize)
	for scanner.Scan() {
		line := scanner.Text()
		var doc interface{}
		err := bson.UnmarshalExtJSON([]byte(line), false, &doc)
		if err != nil {
			log.Printf("Invalid format: %v\n", err)
			continue
		}

		docs = append(docs, doc)

		batchCount++
		// Flush the docs and reset batchCount counter
		if batchCount >= ImportBatchSize {
			if _, err := collection.InsertMany(ctx, docs); err != nil {
				log.Printf("Error: %v\n", err)
			}
			batchCount = 0
			docs = docs[:0]
			fmt.Print(".")
		}
		totalCount++
	}

	// Check the last batch
	if batchCount > 0 {
		log.Printf("docs: %v\n", docs)
		if _, err := collection.InsertMany(ctx, docs); err != nil {
			log.Printf("Error: %v\n", err)
		}
	}
	fmt.Print("\n")

	return totalCount
}

func saveCollectionToFile(db *mongo.Database, collectionName string, dir string) {
	collection := db.Collection(collectionName)
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	filePath := dir + "/" + db.Name() + "." + collectionName + ".bson"
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	count := 0
	for cursor.Next(context.TODO()) {
		var result bson.D
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		res, _ := bson.MarshalExtJSON(result, false, false)
		line := fmt.Sprintf("%+v\n", string(res))
		file.WriteString(line)
		count++
	}
	log.Printf("Exported %d document(s) in collection: %s to file -> %s", count, collectionName, filePath)
}

func listCollections(db *mongo.Database) []string {
	names, err := db.ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		panic(err)
	}
	return names
}

func cleanCollection(db *mongo.Database, name string) int64 {
	collection := db.Collection(name)
	result, err := collection.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		log.Println("Unable to clean collection", err)
		return 0
	}
	return result.DeletedCount
}

func close(client *mongo.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Printf("Failed to close client. %v", err)
	}
}
