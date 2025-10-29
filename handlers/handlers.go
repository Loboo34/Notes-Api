package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Loboo34/Notes-Api/database"
	"github.com/Loboo34/Notes-Api/logger"
	"github.com/Loboo34/Notes-Api/models"
	"github.com/Loboo34/Notes-Api/utils"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

//var notes []models.Note

// add notes to list
func AddNote(w http.ResponseWriter, r *http.Request) {
	//verify http method
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only Post Allowed", "")
		return
	}

	var note models.Note //title,content,tags

	//reads/decode  json
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		logger.Log.Warn("Invalid Json", zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Json", "")
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
		logger.Log.Warn("Error saving Note", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Error Saving note", "")
	}

	//set content header
	//w.Header().Set("Content-Type", "application/json")
	//writes to client-encode to json
	// json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("%v added to notes", note.Title)})
	//fmt.Println("Added task")

	logger.Log.Info("Note added successfully", zap.String("title", note.Title))

	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": fmt.Sprintf("%v added successfully", note.Title)})
}

func UpdateNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only Put Allowed", "")
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		//logger.Log.Warn("Missing id param", zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, "Missing Id Param", "")
		return
	}
	idStr = strings.TrimSpace(idStr)

	objectId, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		logger.Log.Warn("Invalid id format", zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid id format", "")
		return
	}

	var updated models.Note
	if err = json.NewDecoder(r.Body).Decode(&updated); err != nil {
		logger.Log.Warn("Invalid json", zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid json format", "")
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

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectId}, update)
	if err != nil {
		logger.Log.Warn("Failed to update", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Database update error", "")
		return
	}

	if result.MatchedCount == 0 {
		logger.Log.Warn("Failed to update", zap.Error(err))
		utils.RespondWithError(w, http.StatusNotFound, "Note not found", "")
		return
	}

	logger.Log.Info("Updated Task successfully", zap.String("Title", updated.Title))

	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": fmt.Sprintf("%v updated ", updated.Title)})

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(map[string]string{
	// 	"message": "Note updated successfully",
	// })

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
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only Delete Alowed", "")
		return
	}
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing Id param", "")
		return
	}
	idStr = strings.TrimSpace(idStr)

	objectId, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		logger.Log.Error("invalid id format")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid id format", "")
		return
	}

	collection := database.DB.Collection("notes")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		logger.Log.Warn("Failed to delete note")
		utils.RespondWithError(w, http.StatusInternalServerError, "Database delete Error", "")
		return
	}

	if result.DeletedCount == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Note not found", "")
		return
	}

	logger.Log.Warn("Delete successful")

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Note deleted successfully"})
	// json.NewEncoder(w).Encode(map[string]string{
	// 	"message": "Note deleted successfully",
	// })w.Header().Set("Content-Type", "application/json")

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
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only Get Allowed", "")
		return
	}

	collection := database.Client.Database("notesdb").Collection("notes")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		logger.Log.Warn("Failed to fetch notes", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching notes", "")
		return
	}

	defer cursor.Close(ctx)

	var notes []models.Note

	for cursor.Next(ctx) {
		var note models.Note
		if err := cursor.Decode(&note); err != nil {
			logger.Log.Warn("Failed to decoding  notes", zap.Error(err))
			utils.RespondWithError(w, http.StatusInternalServerError, "Error Decoding note", "")
			return
		}
		notes = append(notes, note)
	}

	if err := cursor.Err(); err != nil {
		logger.Log.Warn("Cursor erro", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Cursor Error", "")
	}

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(notes)

	logger.Log.Info("Fetched all tasks")

	utils.RespondWithJSON(w, http.StatusOK, notes)
}

func GetNoteById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only Get Allowed", "")
		return
	}
	vars := mux.Vars(r)
	idStr := vars["is"]
	if idStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing id Param", "")
		return
	}

	objextId, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		logger.Log.Warn("Invalid id format", zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid id format", "")
		return
	}

	collection := database.DB.Collection("notes")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var note models.Note
	err = collection.FindOne(ctx, bson.M{"_id": objextId}).Decode(&note)
	if err != nil {
		logger.Log.Warn("Note not founs", zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, "Note note found", "")
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(note)

	utils.RespondWithJSON(w, http.StatusOK, note)
}
