package controllers

import (
	"TestProject/config"
	"TestProject/helpers"
	"TestProject/models"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func RegisterUser(c *gin.Context) {
	// by this we are telling to mongoDB complete the operation within 5 seconds otherwise it will be cancelled
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// take email and password from request body
	var requestBody models.RegisterUserRequestBody
	if err := c.BindJSON(&requestBody); err != nil { // bind request body to the requestBody struct
		helpers.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	// validate the request body
	if err := validate.Struct(requestBody); err != nil {
		helpers.ErrorResponse(c, 400, err.Error())
		return
	}

	// check if user already exists
	userCollection := config.GetCollection("user")
	res := userCollection.FindOne(ctx, gin.H{"email": requestBody.Email})
	if res.Err() == nil {
		helpers.ErrorResponse(c, 400, "User already exists")
		return
	}

	// hashing the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
	if err != nil {
		helpers.ErrorResponse(c, 500, "Failed to hash password")
		return
	}

	// let's send this to DB
	_, error := userCollection.InsertOne(ctx, gin.H{
		"email":    requestBody.Email,
		"password": hashedPassword,
		"name":     requestBody.Name,
	})
	if error != nil {
		helpers.ErrorResponse(c, 500, "Failed to add user")
		return
	}

	helpers.SuccessResponse(c, 201, "User added successfully", nil)
}

func LoginUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var requestBody models.LoginUserRequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		helpers.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	// validate the request body
	if err := validate.Struct(requestBody); err != nil {
		helpers.ErrorResponse(c, 400, err.Error())
		return
	}

	userCollection := config.GetCollection("user")
	res := userCollection.FindOne(ctx, gin.H{"email": requestBody.Email})
	if res.Err() != nil {
		helpers.ErrorResponse(c, 404, "User not found")
		return
	}

	var user models.LoginUser
	if err := res.Decode(&user); err != nil {
		helpers.ErrorResponse(c, 500, "Failed to decode user data")
		return
	}

	// compare the hashed password with the password from request body
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password)); err != nil {
		helpers.ErrorResponse(c, 401, "Invalid password")
		return
	}

	// generate JWT token
	token, err := helpers.GenerateJWT(user.ID)
	if err != nil {
		helpers.ErrorResponse(c, 500, "Failed to generate token")
		return
	}

	helpers.SuccessResponse(c, 200, "Login successful", gin.H{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"token": token,
	})
}

func GetUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userId, exists := c.Get("userID") // get userID from context set by auth middleware
	if !exists {
		helpers.ErrorResponse(c, 500, "Failed to get user ID from context")
		return
	}

	userCollection := config.GetCollection("user")              //get user collection from DB
	objectID, err := primitive.ObjectIDFromHex(userId.(string)) // convert string ID to ObjectID
	if err != nil {
		helpers.ErrorResponse(c, 400, "Invalid ID")
		return
	}
	res := userCollection.FindOne(ctx, gin.H{"_id": objectID})
	if res.Err() != nil {
		helpers.ErrorResponse(c, 404, "User not found")
		return
	}

	var user models.GetUserResponseBody
	if err := res.Decode(&user); err != nil {
		helpers.ErrorResponse(c, 500, "Failed to decode user data")
		return
	}
	helpers.SuccessResponse(c, 200, "User fetched successfully", user)
}
