package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// struct for storing data
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
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

// READ
func getUserProfile(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")

	var body User
	e := json.NewDecoder(req.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
	}
	var result primitive.M //  an unordered representation of a BSON document which is a Map
	err := userCollection.FindOne(context.TODO(), bson.D{{Key: "name", Value: body.Name}}).Decode(&result)
	if err != nil {

		fmt.Println(err)

	}
	json.NewEncoder(res).Encode(result) // returns a Map containing document

}
