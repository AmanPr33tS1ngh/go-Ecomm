package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/AmanPr33tS1ngh/go-Ecomm.git/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id != "" {
			fmt.Println("user id not provided")
			c.JSON(http.StatusNotFound, gin.H{"error": "User id nont provided"})
			c.Abort()
			return
		}
	
		address, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal server error")
		}
		var addresses models.Address
		addresses.AddressId = primitive.NewObjectID()
		if err = c.BindJSON(&addresses); err != nil{
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5 *time.Second)
		defer cancel()

		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}} // finds user with userId
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}
		pointerCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, "Internal server error!")
		}
		var addressInfo []bson.M
		if err := pointerCursor.All(ctx, &addressInfo); err != nil{
			panic(err)
		}
		var size int32
		for _, address_no := range addressInfo{
			count := address_no["count"]
			size = count.(int32)
		}
		if size >= 2 {
			c.IndentedJSON(http.StatusNotFound, "Not Allowed")
			return
		}
		filter := bson.D{primitive.E{Key: "_id", Value: address}}
		update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		defer cancel()
		ctx.Done()
	}
}
func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {

		user_id := c.Query("id")
		if user_id != "" {
			fmt.Println("user id not provided")
			c.JSON(http.StatusNotFound, gin.H{"error": "User id nont provided"})
			c.Abort()
			return
		}
		userId, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal server error")
		}

		var editaddress models.Address

		if err = c.BindJSON(&editaddress); err != nil{
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5 *time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: userId}}
		// address.0 referes to home address
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: editaddress.House}, 
			{Key: "address.0.street_name", Value: editaddress.Street}, {Key: "address.0.city_name", Value: editaddress.City}, 
			{Key: "address.0.pin_code", Value: editaddress.PinCode}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil{
			fmt.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "Something went wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(http.StatusOK, "Successfully updated the home address!")
	}
}
func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {

		user_id := c.Query("id")
		if user_id != "" {
			fmt.Println("user id not provided")
			c.JSON(http.StatusNotFound, gin.H{"error": "User id nont provided"})
			c.Abort()
			return
		}
		userId, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal server error")
		}

		var editaddress models.Address

		if err = c.BindJSON(&editaddress); err != nil{
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5 *time.Second)
		defer cancel()

		// address.1 referes to work address
		filter := bson.D{primitive.E{Key: "_id", Value: userId}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_name", Value: editaddress.House}, 
			{Key: "address.1.street_name", Value: editaddress.Street}, {Key: "address.1.city_name", Value: editaddress.City}, 
			{Key: "address.1.pin_code", Value: editaddress.PinCode}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil{
			fmt.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "Something went wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(http.StatusOK, "Successfully updated the work address!")
	}
}
func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == ""{
			c.JSON(http.StatusNotFound, "Invalid search index")
			c.Abort()
			return
		}
		addresses := make([]models.Address, 0)
		userId, err := primitive.ObjectIDFromHex(user_id)
		if err != nil{
			fmt.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "Internal server error")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key:"_id", Value: userId}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key:"address", Value: addresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil{
			fmt.Println(err)
			c.IndentedJSON(http.StatusNotFound, "Wrong command")
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(http.StatusOK, "Successfully Deleted!")
	}
}
