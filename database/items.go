package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/hendraprawira/go-procurement/graph/model"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
}

// type EmailCreds struct {
// 	host         string
// 	port         int
// 	senderName   string
// 	authEmail    string
// 	authPassword string
// }

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func Connect() *DB {
	client, err := mongo.NewClient(options.Client().ApplyURI(goDotEnvVariable("MONGODB")))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	return &DB{
		client: client,
	}
}

func (db *DB) Save(input model.NewItem) *model.Item {
	collection := db.client.Database("procurementTripatra").Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, input)
	if err != nil {
		log.Fatal(err)
	}
	return &model.Item{
		ID:          res.InsertedID.(primitive.ObjectID).Hex(),
		NameItem:    input.NameItem,
		Stock:       input.Stock,
		Description: input.Description,
		Price:       input.Price,
	}

}

func (db *DB) Delete(ID string) *bool {
	ObjectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		log.Fatal(err)
	}
	collection := db.client.Database("procurementTripatra").Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, error := collection.DeleteOne(ctx, bson.M{"_id": ObjectID})
	if error != nil {
		log.Fatal(err)
	}
	success := new(bool)
	*success = true
	return success
}

func (db *DB) FindByID(ID string) *model.Item {
	ObjectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		log.Fatal(err)
	}
	collection := db.client.Database("tripartafurniture").Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res := collection.FindOne(ctx, bson.M{"_id": ObjectID})
	item := model.Item{}
	res.Decode(&item)
	return &item
}

func (db *DB) Find(input *model.FilterItem) []*model.Item {
	collection := db.client.Database("procurementTripatra").Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	findQuery := bson.M{}

	if input.NameItem != nil {
		findQuery["NameItem"] = input.NameItem
	}
	cur, err := collection.Find(ctx, findQuery)
	if err != nil {
		log.Fatal(err)
	}
	var items []*model.Item
	for cur.Next(ctx) {
		var item *model.Item
		err := cur.Decode(&item)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}
	return items
}
