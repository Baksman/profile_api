package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// note this should be in  your .env file
var secretKey []byte = []byte("testingsecretkey")

type Db struct {
	session    *mgo.Session
	collection *mgo.Collection
}

type User struct {
	ID       string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type UserProfile struct {
	ID       string   `json:"id" bson:"_id"`
	Username string   `json:"username" bson:"username"`
	hobbies  []string `json:"password" bson:"password"`
}

func (db *Db) createUser(w http.ResponseWriter, r *http.Request) {
	user := User{}
	w.Header().Set("Content-Type", "application/json")
	pBody, err := ioutil.ReadAll(r.Body)
	json.Unmarshal(pBody, &user)
	uuid, _ := uuid.NewV4()
	id := uuid.String()
	fmt.Println(user.Password)
	fmt.Println(user.Username)
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.ID = id
	user.Password = string(bytes)

	err = db.collection.Insert(user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	response := map[string]string{"status": "success", "message": "you can now login via login route"}
	// response := Response{Token: tokenString, Status: "success"}
	responseJSON, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)

	w.Write(responseJSON)
}

func (db *Db) login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := User{}
	pBody, err := ioutil.ReadAll(r.Body)
	json.Unmarshal(pBody, &user)
	fetchedUser := User{}
	err = db.collection.Find(bson.M{"username": user.Username}).One(&fetchedUser)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("User not found"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(fetchedUser.Password), []byte(user.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect password"))
		return
	}
	claims := jwt.MapClaims{
		"username": user.Username,
		"ExpiresAt": 15000,
		"IssuedAt": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)

		w.Write([]byte(err.Error()))
		return
		
	}
	response := map[string]string{"status": "success", "Token": tokenString, "id": fetchedUser.ID}
	// response := Response{Token: tokenString, Status: "success"}
	responseJSON, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	
	w.Write([]byte(responseJSON))
}

func main() {
	db := &Db{}
	var err error
	db.session, err = mgo.Dial("127.0.0.1:27017")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.session.Close()
	db.collection = db.session.DB("usersapi").C("userprofile")
	r := mux.NewRouter()
	r.HandleFunc("/login", db.login)
	r.HandleFunc("/sign-up", db.createUser)

	log.Fatal(http.ListenAndServe(":8080", r))
}
