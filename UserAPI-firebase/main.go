package main

import (
	"encoding/json"
	"fmt"
	"log"
	http "net/http"
	"os"
	"time"
	
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	firebase "firebase.google.com/go"
)

var mongodb_server = "127.0.0.1"
var mongodb_database = "airbnbClonedb"
var mongodb_collection = "user"
var mongo_admin_database = "admin"
var mongo_username = "admin"
var mongo_password = "cmpe281"

var app *App

func pingHandler(w http.ResponseWriter, req *http.Request) {
	log.Print("hello")
	mapD := map[string]string{"message": "API Working"}
	mapB, _ := json.Marshal(mapD)
	ResponseWithJSON(w, mapB, http.StatusOK)
	return
}

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	mapD := map[string]string{"message": message}
	mapB, _ := json.Marshal(mapD)
	ResponseWithJSON(w, mapB, code)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

//new user sign up
func RegisterUser(w http.ResponseWriter, req *http.Request) {
	var user User
	_ = json.NewDecoder(req.Body).Decode(&user)
	unqueId := uuid.Must(uuid.NewV4())
	user.UserId = unqueId.String()
	info := &mgo.DialInfo{
		Addrs:    []string{mongodb_server},
		Timeout:  60 * time.Second,
		Database: mongodb_database,
		Username: mongo_username,
		Password: mongo_password,
	}

	session, err := mgo.DialWithInfo(info)

	if err != nil {
		ErrorWithJSON(w, "Could not connect to database", http.StatusInternalServerError)
		return
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(mongodb_database).C(mongodb_collection)

	err = c.Insert(user)
	if err != nil {
		if mgo.IsDup(err) {
			ErrorWithJSON(w, "User with this ID already exists", http.StatusBadRequest)
			return
		}
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		return
	}

	respBody, err := json.MarshalIndent(user, "", "  ")
	ResponseWithJSON(w, respBody, http.StatusOK)
}

func UserSignIn(w http.ResponseWriter, req *http.Request) {
	var userData User
	_ = json.NewDecoder(req.Body).Decode(&userData)
	info := &mgo.DialInfo{
		Addrs:    []string{mongodb_server},
		Timeout:  60 * time.Second,
		Database: mongodb_database,
		Username: mongo_username,
		Password: mongo_password,
	}

	session, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(err)
		ErrorWithJSON(w, "Could not connect to database", http.StatusInternalServerError)
		return
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(mongodb_database).C(mongodb_collection)
	query := bson.M{"email": userData.Email,
		"password": userData.Password}
	
	var user User
	err = c.Find(query).One(&user)
	if err == mgo.ErrNotFound {
		ErrorWithJSON(w, "Login Failed", http.StatusUnauthorized)
		return
	}
	userData := bson.M{
		"email":   user.Email,
		"FirstName":    user.FirstName,
		"LastName": user.LastName,
		"UserId":      user.UserId}

	// create tokens
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatal("error getting Auth client: %v\n", err)
	}

	token, err := client.CustomToken(context.Background(), user.UserId)
	if err != nil {
		log.Fatal("error minting custom token: %v\n", err)
	}
	log.Printf("Got custom token: %v\n", token)	
		
	respBody, err := json.MarshalIndent(userData, "", "  ")
	ResponseWithJSON(w, respBody, http.StatusOK)
}

func DeleteUser(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	info := &mgo.DialInfo{
		Addrs:    []string{mongodb_server},
		Timeout:  60 * time.Second,
		Database: mongodb_database,
		Username: mongo_username,
		Password: mongo_password,
	}

	session, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(err)
		ErrorWithJSON(w, "Could not connect to database", http.StatusInternalServerError)
		return
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(mongodb_database).C(mongodb_collection)
	query := bson.M{"_id": params["UserId"]}
	err = c.Remove(query)

	if err != nil {
		switch err {
		default:
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed delete user: ", err)
			return
		case mgo.ErrNotFound:
			ErrorWithJSON(w, "User not found", http.StatusNotFound)
			return
		}
	}

	respBody, err := json.MarshalIndent("User deleted", "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	ResponseWithJSON(w, respBody, http.StatusOK)
}

func authenticationMiddleware(au http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//verify token
	token, err := client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
        log.Fatal("error verifying ID token: %v\n", err)
	}

	log.Printf("Verified ID token: %v\n", token)
	})
}

func initfirebase() {

	// initialize firebase sdk with service account
	opt := option.WithCredentialsFile("airbnb-clone.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
        log.Fatal("error initializing app: %v\n", err)
	}
}		
	

func main() {
	log.Print("hello")		
	
	router := mux.NewRouter()
	router.HandleFunc("/users/signup", RegisterUser).Methods("POST")
	router.HandleFunc("/users/signin", UserSignIn).Methods("POST")
	router.HandleFunc("/users/{id}", authenticationMiddleware(DeleteUser)).Methods("DELETE")
	// testing
    router.HandleFunc("/users/ping", pingHandler).Methods("GET")
	
	initfirebase()	
	log.Fatal(http.ListenAndServe(":3000", router))
}