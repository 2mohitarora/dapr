package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	dapr "github.com/dapr/go-sdk/client"
)

const cartActorType = "CartActor"

var daprClient dapr.Client

// AddItemRequest mirrors the actor's input type.
type AddItemRequest struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

func getCart(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")

	resp, err := daprClient.InvokeActor(r.Context(), &dapr.InvokeActorRequest{
		ActorType: cartActorType,
		ActorID:   userID,
		Method:    "GetCart",
	})
	if err != nil {
		log.Printf("frontendsvc: get cart: %s", err)
		http.Error(w, "unable to get cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp.Data)
}

func addItem(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")

	var req AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("frontendsvc: decode add item: %s", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	data, err := json.Marshal(req)
	if err != nil {
		log.Printf("frontendsvc: marshal add item: %s", err)
		http.Error(w, "unable to add item", http.StatusInternalServerError)
		return
	}

	resp, err := daprClient.InvokeActor(r.Context(), &dapr.InvokeActorRequest{
		ActorType: cartActorType,
		ActorID:   userID,
		Method:    "AddItem",
		Data:      data,
	})
	if err != nil {
		log.Printf("frontendsvc: add item: %s", err)
		http.Error(w, "unable to add item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp.Data)
}

func removeItem(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	itemID := r.PathValue("itemId")

	data, _ := json.Marshal(map[string]string{"item_id": itemID})

	resp, err := daprClient.InvokeActor(r.Context(), &dapr.InvokeActorRequest{
		ActorType: cartActorType,
		ActorID:   userID,
		Method:    "RemoveItem",
		Data:      data,
	})
	if err != nil {
		log.Printf("frontendsvc: remove item: %s", err)
		http.Error(w, "unable to remove item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp.Data)
}

func checkout(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")

	resp, err := daprClient.InvokeActor(r.Context(), &dapr.InvokeActorRequest{
		ActorType: cartActorType,
		ActorID:   userID,
		Method:    "Checkout",
	})
	if err != nil {
		log.Printf("frontendsvc: checkout: %s", err)
		http.Error(w, "unable to checkout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp.Data)
}

func clearCart(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")

	_, err := daprClient.InvokeActor(r.Context(), &dapr.InvokeActorRequest{
		ActorType: cartActorType,
		ActorID:   userID,
		Method:    "Clear",
	})
	if err != nil {
		log.Printf("frontendsvc: clear cart: %s", err)
		http.Error(w, "unable to clear cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"cleared"}`)
}

func main() {
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	dc, err := dapr.NewClient()
	if err != nil {
		log.Fatalf("frontendsvc: dapr client: %s", err)
	}
	daprClient = dc
	defer daprClient.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /cart/{userId}", getCart)
	mux.HandleFunc("POST /cart/{userId}/items", addItem)
	mux.HandleFunc("DELETE /cart/{userId}/items/{itemId}", removeItem)
	mux.HandleFunc("POST /cart/{userId}/checkout", checkout)
	mux.HandleFunc("DELETE /cart/{userId}", clearCart)

	log.Printf("frontendsvc: starting service: port %s", appPort)

	if err := http.ListenAndServe(":"+appPort, mux); err != nil {
		log.Fatalf("frontendsvc: %s", err)
	}
}
