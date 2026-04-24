package controllers

import (
	"TestProject/config"
	"TestProject/helpers"
	"TestProject/models"
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context) {
	// by this we are telling to mongoDB complete the operation within 5 seconds otherwise it will be cancelled
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// take email and password from request body
	var requestBody models.RegisterUserRequestBody
	if err := c.BindJSON(&requestBody); err != nil { // bind request body to the requestBody struct
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	fmt.Println("email:", requestBody.Email)
	fmt.Println("password:", requestBody.Password)
	fmt.Println("name:", requestBody.Name)

	// validate the request body
	if requestBody.Email == "" || requestBody.Password == "" || requestBody.Name == "" {
		c.JSON(400, gin.H{"error": "Email, password and name are required"})
		return
	}

	// check if user already exists
	userCollection := config.GetCollection("user")
	res := userCollection.FindOne(ctx, gin.H{"email": requestBody.Email})
	if res.Err() == nil {
		c.JSON(400, gin.H{"error": "User already exists"})
		return
	}

	// hashing the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	// let's send this to DB
	// userCollection := config.GetCollection("user")
	_, error := userCollection.InsertOne(ctx, gin.H{
		"email":    requestBody.Email,
		"password": hashedPassword,
		"name":     requestBody.Name,
	})
	if error != nil {
		c.JSON(500, gin.H{"error": "Failed to add user"})
		return
	}

	c.JSON(200, gin.H{
		"status":  200,
		"message": "User added successfully",
	})
}

func LoginUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var requestBody models.LoginUserRequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// validate the request body
	if requestBody.Email == "" || requestBody.Password == "" {
		c.JSON(400, gin.H{"error": "Email and password are required"})
		return
	}

	userCollection := config.GetCollection("user")
	res := userCollection.FindOne(ctx, gin.H{"email": requestBody.Email})
	if res.Err() != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	var user models.LoginUser
	if err := res.Decode(&user); err != nil {
		c.JSON(500, gin.H{"error": "Failed to decode user data"})
		return
	}

	// compare the hashed password with the password from request body
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid password"})
		return
	}

	// generate JWT token
	token, err := helpers.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(200, gin.H{
		"status":  200,
		"message": "Login successful",
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
			"token": token,
		},
	})
}

func GetUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// take ID from request body
	// var requestBody models.GetUserRequestBody
	// if err := c.BindJSON(&requestBody); err != nil {
	// 	c.JSON(400, gin.H{"error": "Invalid request body"})
	// 	return
	// }
	// fmt.Println("ID:", requestBody.ID)

	userId, exists := c.Get("userID") // get userID from context set by auth middleware
	if !exists {
		c.JSON(500, gin.H{"error": "Failed to get user ID from context"})
		return
	}

	userCollection := config.GetCollection("user")              //get user collection from DB
	objectID, err := primitive.ObjectIDFromHex(userId.(string)) // convert string ID to ObjectID
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}
	fmt.Println("ObjectID:", objectID)
	res := userCollection.FindOne(ctx, gin.H{"_id": objectID})
	fmt.Println("res:", res.Err())
	if res.Err() != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	var user models.GetUserResponseBody
	if err := res.Decode(&user); err != nil {
		c.JSON(500, gin.H{"error": "Failed to decode user data"})
		return
	}

	c.JSON(200, gin.H{
		"status":  200,
		"message": "User fetched successfully!",
		"user":    user,
	})
}
