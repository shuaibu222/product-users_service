package main

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

const mongoURL = "mongodb://host.minikube.internal:27017"

// connect with MongoDB
func init() {
	credential := options.Credential{
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	}
	clientOpts := options.Client().ApplyURI(mongoURL).SetAuth(credential)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Println("Error connecting to MongoDB")
		return
	}

	collection = client.Database("users").Collection("users")

	// collection instance
	log.Println("Collections instance is ready")
}

type UsersEntry struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
}

func (u *UsersEntry) Insert(entry UsersEntry, ctx context.Context) (*mongo.InsertOneResult, error) {
	user, err := collection.InsertOne(ctx, entry)
	if err != nil {
		log.Println("Error inserting into users:", err)
	}

	return user, nil
}

func (u *UsersEntry) GetAllUsers(ctx context.Context) ([]*UsersEntry, error) {

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Finding all users error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*UsersEntry

	for cursor.Next(ctx) {
		var item UsersEntry

		err := cursor.Decode(&item)
		if err != nil {
			log.Print("Error decoding log into slice:", err)
			return nil, err
		} else {
			users = append(users, &item)
		}
	}

	return users, nil
}

func (u *UsersEntry) GetOneUserById(id string, ctx context.Context) (*UsersEntry, error) {
	var user *UsersEntry
	Id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Failed to convert id", err)
		return nil, err
	}

	filter := bson.M{"_id": Id}
	err = collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Println("Failed to get user", err)
		return nil, err
	}

	return user, nil
}

// for authentication purposes only
func (u *UsersEntry) GetUserByUsername(username string, ctx context.Context) (UsersEntry, error) {
	var user UsersEntry
	filter := bson.M{"username": username}
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Println("Failed to decode user: ", err)
	}

	return user, nil
}

func (u *UsersEntry) UpdateUserById(id string, userToUpdate UsersEntry, ctx context.Context) (*mongo.UpdateResult, error) {
	Id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Failed to convert id", err)
		return nil, err
	}

	updatedCount, err := collection.UpdateOne(ctx, bson.M{"_id": Id}, bson.M{"$set": userToUpdate})
	if err != nil {
		log.Println("Failed to update the user: ", err)
		return nil, err
	}
	return updatedCount, nil
}

func (u *UsersEntry) DeleteUser(id string, ctx context.Context) (*mongo.DeleteResult, error) {
	Id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Failed to convert id", err)
		return nil, err
	}

	deleteCount, err := collection.DeleteOne(ctx, bson.M{"_id": Id})
	if err != nil {
		log.Println("Failed to delete user", err)
		return nil, err
	}

	return deleteCount, nil
}
