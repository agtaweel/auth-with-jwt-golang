package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/agtaweel/golang-jwt-project/helpers"
	"github.com/agtaweel/golang-jwt-project/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
)

func AddProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var product models.Product
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
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
		product.Product_id = product.ID.Hex()
		var num = helpers.ToFixed(float64(*product.Price), 2)
		product.Price = &num
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

func IndexProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var products []bson.M
		result, err1 := productCollection.Find(ctx, bson.M{})
		defer cancel()
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
			return
		}
		if err := result.All(ctx, &products); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, products)

	}
}

func DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		productID := c.Param("product_id")
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		result, err := productCollection.DeleteOne(ctx, bson.M{"product_id": productID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cancel()
		if result.DeletedCount == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Not found"})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func GetProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("product_id")
		var product models.Product
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		err := productCollection.FindOne(ctx, bson.M{"product_id": productID}).Decode(&product)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, product)
	}
}

func UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		productID := c.Param("product_id")
		var product models.Product
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer cancel()
		var num = helpers.ToFixed(float64(*product.Price), 2)
		Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		filter := bson.M{"product_id": productID}
		update := bson.M{
			"$set": bson.M{
				"name":        product.Name,
				"buyer":       product.Buyer,
				"price":       &num,
				"phone":       product.Phone,
				"description": product.Description,
				"updated_at":  &Updated_at,
			},
		}
		result, err := productCollection.UpdateOne(
			ctx,
			filter,
			update,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}
