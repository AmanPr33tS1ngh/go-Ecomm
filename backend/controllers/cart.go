package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/AmanPr33tS1ngh/go-Ecomm.git/database"
	"github.com/AmanPr33tS1ngh/go-Ecomm.git/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication (prodCollection, userCollection *mongo.Collection) *Application{
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc  {
	return func(c *gin.Context) {
		var idIsEmpty string = " id is empty"
		productQuery := c.Query("id")
		if productQuery == "" {
			msg := "Product" + idIsEmpty
			log.Print(msg)
			_ = c.AbortWithError(http.StatusBadRequest, errors.New(msg))
			return
		}

		userIdQuery := c.Query("userId")
		if userIdQuery == "" {
			msg := "User" + idIsEmpty
			log.Print(msg)
			_ = c.AbortWithError(http.StatusBadRequest, errors.New(msg))
			return
		}
		productId, err := primitive.ObjectIDFromHex(productQuery)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.AddToCart(ctx, app.prodCollection, app.userCollection, productId, productQuery)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(http.StatusOK, "Successfully added to cart")
	}
}

func (app *Application) RemoveItem()	gin.HandlerFunc  {
	return func(c *gin.Context) {
		var idIsEmpty string = " id is empty"
		productQuery := c.Query("id")
		if productQuery == "" {
			msg := "Product" + idIsEmpty
			log.Print(msg)
			_ = c.AbortWithError(http.StatusBadRequest, errors.New(msg))
			return
		}

		userIdQuery := c.Query("userId")
		if userIdQuery == "" {
			msg := "User" + idIsEmpty
			log.Print(msg)
			_ = c.AbortWithError(http.StatusBadRequest, errors.New(msg))
			return
		}
		productId, err := primitive.ObjectIDFromHex(productQuery)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, productId, productQuery)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(http.StatusOK, "Successfully removed item from cart")
	}
}
func GetItemFromCart()	gin.HandlerFunc  {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == ""{
			c.JSON(http.StatusNotFound, "Invalid search index")
			c.Abort()
			return
		}
		userId, _ := primitive.ObjectIDFromHex(user_id)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var filledCart models.User
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: userId}}).Decode(&filledCart)
		if err != nil{
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "Not Found!")
			return
		}
		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: userId}}}} // finds user with userId
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$userCart"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$userCart.price"}}}}}}
		pointerCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})
		if err != nil{
			log.Println(err)
		}
		var listing []bson.M
		if err := pointerCursor.All(ctx, &listing); err != nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		for _, json := range listing{
			c.IndentedJSON(http.StatusOK, json["total"])
			c.IndentedJSON(http.StatusOK, filledCart.UserCart)
		}
		ctx.Done()
	}
}

func (app *Application )BuyFromCart()	gin.HandlerFunc  {
	return func(c *gin.Context) {
		userIdQuery := c.Query("userId")
		if userIdQuery == "" {
			var msg string = "User id is empty"
			log.Panicln(msg)
			_ = c.AbortWithError(http.StatusBadRequest, errors.New(msg))
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := database.BuyItemFromCart(ctx, app.userCollection, userIdQuery)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(http.StatusOK, "Successfully placed order")
	}
}

func (app *Application)InstantBuy()	gin.HandlerFunc  {
	return func(c *gin.Context) {
		var idIsEmpty string = " id is empty"
		productQuery := c.Query("id")
		if productQuery == "" {
			msg := "Product" + idIsEmpty
			log.Print(msg)
			_ = c.AbortWithError(http.StatusBadRequest, errors.New(msg))
			return
		}

		userIdQuery := c.Query("userId")
		if userIdQuery == "" {
			msg := "User" + idIsEmpty
			log.Print(msg)
			_ = c.AbortWithError(http.StatusBadRequest, errors.New(msg))
			return
		}
		productId, err := primitive.ObjectIDFromHex(productQuery)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.InstantBuyer(ctx, app.prodCollection, app.userCollection, productId, productQuery)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(http.StatusOK, "Successfully placed order")
	}
}