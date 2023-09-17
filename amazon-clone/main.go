package main

import(
	"github.com/anant/amazon-clone/controllers"
	"github.com/anant/amazon-clone/database"
	"github.com/anant/amazon-clone/middleware"
	"github.com/anant/amazon-clone/routes"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == ""{
		port = "8000"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.instantBuy())

	log.Fatal(router.Run(":" + port))

}