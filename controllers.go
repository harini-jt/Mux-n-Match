package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// struct for storing data
type User struct {
	UserID string `json:"userid"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

var userCollection = db().Database("MuxMatchProfiles").Collection("users")

// CREATE
func createProfile(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	// get the post body - user details
	var body User
	err := json.NewDecoder(req.Body).Decode(&body)

	if err != nil {
		fmt.Println(err)
	}

	insertResult, err := userCollection.InsertOne(context.TODO(), body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult)
	json.NewEncoder(res).Encode(insertResult.InsertedID) // return the mongodb ID of generated document
}

// GET BY ID
func getUserProfile(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	// Extract the user ID from the request URL parameters
	userID := req.URL.Query().Get("userid")
	if userID == "" {
		http.Error(res, "No user ID provided", http.StatusBadRequest)
		return
	}
	var result User // Define a variable to hold the user details
	err := userCollection.FindOne(context.TODO(), bson.D{{Key: "userid", Value: userID}}).Decode(&result)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "User not found", http.StatusNotFound)
		return
	}

	// Return the user details as JSON response
	json.NewEncoder(res).Encode(result)
}

// READ
func getAllUsers(res http.ResponseWriter, _ *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var results []primitive.M                                   //slice for multiple documents
	cur, err := userCollection.Find(context.TODO(), bson.D{{}}) //returns a *mongo.Cursor
	if err != nil {
		fmt.Println(err)
	}
	for cur.Next(context.TODO()) {
		//Next() gets the next document for corresponding cursor
		var elem primitive.M
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem) // appending document pointed by Next()
	}
	cur.Close(context.TODO()) // close cursor after streaming documents
	json.NewEncoder(res).Encode(results)
}

// UPDATE
// func updateProfile(res http.ResponseWriter, req *http.Request) {

// 	res.Header().Set("Content-Type", "application/json")

// 	type updateUser struct {
// 		Name string `json:"name"` //value that has to be matched
// 		City string `json:"city"` // value that has to be modified
// 	}
// 	var body updateUser
// 	e := json.NewDecoder(req.Body).Decode(&body)
// 	if e != nil {

// 		fmt.Print(e)
// 	}
// 	filter := bson.D{{Key: "name", Value: body.Name}} // converting value to BSON type
// 	after := options.After                            // for returning updated document
// 	returnOpt := options.FindOneAndUpdateOptions{

// 		ReturnDocument: &after,
// 	}
// 	update := bson.D{{Key: "$set", Value: bson.D{{Key: "city", Value: body.City}}}}
// 	updateResult := userCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

// 	var result primitive.M
// 	_ = updateResult.Decode(&result)

//		json.NewEncoder(res).Encode(result)
//	}
//
// UPDATE
func updateProfile(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	// Extract the user ID from the request URL parameters
	userID := req.URL.Query().Get("userid")
	if userID == "" {
		http.Error(res, "No user ID provided", http.StatusBadRequest)
		return
	}

	// Decode the request body into updateUser struct
	var body map[string]interface{}
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Define the filter to match the user by userID
	filter := bson.D{{Key: "userid", Value: userID}}

	// Define the update operation
	update := bson.D{{Key: "$set", Value: body}}

	// Ensure that the 'userID', and 'email' fields are not modified
	delete(body, "userid")
	delete(body, "email")

	// Set options to return the updated document
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)

	// Perform the update operation
	updateResult := userCollection.FindOneAndUpdate(context.TODO(), filter, update, options)
	if updateResult.Err() != nil {
		http.Error(res, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// Decode the updated document into a map
	var updatedUser map[string]interface{}
	err = updateResult.Decode(&updatedUser)
	if err != nil {
		http.Error(res, "Failed to decode updated user", http.StatusInternalServerError)
		return
	}

	// Return the updated user details as JSON response
	json.NewEncoder(res).Encode(updatedUser)
}

// DELETE
func deleteProfile(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	// Extract the user ID from the request URL parameters
	userID := req.URL.Query().Get("id")
	if userID == "" {
		http.Error(res, "No user ID provided", http.StatusBadRequest)
		return
	}

	// Convert the user ID to the MongoDB ObjectID type
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(res, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Define the filter to match the user by ID
	filter := bson.D{{Key: "_id", Value: objID}}

	// Delete the user document from the database
	deleteResult, err := userCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	// Check if a document was deleted
	if deleteResult.DeletedCount == 0 {
		http.Error(res, "User not found", http.StatusNotFound)
		return
	}

	// Return success message
	json.NewEncoder(res).Encode(map[string]string{"message": "User deleted successfully"})
}
