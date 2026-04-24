package controllers

import (
	"TestProject/config"
	"TestProject/models"
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetNotes(c *gin.Context) {
	// taken user id from parameters
	user_id := c.Param("user_id")

	fmt.Println("User ID from query:", user_id)
	// convert string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(user_id)
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

	var notes []models.GetNotes
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
	user_id := c.Param("user_id")

	// convert string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(user_id)
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
		"user_id": objectID,
		"title":   NotesRequestBody.Title,
		"content": NotesRequestBody.Content,
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
	user_id := c.Param("user_id")
	note_id := c.Param("note_id")

	// convert string to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	// check if user exists
	userCollection := config.GetCollection("user")
	response := userCollection.FindOne(ctx, gin.H{"_id": userObjectID})
	if response == nil {
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
	user_id := c.Param("user_id")
	// convert string to ObjectID
	fmt.Println("userid noteis :", user_id, note_id)
	userObjectID, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}
	fmt.Println("user object ID:", userObjectID)
	// check if user exists
	userCollection := config.GetCollection("user")
	response := userCollection.FindOne(ctx, gin.H{"_id": userObjectID})
	if response == nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	noteObjectID, err := primitive.ObjectIDFromHex(note_id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid note ID"})
		return
	}
	fmt.Println("note object ID:", noteObjectID)
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
			"title":   NotesRequestBody.Title,
			"content": NotesRequestBody.Content,
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
