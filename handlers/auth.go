package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

		"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo"
	

	"github.com/Loboo34/Notes-Api/models"
	"github.com/Loboo34/Notes-Api/database"
	"github.com/Loboo34/Notes-Api/utils"
	"github.com/Loboo34/Notes-Api/logger"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only Post Allowed", "")
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Json", "")
		return
	}

	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		logger.Log.Warn("Failed to hash password")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error hashing pass", "")
		return
	}

	user.Password = hashed

	collection := database.DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_,err = collection.InsertOne(ctx,user)
	if err != nil {
		logger.Log.Warn("Failed to register user")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error saving user", "")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "User registered successfully"})

}

func LoginUser(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only post allowed", "")
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Invalid Json", "")
		return
	}
collection := database.DB.Collection("users")

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

var creds models.User

err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&creds)
if err != nil {
utils.RespondWithError(w, http.StatusUnauthorized, "invalid credentials", "")
return
}

token, err := utils.GenerateJWT(user.Email)
if err != nil {
	utils.RespondWithError(w, http.StatusInternalServerError, "Error generating token", err.Error())
		return
}

utils.RespondWithJSON(w, http.StatusOK, map[string]string{"token": token})

}
