package main

import (
	"encoding/json"
	"fmt"
	"github/ccclin/rakuten_f2e/model"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var (
	originAllowed = os.Getenv("ORIGIN_ALLOWED")
	port          = os.Getenv("PORT")
)

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{originAllowed},
		AllowedMethods: []string{"GET", "OPTIONS", "POST", "PATCH", "DELETE"},
	})
	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1/").Subrouter()
	api.HandleFunc("/users", usersIndexHandel).Methods("GET", "OPTIONS")
	api.HandleFunc("/users", usersCreateHandel).Methods("POST")
	api.HandleFunc("/users/{id}", usersUpdateHandel).Methods("PATCH", "OPTIONS")
	api.HandleFunc("/users/{id}", usersDeleteHandel).Methods("DELETE")

	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), c.Handler(api)))
}

func usersIndexHandel(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}

	var usersDoc model.UsersDoc
	users, err := usersDoc.All()
	if err != nil {
		log.Printf("usersDoc error: %v", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	usersJSON, _ := json.Marshal(users)
	_, _ = w.Write(usersJSON)
}

func usersCreateHandel(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %v", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}

	var user model.User
	err = json.Unmarshal(data, &user)
	if err != nil {
		log.Printf("json.Unmarshal error: %v", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}
	if !user.CheckPhone() {
		http.Error(w, "Wrong phrone format", http.StatusBadRequest)
		return
	}

	if isPass, _ := user.CheckName(); !isPass {
		http.Error(w, "Same name exist", http.StatusBadRequest)
		return
	}

	err = user.Create()
	if err != nil {
		log.Printf("user create error: %v", err)
		http.Error(w, "Error created request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	userJSON, _ := json.Marshal(user)
	_, _ = w.Write(userJSON)
}

func usersUpdateHandel(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}

	vars := mux.Vars(r)
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %v", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}
	var user model.User
	err = user.Find(vars["id"])
	if err != nil {
		log.Printf("user find error: %v", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}

	var newUser model.User
	err = json.Unmarshal(data, &newUser)
	if err != nil {
		log.Printf("json.Unmarshal error: %v", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}
	if !newUser.CheckPhone() {
		http.Error(w, "Wrong phrone format", http.StatusBadRequest)
		return
	}

	if isPass, _ := newUser.CheckName(); !isPass {
		http.Error(w, "Same name exist", http.StatusBadRequest)
		return
	}

	user.Merge(newUser)
	err = user.Update(vars["id"])
	if err != nil {
		log.Printf("user update error: %v", err)
		http.Error(w, "Error update request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	userJSON, _ := json.Marshal(user)
	_, _ = w.Write(userJSON)
}

func usersDeleteHandel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var user model.User
	err := user.Delete(vars["id"])
	if err != nil {
		http.Error(w, "Error delete request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"messag": "Delete finish"}`))
}
