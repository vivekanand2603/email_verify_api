package main

import (
	"encoding/json"
	"log"
	"sync"

	emailVerifier "github.com/AfterShip/email-verifier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func processQueue(client *mongo.Client) {
	log.Println("Processing queue")
	queue, err := GetQueue(client)
	if err != nil {
		log.Println(err)
		return
	}

	// Limit the number of concurrent goroutines
	const maxConcurrency = 10
	semaphore := make(chan struct{}, maxConcurrency)

	var wg sync.WaitGroup
	for _, q := range queue {
		semaphore <- struct{}{} // Acquire a token
		wg.Add(1)
		go func(q VerificationQueue) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release the token

			lead, err := GetLead(client, q.LeadID)
			if err != nil {
				log.Println(err)
			}

			verifier := emailVerifier.NewVerifier().EnableSMTPCheck().EnableAutoUpdateDisposable().EnableCatchAllCheck().EnableDomainSuggest().EnableGravatarCheck().HelloName("mx.google.com").FromEmail("vivek@thinksurfmedia.info")
			ret, err := verifier.Verify(lead.Email)
			if err != nil {
				log.Println(err)
			}

			bytes, err := json.Marshal(ret)
			if err != nil {
				log.Println(err)
			}
			reason := bson.M{"reachable": ret.Reachable, "reason": string(bytes)}
			err = Dequeue(client, q.ID, ret.Reachable, reason)
			if err != nil {
				log.Println(err)
			}
		}(q)
	}
}
