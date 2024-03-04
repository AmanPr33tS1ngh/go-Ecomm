package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AmanPr33tS1ngh/go-Ecomm.git/database"
	"github.com/AmanPr33tS1ngh/go-Ecomm.git/models"
	generate "github.com/AmanPr33tS1ngh/go-Ecomm.git/tokens"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")
var Validate = validator.New()


func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil{
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(givenPassword))
	valid := true
	msg := ""
	if err != nil {
		log.Panic(err)
		msg = "Login or password incorrect"
		valid = false
	}
	return valid, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validateErr := Validate.Struct(user)
		if validateErr != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr})
		}

		// check if same email already exists or not
		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		// email count > 0
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists!"})
		}


		// check if same mobile already exists or not
		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()

		if err != nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		// phone count > 0
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Number already in use!"})
		}
		password := HashPassword(*user.Password)
		user.Password = &password

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserId = user.ID.Hex()
		token, refreshtoken, _ := generate.TokenGenerator(*user.Email, *user.FirstName, *user.LastName, user.UserId)
		user.Token = &token
		user.RefreshToken = &refreshtoken
		user.UserCart = make([]models.ProductUser, 0)
		user.AddressDetails = make([]models.Address, 0)
		user.OrderStatus = make([]models.Order, 0)
		_, inserterr := UserCollection.InsertOne(ctx, user)
		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusCreated, "Successfully Signed Up!!")
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email or password wrong"})
			return
		}

		isPasswordInvalid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()

		if !isPasswordInvalid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}
		token, refreshtoken, _ := generate.TokenGenerator(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, foundUser.UserId)
		defer cancel()

		generate.UpdateAllTokens(token, refreshtoken, foundUser.UserId)
		c.JSON(http.StatusFound, foundUser)
	}
}

func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var products models.Product
		defer cancel()
		if err := c.BindJSON(&products); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		products.ProductId = primitive.NewObjectID()
		_, anyerr := ProductCollection.InsertOne(ctx, products)
		if anyerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, "Successfully added our Product Admin!!")
	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productList []models.Product
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Something went wrong please try again in some time!")
			return
		}
		err = cursor.All(ctx, &productList)
		if err != nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		if err := cursor.Err(); err != nil{
			log.Println(err)
			c.IndentedJSON(http.StatusNotFound, "Invalid!")
			return
		}
		defer cancel()

		c.IndentedJSON(http.StatusOK, productList)
	}
}

func SearchProductByAdmin() gin.HandlerFunc {
	
	return func(c *gin.Context) {
		var productList []models.Product
		queryParam := c.Query("name")


		if queryParam == ""{
			log.Println("No query provided!")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid search index"})
			c.Abort()
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		searchedQuery, err := ProductCollection.Find(ctx, bson.M{"productName": bson.M{"$regex": queryParam}})
		if err != nil{
			c.IndentedJSON(http.StatusNotFound, "Something went wrong while fetching data!")
			return
		}
		err = searchedQuery.All(ctx, &productList)
		if err != nil{
			log.Println(err)
			c.IndentedJSON(http.StatusNotFound, "Invalid!")
			return
		}
		defer searchedQuery.Close(ctx)

		if err := searchedQuery.Err(); err != nil{
			log.Println(err)
			c.IndentedJSON(http.StatusNotFound, "Invalid!")
			return
		}
		defer cancel()
		c.IndentedJSON(http.StatusOK, productList)
	}
}