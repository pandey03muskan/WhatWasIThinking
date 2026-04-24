package controllers

import (
	"TestProject/config"
	"TestProject/models"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetNotes(c *gin.Context) {
	// taken user id from parameters
	userId, exists := c.Get("userID")

	if !exists {
		c.JSON(500, gin.H{"error": "Failed to get user ID from context"})
		return
	}
	// convert string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userId.(string))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	notesCollection := config.GetCollection("notes")
	// find all notes for the user
	res, err := notesCollection.Find(c, gin.H{"user_id": objectID})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch notes"})
		return
	}

	var notes []models.Note
	if err := res.All(c, &notes); err != nil {
		c.JSON(500, gin.H{"error": "Failed to decode notes"})
		return
	}

	c.JSON(200, gin.H{
		"status":  200,
		"message": "Notes fetched successfully!",
		"data":    notes,
	})
}

func CreateNote(c *gin.Context) {
	// taken user id from query parameters
	userId, exists := c.Get("userID")

	if !exists {
		c.JSON(500, gin.H{"error": "Failed to get user ID from context"})
		return
	}

	// convert string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userId.(string))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	// check if user exists
	userCollection := config.GetCollection("user")
	res := userCollection.FindOne(c, gin.H{"_id": objectID})
	if res == nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	var NotesRequestBody models.CreateNoteRequestBody
	if err := c.BindJSON(&NotesRequestBody); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// validate request body
	if NotesRequestBody.Title == "" || NotesRequestBody.Content == "" {
		c.JSON(400, gin.H{"error": "Title and content are required"})
		return
	}

	notesCollection := config.GetCollection("notes")
	// create a new note for the user
	_, err = notesCollection.InsertOne(c, gin.H{
		"user_id":    objectID,
		"title":      NotesRequestBody.Title,
		"content":    NotesRequestBody.Content,
		"created_at": time.Now(),
		"updated_at": time.Now(),
	})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create note"})
		return
	}
	c.JSON(200, gin.H{
		"status":  200,
		"message": "Note created successfully!",
	})
}

func DeleteNote(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// taken user id and note id from query parameters
	note_id := c.Param("note_id")
	userId, exists := c.Get("userID")

	if !exists {
		c.JSON(500, gin.H{"error": "Failed to get user ID from context"})
		return
	}

	// convert string to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userId.(string))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	// check if user exists
	userCollection := config.GetCollection("user")
	if err := userCollection.FindOne(ctx, gin.H{"_id": userObjectID}).Err(); err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	noteObjectID, err := primitive.ObjectIDFromHex(note_id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid note ID"})
		return
	}
	// delete the note for the user
	notesCollection := config.GetCollection("notes")
	res, err := notesCollection.DeleteOne(ctx, gin.H{"_id": noteObjectID, "user_id": userObjectID})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete note"})
		return
	}
	if res.DeletedCount == 0 {
		c.JSON(404, gin.H{"error": "Note not found"})
		return
	}
	c.JSON(200, gin.H{
		"status":  200,
		"message": "Note deleted successfully!",
	})
}

func UpdateNote(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// taken user id and note id from query parameters
	note_id := c.Param("note_id")
	userId, exists := c.Get("userID")

	if !exists {
		c.JSON(500, gin.H{"error": "Failed to get user ID from context"})
		return
	}

	// convert string to ObjectID
	noteObjectID, err := primitive.ObjectIDFromHex(note_id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid note ID"})
		return
	}
	userObjectID, err := primitive.ObjectIDFromHex(userId.(string))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}
	// check if user exists
	userCollection := config.GetCollection("user")
	if err := userCollection.FindOne(ctx, gin.H{"_id": userObjectID}).Err(); err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	var NotesRequestBody models.UpdateNoteRequestBody
	if err := c.BindJSON(&NotesRequestBody); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	// validate request body
	if NotesRequestBody.Title == "" || NotesRequestBody.Content == "" {
		c.JSON(400, gin.H{"error": "Title and content are required"})
		return
	}
	// update the note for the user
	notesCollection := config.GetCollection("notes")
	res, err := notesCollection.UpdateOne(ctx, gin.H{"_id": noteObjectID, "user_id": userObjectID}, gin.H{
		"$set": gin.H{
			"title":      NotesRequestBody.Title,
			"content":    NotesRequestBody.Content,
			"updated_at": time.Now(),
		},
	})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update note"})
		return
	}
	if res.MatchedCount == 0 {
		c.JSON(404, gin.H{"error": "Note not found"})
		return
	}
	c.JSON(200, gin.H{
		"status":  200,
		"message": "Note updated successfully!",
	})
}
