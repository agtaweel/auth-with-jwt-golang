package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/agtaweel/golang-jwt-project/database"
	"github.com/agtaweel/golang-jwt-project/helpers"
	"github.com/agtaweel/golang-jwt-project/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

var productCollection *mongo.Collection = database.OpenCollection(database.Client, "product")

func GetUsers() gin.HandlerFunc {

	return func(c *gin.Context) {

		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		// startIndex := (page - 1) * recordPerPage
		startIndex, err := strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage,
		})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while listing users"})
			return
		}

		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil {
			log.Fatal(err)
			log.Panic(startIndex)
		}
		c.JSON(http.StatusOK, allUsers[0])
	}

}
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func addProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var product models.Product
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		validationErr := validate.Struct(product)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		defer cancel()
		Created_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		product.Created_at = &Created_at
		Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		product.Updated_at = &Updated_at
		product.ID = primitive.NewObjectID()
		insertionId, insertErr := productCollection.InsertOne(ctx, product)

		if insertErr != nil {
			msg := "Product was not inserted"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, insertionId)

	}
}
