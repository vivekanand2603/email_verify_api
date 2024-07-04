package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type List struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
}

type Lead struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	Email              string             `bson:"email"`
	ListID             primitive.ObjectID `bson:"list_id"`
	LeadData           any                `bson:"lead_data"`
	EmailVerified      bool               `bson:"email_verified"`
	EmailIsValid       string             `bson:"email_is_valid"`
	VerificationResult any                `bson:"verification_result"`
}

type VerificationQueue struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Email  string             `bson:"email"`
	LeadID primitive.ObjectID `bson:"lead_id"`
	ListID primitive.ObjectID `bson:"list_id"`
}
