package db

import (
	"context"
	"fmt"
	configs "helpdesk_backend/config"
	"helpdesk_backend/logger"
	"helpdesk_backend/utilities"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var MongoClient *mongo.Client

// Intialise and setup Mongo connection
func SetupMongoDB() {
	// config := config
	var err error
	uri := "mongodb://127.0.0.1:27017/test?directConnection=true"
	// uri := "mongodb+srv://umedrajawat26:umed271998@cluster0.hthaetu.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
	MongoClient, err = mongo.NewClient(options.Client().ApplyURI(uri).SetMaxPoolSize(1000))
	if err != nil {
		logger.Logger.Println("Error setting up mongo client")
		panic("[error] Error setting up mongo client" + uri)
	}
	fmt.Println("connected to Mongo", MongoClient)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = MongoClient.Connect(ctx)
	if err != nil {
		logger.Logger.Println("Error connecting to mongo server")
		panic("[error] Error connecting to mongo server")
	}
}

func GetCollection(coll string) *mongo.Collection {
	return MongoClient.Database(configs.MONGO_DB).Collection(coll)
}
func InsertMongoDocument(coll string, doc interface{}) error {
	db := configs.MONGO_DB
	var docInterface map[string]interface{}
	var err error

	out, err := bson.Marshal(doc)
	if err != nil {
		panic(err)
	}
	err = bson.Unmarshal(out, &docInterface)
	if err != nil {
		logger.Logger.Println("[error] | Error unmarshalling ", coll, err)
	}
	t := time.Now()
	docKeys := utilities.GetMapInterfaceKeys(doc)
	for _, key := range docKeys {
		if key == "created_at" || key == "updated_at" || key == "last_updated" {
			docInterface[key] = utilities.CLTime(t)
		}
	}

	retry := 0
	for {
		if retry >= 3 {
			logger.ZapLogger.Error("Error after 3 retries. Stopping find op", coll)
			break
		}
		_, err = MongoClient.Database(db).Collection(coll).InsertOne(context.TODO(), docInterface)
		if err != nil && err != mongo.ErrNoDocuments {
			logger.ZapLogger.Error("Error inserting", coll, err, retry)
			retry += 1
		} else if err == mongo.ErrNoDocuments {
			return err
		} else {
			break
		}
	}
	return err
}

func InsertManyMongo(coll string, docs []interface{}) error {
	db := configs.MONGO_DB
	var err error
	retry := 0
	for {
		if retry >= 3 {
			logger.ZapLogger.Error("Error after 3 retries. Stopping find op", coll)
			break
		}
		_, err := MongoClient.Database(db).Collection(coll).InsertMany(context.TODO(), docs)
		if err != nil && err != mongo.ErrNoDocuments {
			logger.ZapLogger.Error("Error inserting many ", coll, err, retry)
			retry += 1
		} else if err == mongo.ErrNoDocuments {
			return err
		} else {
			break
		}
	}
	return err
}

func UpdateMongoDocument(coll string, filter map[string]interface{}, update interface{}) error {
	db := configs.MONGO_DB
	var updateInterface map[string]interface{}
	var err error

	out, err := bson.Marshal(update)
	if err != nil {
		panic(err)
	}
	err = bson.Unmarshal(out, &updateInterface)
	if err != nil {
		logger.Logger.Println("[error] | Error unmarshalling ", coll, err)
		return err
	}
	t := time.Now()
	docKeys := utilities.GetMapInterfaceKeys(update)
	for _, key := range docKeys {
		if key == "updated_at" || key == "last_updated" {
			updateInterface[key] = utilities.CLTime(t)
		}
	}

	f := bson.M{}
	for k, v := range filter {
		f[k] = v
	}

	u := bson.M{"$set": updateInterface}

	retry := 0
	for {
		if retry >= 3 {
			logger.ZapLogger.Error("Error after 3 retries. Stopping find op", coll)
			break
		}
		_, err = MongoClient.Database(db).Collection(coll).UpdateOne(context.TODO(), f, u)
		if err != nil && err != mongo.ErrNoDocuments {
			logger.ZapLogger.Error("Error updating ", coll, err, retry)
			retry += 1
		} else if err == mongo.ErrNoDocuments {
			return err
		} else {
			break
		}
	}

	return err
}

func FindManyMongo(coll string, filter map[string]interface{}, options *options.FindOptions) (*mongo.Cursor, error) {
	db := configs.MONGO_DB
	q := bson.M{}
	for k, v := range filter {
		q[k] = v
	}
	var cursor *mongo.Cursor
	var err error
	retry := 0
	for {
		if retry >= 3 {
			logger.ZapLogger.Error("Error after 3 retries. Stopping find op", coll)
			break
		}
		cursor, err = MongoClient.Database(db).Collection(coll).Find(context.TODO(), q, options)
		if err != nil && err != mongo.ErrNoDocuments {
			logger.ZapLogger.Error("Error finding many ", coll, err, retry)
			retry += 1
		} else if err == mongo.ErrNoDocuments {
			return cursor, err
		} else {
			break
		}
	}

	return cursor, err
}

func FindOneMongo(coll string, filter map[string]interface{}, res interface{}) error {
	db := configs.MONGO_DB
	q := bson.M{}
	for k, v := range filter {
		q[k] = v
	}

	retry := 0
	var err error
	for {
		if retry >= 3 {
			logger.ZapLogger.Error("Error after 3 retries. Stopping find op", coll, filter)
			break
		}
		err := MongoClient.Database(db).Collection(coll).FindOne(context.TODO(), q).Decode(res)
		if err != nil && err != mongo.ErrNoDocuments {
			logger.ZapLogger.Error("Error finding", coll, filter, err, retry)
			retry += 1
		} else if err == mongo.ErrNoDocuments {
			return err
		} else {
			break
		}
	}
	return err
}

func DeleteManyMongo(coll string, filter map[string]interface{}) {
	db := configs.MONGO_DB
	q := bson.M{}
	for k, v := range filter {
		q[k] = v
	}
	retry := 0
	var delResults *mongo.DeleteResult
	var err error
	for {
		if retry >= 3 {
			logger.ZapLogger.Error("Error after 3 retries. Stopping find op", coll, filter)
			break
		}
		delResults, err = MongoClient.Database(db).Collection(coll).DeleteMany(context.TODO(), q)
		if err != nil && err != mongo.ErrNoDocuments {
			logger.ZapLogger.Error("Error deleting", coll, filter, err, retry)
			retry += 1
		} else {
			logger.ZapLogger.Info("While deleting", coll, filter, err, retry, delResults)
			break
		}
	}
}

func DeleteOneMongo(coll string, filter map[string]interface{}) (*mongo.DeleteResult, error) {
	db := configs.MONGO_DB
	q := bson.M{}
	for k, v := range filter {
		q[k] = v
	}

	retry := 0
	var delResults *mongo.DeleteResult
	var err error
	for {
		if retry >= 3 {
			logger.ZapLogger.Error("Error after 3 retries. Stopping find op", coll, filter)
			break
		}
		delResults, err = MongoClient.Database(db).Collection(coll).DeleteOne(context.TODO(), q)
		if err != nil && err != mongo.ErrNoDocuments {
			logger.ZapLogger.Error("Error deleting", coll, filter, err, retry)
			retry += 1
		} else {
			logger.ZapLogger.Info("While deleting", coll, filter, err, retry, delResults)
			break
		}
	}
	return delResults, err
}

func CountMongo(coll string, filter map[string]interface{}) (int64, error) {
	db := configs.MONGO_DB
	q := bson.M{}
	for k, v := range filter {
		q[k] = v
	}

	retry := 0
	var count int64
	var err error
	for {
		if retry >= 3 {
			logger.ZapLogger.Error("Error after 3 retries. Stopping find op", coll, filter)
			break
		}
		count, err = MongoClient.Database(db).Collection(coll).CountDocuments(context.TODO(), q)
		if err != nil && err != mongo.ErrNoDocuments {
			logger.ZapLogger.Error("Error counting", coll, filter, err, retry)
			retry += 1
		} else {
			logger.ZapLogger.Info("While counting", coll, filter, err, retry, count)
			break
		}
	}

	return count, err
}
