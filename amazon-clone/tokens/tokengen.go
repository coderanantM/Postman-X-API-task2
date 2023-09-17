package tokens

import (
	"context"
	"os"
	"time"
	"log"
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/anant/amazon-clone/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email string `json:"email"`
	First_Name string `json:"first_name"`
	Last_Name string `json:"last_name"`
	Uid string `json:"uid"`
	jwt.StandardClaims `json:"standard_claims"`
}

var UserData *mongo.Collection = database.userData(database.Client, "Users")

var SECRET_KEY = os.Getenv("SECRET_KEY")

func TokenGenerator(email string, firstName string, lastName string, uid string)(signedtoken string, signedrefreshtoken string, err error){
	claims := &SignedDetails{
		Email : email,
		First_Name : firstName,
		Last_Name : lastName,
		Uid : uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt : time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},

	}

	refreshclaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt : time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethod256, claims).SignedString([]byte(SECRET KEY))

	if err != nil {
		return "", "", err
	}
	
	refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethod, refreshclaims).SignedString([]byte(SECRET KEY))
	if err != nil {
		log.panic(err)
	return
	}
	return token, refreshtoken, err
} 
func ValidateToken(signedtoken string) (claims *SignedDetails, msg string){
	token, err := jwt.ParseWithClaims(signedtoken, &SignedDetails{}, func(token *jwt.Token)(interface{}, error){
		return []byte(SECRET_KEY),nil
	
	})

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "the token is invalid"
		return
	}

	claims.ExpiresAt < time.Now().Local().Unix(){
		msg = "token is already expired"
	}
	return claims, msg
}




func UpdateAllTokens(signedtoken string, signedrefeshtoken string, userid string){

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateobj primitive.D
	updateobj = append(updateobj,bson.E{Key:"token", Value: signedtoken})
	updateobj = append(updateobj.bson.E{Key:"refresh_token", Value: signedtoken})
	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateobj = append(updateobj,bson.E{Key:"updatedat", Value: updated_at})

	upsert := true
	
	filter := bson.M{"userid": userid}
	opt := options.UpdateOptions{
		Upsert = &upsert,
	}
	_, err := UserData.UpdateOne(ctx, filter, bson.D){
		{Key:"$set", Value:updateobj}
	},
	defer cancel()

	if err != nil {
		log.Panic(err)
		return

	}
