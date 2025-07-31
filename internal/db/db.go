package db

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var MG MongoInstance

const dbName = "hrms"
const MongoURI = "mongodb://localhost:27017/" + dbName
const EmployeeCollection = "employees"

func ConnectDb() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoURI))
	if err != nil {
		return err
	}
	db := client.Database(dbName)
	MG = MongoInstance{
		Client: client,
		Db:     db,
	}
	return nil
}

func GetEmployeeCollection() *mongo.Collection {
	return MG.Db.Collection(EmployeeCollection)
}

func GetUploadDir() (string, error) {
	basePath, err := os.Getwd()
	if err != nil {
		return "/", err
	}
	filePath := filepath.Join(basePath, "internal", "db", "upload")
	return filePath, nil
}
