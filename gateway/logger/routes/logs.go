/* Logs

Contains the API routes to handle kong logs
- PostLog: Adds a Kong log to the database

Joaquin Badillo
2024-04-14
*/

package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"ccs/logger/db"
	"ccs/logger/lib"
	"ccs/logger/models"
)

func PostLog(w http.ResponseWriter, r *http.Request) {
	client, err := db.GetMongoClient()
	if err != nil {
		log.Printf("MongoDB error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	
	collection := client.Database("logs").Collection("gateway")	

	var rawLog models.KongLog
	if err := json.NewDecoder(r.Body).Decode(&rawLog); err != nil {
		log.Printf("Failed to marshal response for logging: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	entry := models.MapKongLogToEntry(&rawLog)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, entry)
	if err != nil {
		log.Printf("MongoDB insert error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Inserted log entry with ID: %v", res.InsertedID)

	lib.WriteResponse(res.InsertedID, w, http.StatusOK)
}
