package main

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderService struct {
	collection *mongo.Collection
}

func NewOrderService(client *mongo.Client, dbName, collectionName string) *OrderService {
	collection := client.Database(dbName).Collection(collectionName)
	return &OrderService{collection: collection}
}

func (s *OrderService) Create(ctx context.Context, order Order) (Order, error) {
	order.ID = primitive.NewObjectID()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	_, err := s.collection.InsertOne(ctx, order)
	if err != nil {
		var writeException mongo.WriteException
		if errors.As(err, &writeException) {
			for _, e := range writeException.WriteErrors {
				if e.Code == 11000 { // Duplicate key error
					return Order{}, errors.New("Order with this order_code already exists.")
				}
			}
		}
		return Order{}, err
	}
	return order, nil
}

func (s *OrderService) FindAll(ctx context.Context) ([]Order, error) {
	var orders []Order
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *OrderService) FindOne(ctx context.Context, id string) (Order, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Order{}, errors.New("Invalid ID format")
	}

	var order Order
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Order{}, errors.New("Order not found")
		}
		return Order{}, err
	}
	return order, nil
}

func (s *OrderService) Update(ctx context.Context, id string, updatedOrder Order) (Order, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Order{}, errors.New("Invalid ID format")
	}

	updatedOrder.UpdatedAt = time.Now()

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": updatedOrder}

	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return Order{}, err
	}

	if result.MatchedCount == 0 {
		return Order{}, errors.New("Order not found")
	}

	// Retrieve the updated document to return it
	var order Order
	err = s.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		return Order{}, err
	}

	return order, nil
}

func (s *OrderService) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("Invalid ID format")
	}

	filter := bson.M{"_id": objID}
	result, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("Order not found")
	}

	return nil
}
