package controllers

import (
	"context"
	"net/http"
	"pokemon-api/configs"
	"pokemon-api/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		defer cancel()

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"data": err.Error()})
			return
		}

		if validationErr := validate.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"data": validationErr.Error()})
			return
		}

		newUser := models.User{
			Id:    primitive.NewObjectID(),
			Name:  user.Name,
			Score: user.Score,
		}

		result, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"data": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"data": result})
	}
}

func GetAUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		var user models.User
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"data": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": user})
	}
}

func EditAUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		var user models.User
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(userId)

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"data": err.Error()})
			return
		}

		if validationErr := validate.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"data": validationErr.Error()})
			return
		}

		update := bson.M{"name": user.Name, "score": user.Score}
		result, err := userCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"data": err.Error()})
			return
		}

		var updatedUser models.User
		if result.MatchedCount == 1 {
			err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"data": err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"data": updatedUser})
	}
}

func DeleteAUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"data": err.Error()})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, gin.H{"data": "User with specified ID not found!"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": "User successfully deleted!"})
	}
}

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var users []models.User
		defer cancel()

		results, err := userCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"data": err.Error()})
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser models.User
			if err = results.Decode(&singleUser); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"data": err.Error()})
			}

			users = append(users, singleUser)
		}

		c.JSON(http.StatusOK, gin.H{"data": users})

	}
}

func GetAllUsersTop() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var users []models.User
		defer cancel()

		findOptions := options.Find()
		findOptions.SetLimit(10)
		findOptions.SetSort(bson.D{{Key: "score", Value: -1}})
		results, err := userCollection.Find(ctx, bson.D{}, findOptions)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"data": err.Error()})
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser models.User
			if err = results.Decode(&singleUser); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"data": err.Error()})
			}

			users = append(users, singleUser)
		}

		c.JSON(http.StatusOK, gin.H{"data": users})

	}
}
