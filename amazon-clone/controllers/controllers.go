package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/go-playground/validator/v10"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson" 
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.userrData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductrData(database.Client, "Products")
var Validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)


}

func verifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""

	if err != nil {
		msg = "Login or Password is incorrect"
		valid = false
	}
	return valid, msg

}

func SignUp() gin.HandlerFunc {

	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email}) 
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error":"user already exists"})
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone":user.Phone})

		defer cancel()
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error":err})
			return
		}

		if count>0 {
			c.JSON(http.StatusBadRequest, gin.H{"error":" this phone no. is already in use"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_At = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitve.newObjectID()
		user.User_ID = User.ID.Hex()
		token, refreshtoken, _ := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, *user.ID)
		user.Token = &token
		user.Refresh_Token = &refreshtoken
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)
		_, inserterr := userCollection.InsertOne(ctx, user)
		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"the user did not get created"})
			return
		}
		defer cancel()

		c.JSON(http.StatusCreated, "Successfully signed in!")
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":err})
			return
		}
		
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"login or password incorrect"})
			return
		}

		PasswordIsValid, msg := verifyPassword(*user.Password, *founduser.Password)

		defer cancel()

		if !PasswordIsValid {
			c.JSON{http.StatusInternalServerError, gin.H{"error": msg}}
			fmt.Println(msg)
			return
		}

		token, refreshToken, _ := generate.TokenGenerator(*founderuser.Email, *founduser.First_Name, *founderuser.Last_Name, founduser.User_ID)
		defer cancel()

		generate.UpdateAllTokens(token, refreshToken, founderuser.User_ID)

		c.JSON(http.StatusFound, founduser)
	

	}
}

func ProductViewerAdmin() gin.HandlerFunc {

}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		var productlist []models.Product 
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err :ProductCollection.Find(ctx, bson.D{{}})

		if err!=nil {
			c.IndentedJSON(http.StatusInternalServerError, "something went wrong, please try again")
			return
		}

		err = cursor.All(ctx, &productlist)

		iff err!= nil {
			log,Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	defer cursor.Close()

	if err := cursor.err(); err != nil {
		log.println(err)
		c.IndentedJSON(400, "invalid")
		return
	}
	defer cancel()
	c.IndentedJSON(200, productlist)
	}

}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context){
		var searchProducts []models.Product
		queryParam := c.Query("name")
		
		if queryParam == "" {
			log.Println("query is empty")
			c.Header("content-type", "application/json")
			c.JSON(hhtp.StatusNotFound, gin.H{"Error":"Invalid search index"})
			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		searchquerydb, err := prodCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex":queryParam}})

		if err!= nil{
			c.IndentedJSON(404, "something went wrong while fetching the data")
			return
		}

		searchquerydb.All(ctx, &searchproducts)
		if err := nil{
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}

		defer searchquerydb.Close(ctx)

		if err := searchquerydb.Err(); err := nil{
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}

		defer cancel()
		c.IndentedJSON(200)
		
	}

}
