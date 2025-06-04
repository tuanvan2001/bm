package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var orderService *OrderService

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://admin:c21a781d850c9b4e69c4627c801200c0c1f052fdc0aa6fd0@157.245.48.12")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	orderService = NewOrderService(client, "orderdb", "orders")
}

func getOrders(c *gin.Context) {
	orders, err := orderService.FindAll(context.TODO())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, orders)
}

func createOrder(c *gin.Context) {
	var newOrder Order

	if err := c.BindJSON(&newOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdOrder, err := orderService.Create(context.TODO(), newOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, createdOrder)
}

func getOrderById(c *gin.Context) {
	id := c.Param("id")

	order, err := orderService.FindOne(context.TODO(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, order)
}

func updateOrder(c *gin.Context) {
	id := c.Param("id")
	var updatedOrder Order

	if err := c.BindJSON(&updatedOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := orderService.Update(context.TODO(), id, updatedOrder)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, order)
}

func deleteOrder(c *gin.Context) {
	id := c.Param("id")

	err := orderService.Delete(context.TODO(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order deleted"})
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.GET("/orders", getOrders)
	router.POST("/orders", createOrder)
	router.GET("/orders/:id", getOrderById)
	router.PUT("/orders/:id", updateOrder)
	router.DELETE("/orders/:id", deleteOrder)

	router.Run("localhost:8080")
}
