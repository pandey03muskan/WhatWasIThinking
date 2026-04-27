package controllers

import (
	"TestProject/config"
	"TestProject/helpers"
	"TestProject/models"
	"context"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetNotes(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// taken user id from parameters
	userId, exists := c.Get("userID")

	if !exists {
		helpers.ErrorResponse(c, 500, "Failed to get user ID from context")
		return
	}
	// convert string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userId.(string))
	if err != nil {
		helpers.ErrorResponse(c, 400, "Invalid user ID")
		return
	}

	notesCollection := config.GetCollection("notes")
	// find all notes for the user
	res, err := notesCollection.Find(ctx, gin.H{"user_id": objectID})
	if err != nil {
		helpers.ErrorResponse(c, 500, "Failed to fetch notes")
		return
	}

	var notes []models.Note
	if err := res.All(ctx, &notes); err != nil {
		helpers.ErrorResponse(c, 500, "Failed to decode notes")
		return
	}

	helpers.SuccessResponse(c, 200, "Notes fetched successfully", notes)
}

func CreateNote(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// taken user id from query parameters
	userId, exists := c.Get("userID")

	if !exists {
		helpers.ErrorResponse(c, 500, "Failed to get user ID from context")
		return
	}

	// convert string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userId.(string))
	if err != nil {
		helpers.ErrorResponse(c, 400, "Invalid user ID")
		return
	}

	// check if user exists
	userCollection := config.GetCollection("user")
	res := userCollection.FindOne(ctx, gin.H{"_id": objectID})
	if res.Err() != nil {
		helpers.ErrorResponse(c, 404, "User not found")
		return
	}

	var NotesRequestBody models.CreateNoteRequestBody
	if err := c.BindJSON(&NotesRequestBody); err != nil {
		helpers.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	// validate request body
	if err := validate.Struct(NotesRequestBody); err != nil {
		helpers.ErrorResponse(c, 400, err.Error())
		return
	}

	notesCollection := config.GetCollection("notes")
	// create a new note for the user
	_, err = notesCollection.InsertOne(ctx, gin.H{
		"user_id":    objectID,
		"title":      NotesRequestBody.Title,
		"content":    NotesRequestBody.Content,
		"created_at": time.Now(),
		"updated_at": time.Now(),
	})
	if err != nil {
		helpers.ErrorResponse(c, 500, "Failed to create note")
		return
	}
	helpers.SuccessResponse(c, 201, "Note created successfully!", nil)
}

func DeleteNote(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// taken user id and note id from query parameters
	note_id := c.Param("note_id")
	userId, exists := c.Get("userID")

	if !exists {
		helpers.ErrorResponse(c, 500, "Failed to get user ID from context")
		return
	}

	// convert string to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userId.(string))
	if err != nil {
		helpers.ErrorResponse(c, 400, "Invalid user ID")
		return
	}

	// check if user exists
	userCollection := config.GetCollection("user")
	if err := userCollection.FindOne(ctx, gin.H{"_id": userObjectID}).Err(); err != nil {
		helpers.ErrorResponse(c, 404, "User not found")
		return
	}

	noteObjectID, err := primitive.ObjectIDFromHex(note_id)
	if err != nil {
		helpers.ErrorResponse(c, 400, "Invalid note ID")
		return
	}
	// delete the note for the user
	notesCollection := config.GetCollection("notes")
	res, err := notesCollection.DeleteOne(ctx, gin.H{"_id": noteObjectID, "user_id": userObjectID})
	if err != nil {
		helpers.ErrorResponse(c, 500, "Failed to delete note")
		return
	}
	slog.Info("Delete result", "deletedCount", res.DeletedCount)
	if res.DeletedCount == 0 {
		helpers.ErrorResponse(c, 404, "Note not found for the user")
		return
	}
	helpers.SuccessResponse(c, 200, "Note deleted successfully!", nil)
}

func UpdateNote(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// taken user id and note id from query parameters
	note_id := c.Param("note_id")
	userId, exists := c.Get("userID")

	if !exists {
		helpers.ErrorResponse(c, 500, "Failed to get user ID from context")
		return
	}

	// convert string to ObjectID
	noteObjectID, err := primitive.ObjectIDFromHex(note_id)
	if err != nil {
		helpers.ErrorResponse(c, 400, "Invalid note ID")
		return
	}
	userObjectID, err := primitive.ObjectIDFromHex(userId.(string))
	if err != nil {
		helpers.ErrorResponse(c, 400, "Invalid user ID")
		return
	}
	// check if user exists
	userCollection := config.GetCollection("user")
	if err := userCollection.FindOne(ctx, gin.H{"_id": userObjectID}).Err(); err != nil {
		helpers.ErrorResponse(c, 404, "User not found")
		return
	}
	var NotesRequestBody models.UpdateNoteRequestBody
	if err := c.BindJSON(&NotesRequestBody); err != nil {
		helpers.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	// validate request body
	if err := validate.Struct(NotesRequestBody); err != nil {
		helpers.ErrorResponse(c, 400, err.Error())
		return
	}

	slog.Info("Updating note", "noteID", note_id, "userID", userId.(string), "title", NotesRequestBody.Title, "content", NotesRequestBody.Content)

	notesCollection := config.GetCollection("notes")
	res, err := notesCollection.UpdateOne(ctx, gin.H{"_id": noteObjectID, "user_id": userObjectID}, gin.H{
		"$set": gin.H{
			"title":      NotesRequestBody.Title,
			"content":    NotesRequestBody.Content,
			"updated_at": time.Now(),
		},
	})
	if err != nil {
		helpers.ErrorResponse(c, 500, "Failed to update note")
		return
	}

	slog.Info("Update result", "matchedCount", res.MatchedCount, "modifiedCount", res.ModifiedCount)

	if res.MatchedCount == 0 {
		helpers.ErrorResponse(c, 404, "Note not found for the user")
		return
	}
	helpers.SuccessResponse(c, 200, "Note updated successfully!", nil)
}
