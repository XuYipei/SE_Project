package tools

import (
	"context"

	"github.com/tuplz/tuplz-be/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDb() (*mongo.Client, error) {
	clientUri := "mongodb://" + config.ClientUser + ":" + config.ClientPassword + "@localhost:27017"
	clientOpts := options.Client().ApplyURI(clientUri)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		return client, err
	}

	err = client.Ping(context.TODO(), nil)
	return client, err
}

func ClearCollection(table string) error {
	client, err := ConnectDb()
	if err != nil {
		return err
	}
	collection := client.Database(config.Database).Collection(table)
	_, errDelete := collection.DeleteMany(context.TODO(), bson.M{})
	return errDelete
}

func DeleteCollection(table string) error {
	dbClient, errConnect := ConnectDb()
	if errConnect != nil {
		return errConnect
	}
	collection := dbClient.Database(config.Database).Collection(table)
	err := collection.Drop(context.Background())
	return err
}

func CreateCollection(table string) error {
	client, err := ConnectDb()
	if err != nil {
		return err
	}
	database := client.Database(config.Database)
	err = database.CreateCollection(context.Background(), table)
	return err
}

func InitCollection(table string) error {
	err := DeleteCollection(table)
	if err != nil {
		// log.print(err)
		return err
	}
	err = CreateCollection(table)
	return err
}

func FindOne(db string, table string, filter bson.M) *mongo.SingleResult {
	client, err := ConnectDb()
	if err != nil {
		return nil
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(db).Collection(table)
	return collection.FindOne(context.Background(), filter)
}

func FindMany(db string, table string, filter bson.M) (*mongo.Cursor, error) {
	client, err := ConnectDb()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(db).Collection(table)
	cur, err := collection.Find(context.Background(), filter)
	return cur, err
}

func FindManyOpt(db string, table string, filter bson.M, opts *options.FindOptions) (*mongo.Cursor, error) {
	client, err := ConnectDb()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(db).Collection(table)
	cur, err := collection.Find(context.Background(), filter, opts)
	return cur, err
}

func InsertOne(db string, table string, data interface{}) error {
	client, err := ConnectDb()
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(db).Collection(table)
	_, err = collection.InsertOne(context.Background(), data)
	return err
}

func UpdateOne(db string, table string, filter bson.M, data interface{}) error {
	dbClient, err := ConnectDb()
	if err != nil {
		return err
	}
	defer dbClient.Disconnect(context.Background())

	collection := dbClient.Database(db).Collection(table)
	_, err = collection.UpdateOne(context.Background(), filter, data)
	return err
}

func DeleteOne(db string, table string, filter bson.M) error {
	dbClient, errConnect := ConnectDb()
	if errConnect != nil {
		return nil
	}
	defer dbClient.Disconnect(context.Background())

	collection := dbClient.Database(db).Collection(table)
	_, errDelete := collection.DeleteOne(context.Background(), filter)
	return errDelete
}

func DeleteMany(dbName string, collectionName string, filter bson.M) (int, error) {
	dbClient, err := ConnectDb()
	if err != nil {
		return 0, err
	}
	defer dbClient.Disconnect(context.Background())

	collection := dbClient.Database(dbName).Collection(collectionName)
	result, err := collection.DeleteMany(context.Background(), filter)
	return int(result.DeletedCount), err
}

func DbMax(dbName string, collectionName string, index string) (*mongo.Cursor, error) {
	dbClient, err := ConnectDb()
	if err != nil {
		return nil, err
	}

	opt := options.Find()
	opt.SetSort(bson.M{index: -1})
	opt.SetLimit(1)

	collection := dbClient.Database(dbName).Collection(collectionName)
	cur, err := collection.Find(context.Background(), bson.M{}, opt)
	return cur, err
}

func CollStatus(dbName string, collName string) (*mongo.SingleResult, error) {
	dbClient, err := ConnectDb()
	if err != nil {
		return nil, err
	}

	database := dbClient.Database(dbName)
	result := database.RunCommand(context.Background(), bson.M{"collStats": collName})
	return result, nil
}

func RunCommandDb(db string, opt interface{}) error {
	dbClient, err := ConnectDb()
	if err != nil {
		return err
	}
	database := dbClient.Database(db)
	result := database.RunCommand(context.Background(), opt)
	// log.print(result.Err())
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func LockDb() error {
	arg := struct {
		Fsync int `bson:"fsync"`
		Lock  int `bson:"lock"`
	}{Fsync: 1, Lock: 1}
	return RunCommandDb("admin", arg)
}

func UnlockDb() error {
	arg := struct {
		Unlock  int    `bson:"fsyncUnlock"`
		Comment string `bson:"comment"`
	}{Unlock: 1, Comment: "unlock"}
	return RunCommandDb("admin", arg)
}
