package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddListToQueue(client *mongo.Client, listID primitive.ObjectID) error {
	// get all leads in the list and add them to the queue only if they are not already in the queue
	collection := client.Database("email_verify").Collection("leads")
	// add leads concurrently to make it faster
	cursor, err := collection.Find(context.TODO(), bson.M{"list_id": listID})
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())

	queueCollection := client.Database("email_verify").Collection("verification_queue")
	var queue []VerificationQueue
	for cursor.Next(context.Background()) {
		var lead Lead
		cursor.Decode(&lead)
		queue = append(queue, VerificationQueue{Email: lead.Email, LeadID: lead.ID, ListID: lead.ListID})
	}

	// insert the leads into the queue
	var queueDocuments []interface{}
	for _, q := range queue {
		queueDocuments = append(queueDocuments, q)
	}
	_, err = queueCollection.InsertMany(context.TODO(), queueDocuments)
	return err
}

func IsListInQueue(client *mongo.Client, listID primitive.ObjectID) (bool, error) {
	collection := client.Database("email_verify").Collection("verification_queue")
	count, err := collection.CountDocuments(context.TODO(), bson.M{"list_id": listID})
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func RemoveListFromQueue(client *mongo.Client, listID primitive.ObjectID) error {
	// remove all leads in the list from the queue
	queueCollection := client.Database("email_verify").Collection("verification_queue")
	_, err := queueCollection.DeleteMany(context.TODO(), bson.M{"lead_id": bson.M{"$in": listID}})
	return err
}

func GetQueue(client *mongo.Client) ([]VerificationQueue, error) {
	collection := client.Database("email_verify").Collection("verification_queue")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var queue []VerificationQueue
	for cursor.Next(context.Background()) {
		var q VerificationQueue
		cursor.Decode(&q)
		queue = append(queue, q)
	}
	return queue, nil
}

func GetQueueCount(client *mongo.Client) (int64, error) {
	collection := client.Database("email_verify").Collection("verification_queue")
	return collection.CountDocuments(context.TODO(), bson.M{})
}

func Dequeue(client *mongo.Client, queueItemId primitive.ObjectID, EmailIsValid string, VerificationResult bson.M) error {
	collection := client.Database("email_verify").Collection("verification_queue")
	queueItem := VerificationQueue{}
	err := collection.FindOne(context.TODO(), bson.M{"_id": queueItemId}).Decode(&queueItem)
	if err != nil {
		return err
	}

	// update the lead
	leadsCollection := client.Database("email_verify").Collection("leads")
	_, err = leadsCollection.UpdateOne(context.TODO(), bson.M{"_id": queueItem.LeadID}, bson.M{"$set": bson.M{"email_is_valid": EmailIsValid, "verification_result": VerificationResult, "email_verified": true}})
	if err != nil {
		return err
	}

	// remove the queue item
	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": queueItemId})
	return err
}
