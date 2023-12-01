package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shuaibu222/p-users/users"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// here should be the grpc handlers
type UsersServer struct {
	users.UsersServiceServer
}

func (u *UsersServer) CreateUser(ctx context.Context, req *users.UserRequest) (*users.UserResponse, error) {
	input := req.GetUserEntry()

	// hash the user entered password for security purposes
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error generating password: ", err)
	}

	input.Password = string(hashedPassword)

	userEntry := UsersEntry{
		Username: input.Username,
		Password: input.Password,
	}

	user, err := userEntry.Insert(userEntry, ctx)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
	}

	log.Println("User inserted successfully: ", user)

	var entry UsersEntry
	err = collection.FindOne(ctx, bson.M{"_id": user.InsertedID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	app := Config{}
	app.SendUserToRabbitmq(entry)

	insertedUser := &users.User{
		Id:       userEntry.ID.Hex(),
		Username: userEntry.Username,
		Password: userEntry.Password,
	}

	// return response
	res := &users.UserResponse{Response: insertedUser}
	return res, nil
}

func (u *UsersServer) GetUserByUsername(ctx context.Context, req *users.User) (*users.UserResponse, error) {
	userEntry := UsersEntry{}

	result, err := userEntry.GetUserByUsername(req.Username, ctx)
	if err != nil {
		log.Println(err)
	}

	userResult := &users.UserResponse{
		Response: &users.User{
			Id:       result.ID.String(),
			Username: result.Username,
			Password: result.Password,
		},
	}

	return userResult, nil

}

func (u *UsersServer) GetAllUsers(ctx context.Context, req *users.NoParams) (*users.UsersLists, error) {
	userEntry := UsersEntry{}

	usersResult, err := userEntry.GetAllUsers(ctx)
	if err != nil {
		log.Println("Failed to get users", err)
	}

	// Create the response containing the list of users
	var userResponses []*users.User
	for _, userEntry := range usersResult {
		userResponse := &users.User{
			Id:       userEntry.ID.Hex(),
			Username: userEntry.Username,
			Password: userEntry.Password,
		}
		userResponses = append(userResponses, userResponse)
	}

	// Return the response
	return &users.UsersLists{
		Response: userResponses,
	}, nil
}

func (u *UsersServer) GetUserById(ctx context.Context, req *users.UserId) (*users.UserResponse, error) {
	userEntry := UsersEntry{}

	result, err := userEntry.GetOneUserById(req.Id, ctx)
	if err != nil {
		log.Println(err)
	}

	userResult := &users.UserResponse{
		Response: &users.User{
			Id:       result.ID.String(),
			Username: result.Username,
			Password: result.Password,
		},
	}

	return userResult, nil
}

func (u *UsersServer) UpdateUser(ctx context.Context, req *users.User) (*users.Count, error) {
	userEntry := UsersEntry{
		Username: req.Username,
	}

	updateCount, err := userEntry.UpdateUserById(req.Id, userEntry, ctx)
	if err != nil {
		log.Println(err)
	}

	log.Println(updateCount.ModifiedCount)

	userResult := &users.Count{
		Count: fmt.Sprint(updateCount.ModifiedCount),
	}

	return userResult, nil
}

func (u *UsersServer) DeleteUser(ctx context.Context, req *users.UserId) (*users.Count, error) {
	userEntry := UsersEntry{}

	deletedCount, err := userEntry.DeleteUser(req.Id, ctx)
	if err != nil {
		log.Println(err)
	}

	res := &users.Count{
		Count: fmt.Sprint(deletedCount.DeletedCount),
	}

	return res, nil
}
