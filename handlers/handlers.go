package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"notes/database"
	"notes/models"
	"notes/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//var notes []models.Note

// add notes to list
func AddNote(w http.ResponseWriter, r *http.Request) {
	//verify http method
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only Post Allowed", "")
		return
	}

	var note models.Note

	//reads/decode  json
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	note.ID = primitive.NewObjectID() //add mongodb id
	note.CreatedAt = time.Now()       //add time stamp
	//notes = append(notes, note)

	//connect to collection
	collection := database.Client.Database("notesdb").Collection("notes")   //gets reference of the notes in the notesDb
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //creates context with 5 min time out
	defer cancel()                                                          //ensures context is cleaned when function exits

	//inserts note into db
	_, err = collection.InsertOne(ctx, note)
	if err != nil {
		http.Error(w, "Error Saving note", http.StatusInternalServerError)
	}

	//set content header
	w.Header().Set("Content-Type", "application/json")
	//writes to client-encode to json
	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("%v added to notes", note.Title)})

	fmt.Println("Added task")
}

func UpdateNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Only Put Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing Id Param", http.StatusBadRequest)
		return
	}
	idStr = strings.TrimSpace(idStr)

	objectId, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return
	}

	var updated models.Note
	if err = json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid json format", http.StatusBadRequest)
		return
	}

	collection := database.DB.Collection("notes")

	update := bson.M{
		"$set": bson.M{
			"title":     updated.Title,
			"content":   updated.Content,
			"tags":      updated.Tags,
			"updatedAt": time.Now(),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("Received ID:", idStr)
	fmt.Println("Converted ObjectID:", objectId)
	fmt.Println("Updating fields:", update)

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectId}, update)
	if err != nil {
		http.Error(w, "Database update error", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Note updated successfully",
	})

	// for idx, n := range notes {
	// 	if n.ID == primitive.NewObjectID() {
	// 		if updated.Title != "" {
	// 			notes[idx].Title = updated.Title
	// 		}
	// 		if updated.Content != "" {
	// 			notes[idx].Content = updated.Content
	// 		}
	// 		if len(updated.Tags) > 0 {
	// 			notes[idx].Tags = updated.Tags
	// 		}

	// 		w.Header().Set("Content-Type", "application/json")
	// 		json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("%v note updated", id)})
	// 		return
	// 	}
	// }
	//http.Error(w, "Note Not found", http.StatusNotFound)
}

func DeleteNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Only Delete Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing Id param", http.StatusBadRequest)
		return
	}
	idStr = strings.TrimSpace(idStr)

	objectId, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return
	}

	collection := database.DB.Collection("notes")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		http.Error(w, "Database delete Error", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Note deleted successfully",
	})

	// for idx, n := range notes {
	// 	if n.ID == primitive.NewObjectID() {
	// 		notes = append(notes[:idx], notes[idx+1:]...)
	// 		w.Header().Set("Content-Type", "application/json")
	// 		json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("%v deleted ", id)})
	// 		return
	// 	}
	// }
	//http.Error(w, "Note Not found", http.StatusNotFound)
}

func GetNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only Get Allowed", http.StatusMethodNotAllowed)
		return
	}

	collection := database.Client.Database("notesdb").Collection("notes")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Error fetching Notes", http.StatusInternalServerError)
		return
	}

	defer cursor.Close(ctx)

	var notes []models.Note

	for cursor.Next(ctx) {
		var note models.Note
		if err := cursor.Decode(&note); err != nil {
			http.Error(w, "Error Decoding note", http.StatusInternalServerError)
			return
		}
		notes = append(notes, note)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, "Cursor Error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)

	fmt.Println("Fetched all notes")
}


func GetNoteById(w http.ResponseWriter, r *http.Request) {
if	r.Method == http.MethodGet {
		http.Error(w, "Only Get Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == ""{
		http.Error(w, "Missing id Param", http.StatusBadRequest)
		return
	}

	objextId, err := primitive.ObjectIDFromHex(idStr)
	if err !=nil {
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return
	}

	collection := database.DB.Collection("notes")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var note models.Note
	err = collection.FindOne(ctx, bson.M{"_id": objextId}).Decode(&note)
	if err != nil {
		http.Error(w, "Note note found", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}