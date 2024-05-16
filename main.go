package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {

	// create router
	route := mux.NewRouter()
	router := route.PathPrefix("/api/v1").Subrouter() //Base Path
	router.HandleFunc("/createProfile", createProfile).Methods("POST")
	router.HandleFunc("/getUserProfile", getUserProfile).Methods("GET")
	router.HandleFunc("/getAllUsers", getAllUsers).Methods("GET")
	router.HandleFunc("/updateProfile", updateProfile).Methods("PATCH")
	router.HandleFunc("/deleteProfile", deleteProfile).Methods("GET")

	// Run Server
	log.Fatal(http.ListenAndServe(":8000", jsonContentTypeMiddleware(router)))

}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
