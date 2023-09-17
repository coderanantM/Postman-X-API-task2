package controllers

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/anant/amazon-clone/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (

	ErrorCantFindProduct = errors.New("cannot find product")
	ErrorCantDecodeproducts = errors.New("cannot find product")
	ErroruserIDIsNotValid = errors.New("this user is not valid")
	ErrorCantUpdateUSer = errors.New("cannot add this product to the cart")
	ErrorCantRemoveItemCart = errors.New("cannot remove this product from the cart")
	ErrorCannotGetItemCart = errors.New("unable to get this product from the cart")
	ErrorCannotBuyCartItem = errors.New("cannot update this purchase")

)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	searchfromdb, err := prodCollection.Find(ctx, bson.M{"_id": productID})
	if err != nil {
		log.Println(err)
		return ErrorCantFindProduct
	}
	
	var productCart []models.ProductUser
	err = searchfromdb.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrorCantDecodeproducts
	}

	

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErroruserIDIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key:"$push", Value: bson.D{primitive.E{Key:"usercart", Value: bson.D{{Key:"$each", Value: productCart}}}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err!=nil {
		return ErrorCantUpdateUSer
	}
	return nil
}

func RemoveCartItem(ctx, context.Context, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error{
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErruserIDIsNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value : id}}
	update := bson.M{"$pull":bson.M{"usercart":bson.M{"_id":productID}}}
	_, err = UpdateMany(ctx, filter, update)
	if err != nil{
		return ErrorCantRemoveItemCart
	}
	return nil


}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string)error{
	id, err := primitive.ObjectIDFromHex()
	if err != nil{
		return ErrorUserIDIsNotValid
	}

	var getcartitems models.User
	var ordercart models.Order

	ordercart.Order_ID = primitive.NewObjectID()
	ordercart.Ordered_At = time.Now()
	ordercart.Order_Cart = make([]models.ProductUser, 0)
	ordercart.Payment_Method.COD = true

	unwind := bson.D{{Key"$unwind", Value:bson.D{primitive.E{Key:"path", Value"$usercart"}}}}
	grouping := bson.D{{Key"$group", Value:bson.D{primitive.E{Key:"_id", Value:"$_id", {Key:"total", Value:"total", Value: bson.D{primitive.E{Key:"$sum", Value:"$usercart.price"}}}}}}}
	currentresults, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	ctx.Done()
	if err != nil {
		panic(err)
	}

	var getusercart []bson.M
	if err = currentresults.All(ctx, &getusercart); err != nil {
		panic(err)
	}
	var total_price int32ctx

	for _, user_item := range getusercart {
		price := user_item["total"]
		total_price = price.(int32)
	}
	ordercart.Price = int(total_price)

	filter := bson.D{primitive.E{Key:"_id", Value:id}}
	update := bson.D{{Key:"$push", Value:bson.D{primitive.E{Key:"orders", Value:ordercart}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	userCollection>findOne(ctx, bson.D{primitive.E{Key:"_id", Value:id}}).Decode(&getcartitems)
	if err != nil{
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key:"_id", Value:id}}
	update2 := bson.M{"$push":bson.M{"orders.$[].order_list":bson.M{"$each":getcartitems.UserCart}}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil{
		log.Println(err)
	}

	usercart_empty := make([]models.productUser, 0)
	filter3 := bson.D{primitive.E{Key:"_id", Value: id}}
	update3 := bson.D{{Key:"$set", value:bson.D{primitive.E{Key:"usercart", Value:usercart_empty}}}}
	_, err = userCollection.UpdateOne(ctx, filter3, update3)
	if err != nil{
		return ErrorCannotBuyCartItem
	}
	return nil



	
}
func InstantBuyer(ctx context.Context, prodCollection, userCollection, productID primitive.ObjectID, UserID string) error {
	id, err := primitive.ObjectIDFromHex(UserID)

	if err!=nil{
		log.Println(err)
		return ErroruserIDIsNotValid
	}

	var product_details models.ProductUser
	var order_detail models.Order

	order_detail.Order_ID = primitive.NewObjectID()
	order_details.Ordered_At = time.Now()
	order_detail.Order_Cart = make([]models.ProductUser, 0)
	order_detail.Payment_Method.COD = true
	prodCollection.FindOne(ctx, bson.D{primitive.E{Key:"_id", Value: productID}}).Decode(&product_details)
	if err!=nil{
		log.Println(err)
	}
	order_detail.Price = product_details.Price

	filter := bson.D{primitive.E{Key:"_id", Value:id}}
	update := bson.D{{Key:"$push", Value:bson.D{primitive.E{Key:"orders", Value:order_detail}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err!=nil{
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key:"_id", Value:id}}
	update2 := bson.M{"$push":bson.M{"orders.$[].order_list":product_details}}

	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err!=nil{
		log.Println(err)
	}
	return nil

}