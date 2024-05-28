package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Restaurant struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

var restaurants []Restaurant
var restaurantIDCounter = 1

func main() {
	// Rotas da API
	http.HandleFunc("/restaurants", getRestaurants)
	http.HandleFunc("/restaurants/add", addRestaurant)
	http.HandleFunc("/restaurants/update", updateRestaurant)
	http.HandleFunc("/restaurants/delete", deleteRestaurant)

	// Inicie o servidor na porta 8080
	fmt.Println("Servidor rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getRestaurants(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(restaurants)
}

func addRestaurant(w http.ResponseWriter, r *http.Request) {
	var newRestaurant Restaurant
	err := json.NewDecoder(r.Body).Decode(&newRestaurant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newRestaurant.ID = restaurantIDCounter
	restaurantIDCounter++
	restaurants = append(restaurants, newRestaurant)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newRestaurant)
}

func updateRestaurant(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}

	var updatedRestaurant Restaurant
	err = json.NewDecoder(r.Body).Decode(&updatedRestaurant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, restaurant := range restaurants {
		if restaurant.ID == id {
			restaurants[i] = updatedRestaurant
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedRestaurant)
			return
		}
	}

	http.Error(w, "Restaurant not found", http.StatusNotFound)
}

func deleteRestaurant(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}

	for i, restaurant := range restaurants {
		if restaurant.ID == id {
			restaurants = append(restaurants[:i], restaurants[i+1:]...)
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(w, "Restaurant not found", http.StatusNotFound)
}
