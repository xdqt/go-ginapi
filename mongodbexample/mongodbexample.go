package mongodbexample

import (
	"context"
	"fmt"
	"ginapi/structs"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

var testCollection *mongo.Collection

func InitMongo() {
	clientOptions := options.Client().ApplyURI("mongodb://ellis:ellischen@192.168.214.133:32000/")
	mongoClient, _ := mongo.Connect(context.TODO(), clientOptions)
	testCollection = mongoClient.Database("baz").Collection("qux")
}

func InsertOneByStruct() {
	res, err := testCollection.InsertOne(context.Background(), &structs.MongoStruct{Id: primitive.NewObjectID(), UserName: "ellis", Email: "849773373@qq.com"})
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	id := res.InsertedID
	fmt.Printf("id: %v\n", id)
}

func InsertManyByStructs() {
	values := []interface{}{structs.MongoStruct{Id: primitive.NewObjectID(), UserName: "1", Email: "1"}, structs.MongoStruct{Id: primitive.NewObjectID(), UserName: "2", Email: "2"}}
	imr, _ := testCollection.InsertMany(context.Background(), values)
	fmt.Printf("imr.InsertedIDs: %v\n", imr.InsertedIDs)
}

func FindALL() {
	ctx, channel := context.WithTimeout(context.Background(), 30*time.Second)
	defer channel()

	// cur, _ := testCollection.Find(ctx, bson.M{"username": "1"})
	cur, _ := testCollection.Find(ctx, bson.D{{"username", "1"}})
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var value structs.MongoStruct
		cur.Decode(&value)
		fmt.Printf("value: %v\n", value)
	}
}

func UpdateMany() {
	ctx, channel := context.WithTimeout(context.Background(), 30*time.Second)
	defer channel()
	ur, err := testCollection.UpdateMany(ctx, bson.D{{"username", "vv"}}, bson.D{{"$set", bson.D{{"username", "ellis"}, {"email", "haha"}}}})
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("ur.MatchedCount: %v\n", ur.MatchedCount)
}

func DeleteOne() {
	ctx, channel := context.WithTimeout(context.Background(), 30*time.Second)
	defer channel()
	dr, err := testCollection.DeleteOne(ctx, bson.D{{"username", "1"}})
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("dr.DeletedCount: %v\n", dr.DeletedCount)
}

// func main() {
// 	InitMongo()
// 	// InsertOneByStruct()
// 	// InsertManyByStructs()
// 	// FindALL()
// 	// UpdateMany()
// 	DeleteOne()
// }
