package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateList(client *mongo.Client, name string) (primitive.ObjectID, error) {
	collection := client.Database("email_verify").Collection("lists")
	res, err := collection.InsertOne(context.TODO(), bson.M{"name": name})
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func GetList(client *mongo.Client, id primitive.ObjectID) (List, error) {
	collection := client.Database("email_verify").Collection("lists")
	var list List
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&list)
	if err != nil {
		return List{}, err
	}
	return list, nil
}

func GetLists(client *mongo.Client) ([]List, error) {
	collection := client.Database("email_verify").Collection("lists")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var lists []List
	for cursor.Next(context.Background()) {
		var list List
		cursor.Decode(&list)
		lists = append(lists, list)
	}
	return lists, nil
}

func DeleteList(client *mongo.Client, id primitive.ObjectID) error {
	collection := client.Database("email_verify").Collection("lists")
	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}
