package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateLead(client *mongo.Client, email string, listID primitive.ObjectID, leadData bson.M) (primitive.ObjectID, error) {
	collection := client.Database("email_verify").Collection("leads")
	res, err := collection.InsertOne(context.TODO(), bson.M{"email": email, "list_id": listID, "lead_data": leadData})
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func GetLead(client *mongo.Client, id primitive.ObjectID) (Lead, error) {
	collection := client.Database("email_verify").Collection("leads")
	var lead Lead
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&lead)
	if err != nil {
		return Lead{}, err
	}
	return lead, nil
}

func GetLeads(client *mongo.Client, listID primitive.ObjectID) ([]Lead, error) {
	collection := client.Database("email_verify").Collection("leads")
	cursor, err := collection.Find(context.TODO(), bson.M{"list_id": listID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var leads []Lead
	for cursor.Next(context.Background()) {
		var lead Lead
		cursor.Decode(&lead)
		leads = append(leads, lead)
	}
	return leads, nil
}

func GetLeadsCount(client *mongo.Client, listID primitive.ObjectID) (int64, error) {
	collection := client.Database("email_verify").Collection("leads")
	return collection.CountDocuments(context.TODO(), bson.M{"list_id": listID})

}

func CountEmailVerified(client *mongo.Client, listID primitive.ObjectID) (int64, error) {
	collection := client.Database("email_verify").Collection("leads")
	return collection.CountDocuments(context.TODO(), bson.M{"list_id": listID, "email_verified": true})
}

func DeleteLead(client *mongo.Client, id primitive.ObjectID) error {
	collection := client.Database("email_verify").Collection("leads")
	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}

func AddLeadsFromCSV(client *mongo.Client, listID primitive.ObjectID, request *http.Request) error {
	collection := client.Database("email_verify").Collection("leads")

	// Read csv from request
	file, _, err := request.FormFile("csvfile")
	if err != nil {
		return err
	}
	defer file.Close()

	// Parse csv
	csvReader := csv.NewReader(file)
	// find email index in csv
	header, err := csvReader.Read()
	if err != nil {
		return err
	}
	emailIdx := -1
	for i, h := range header {
		if strings.ToLower(h) == "email" {
			emailIdx = i
			break
		}
	}
	// add leads concurrently to make it faster make an array of leads and insert them all at once

	var leads []interface{}
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Get email from csv
		email := record[emailIdx]

		// every other field is lead data

		leadData := make(map[string]interface{})

		for i, h := range header {
			if i == emailIdx {
				continue
			}
			leadData[h] = record[i]
		}

		// Create Lead object
		lead := Lead{
			ID:       primitive.NewObjectID(),
			Email:    email,
			ListID:   listID,
			LeadData: leadData,
		}

		// Add lead to leads array
		leads = append(leads, lead)
	}

	// Insert leads into database
	_, err = collection.InsertMany(context.TODO(), leads)

	return err
}

func DownloadLeadsAsCSV(client *mongo.Client, listID primitive.ObjectID, w http.ResponseWriter) error {
	collection := client.Database("email_verify").Collection("leads")
	cursor, err := collection.Find(context.TODO(), bson.M{"list_id": listID})
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=leads.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// write header from first lead

	if cursor.Next(context.Background()) {
		var lead Lead
		cursor.Decode(&lead)
		// lead.LeadData is a primitive.D
		// convert it to map[string]interface{}
		leadData, err := json.Marshal(lead.LeadData)
		if err != nil {
			return err
		}
		// define and array leadDataKeys to store keys of leadData sample leadDataKeys = [{}]
		var leadDataKeys []struct {
			Key   string
			Value string
		}
		// unmarshal leadData to leadDataKeys
		err = json.Unmarshal(leadData, &leadDataKeys)
		if err != nil {
			return err
		}
		// write header to csv
		var header []string
		header = append(header, "email")
		for _, key := range leadDataKeys {
			header = append(header, key.Key)
		}
		header = append(header, "email_is_valid")
		header = append(header, "verification_result")
		err = writer.Write(header)
		if err != nil {
			return err
		}
		cursor.Close(context.Background())
		cursor, err = collection.Find(context.TODO(), bson.M{"list_id": listID})
		if err != nil {
			return err
		}
		defer cursor.Close(context.TODO())
	}

	for cursor.Next(context.Background()) {
		var lead Lead
		cursor.Decode(&lead)
		// lead.LeadData is a primitive.D
		// convert it to map[string]interface{}
		leadData, err := json.Marshal(lead.LeadData)
		if err != nil {
			return err
		}
		// define and array leadDataKeys to store keys of leadData sample leadDataKeys = [{}]
		var leadDataKeys []struct {
			Key   string
			Value string
		}
		// unmarshal leadData to leadDataKeys
		err = json.Unmarshal(leadData, &leadDataKeys)
		if err != nil {
			return err
		}
		// write lead to csv
		var record []string
		record = append(record, lead.Email)
		for _, key := range leadDataKeys {
			record = append(record, key.Value)
		}
		record = append(record, lead.EmailIsValid)
		verificationResult, err := json.Marshal(lead.VerificationResult)
		if err != nil {
			return err
		}
		record = append(record, string(verificationResult))
		err = writer.Write(record)
		if err != nil {
			return err
		}
	}

	return nil
}
