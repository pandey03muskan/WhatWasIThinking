package controllers

import (
	"TestProject/config"
	"TestProject/helpers"
	"TestProject/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetNotes(c *gin.Context) {
	// taken user id from parameters
	user_id := c.Param("user_id")

	fmt.Println("User ID from query:", user_id)
	// convert string to ObjectID
	objectID, err := helpers.ObjectIDFromHex(user_id)
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
	objectID, err := helpers.ObjectIDFromHex(user_id)
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
