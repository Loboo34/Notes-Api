package database

import (
	"context"//controlls lifetime of a db
	"fmt"
	"time"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client//reps mongo client connection
var DB *mongo.Database//reps actual db


func ConnectDB() *mongo.Database {	
	var mongoUri = os.Getenv("MONGO_URI") 
	clientOptions := options.Client().ApplyURI(mongoUri)//connects to local db
	//options.Client()-creates a stuct/obj that holds:-connection uri,timeout settings, pool size, Auth creds

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//creates context and cancels it after 10secs
	//context.Bacground()-creates base context-used to controll timeouts,cancel ops and to pass meta data
	//wrapping using context.WithTimeout() creates child context(ctx) that will cancel after 10 secs
	//ctx-context to pass to mongo ops

	defer cancel()//ensures resources are cleared up after connection
	

	client, err := mongo.Connect(ctx, clientOptions)//connects to db using given context(ctx) and potions
	//pass clientOptions to mongo.connect creating an actual mongo client connection
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, nil)//tests connection ie checks if mongo responds after a single ping
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to db Successfully")
	Client = client//stores client globally for reuse
	DB = client.Database("notesdb")//specifies db being called
	return DB
}
