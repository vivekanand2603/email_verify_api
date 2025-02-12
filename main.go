package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func enableCORSAndJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017/email"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	router := httprouter.New()
	// list crud routes
	router.GET("/lists", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		lists, err := GetLists(client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(lists)
	})

	router.POST("/lists", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var list List
		json.NewDecoder(r.Body).Decode(&list)
		id, err := CreateList(client, list.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		list.ID = id
		json.NewEncoder(w).Encode(list)
	})

	router.GET("/lists/:id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		list, err := GetList(client, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(list)
	})

	router.DELETE("/lists/:id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = DeleteList(client, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	// lead crud routes

	router.GET("/leads", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		leads, err := GetLeads(client, primitive.NilObjectID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(leads)
	})

	router.POST("/leads", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var reqBody struct {
			Email  string `json:"email"`
			ListID string `json:"list_id"`
		}
		json.NewDecoder(r.Body).Decode(&reqBody)
		var lead Lead
		json.NewDecoder(r.Body).Decode(&lead)
		lead.ListID, err = primitive.ObjectIDFromHex(reqBody.ListID)
		lead.Email = reqBody.Email
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		leadData := make(map[string]interface{})
		id, err := CreateLead(client, lead.Email, lead.ListID, leadData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		lead.ID = id
		json.NewEncoder(w).Encode(lead)
	})

	router.GET("/leads/:id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		lead, err := GetLead(client, id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(lead)
	})

	router.DELETE("/leads/:id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = DeleteLead(client, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	// get leads by list id

	router.GET("/lists/:id/leads", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		leads, err := GetLeads(client, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(leads)
	})

	// get leads count by list id

	router.GET("/lists/:id/leads/count", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		count, err := GetLeadsCount(client, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(count)
	})

	// get leads count by list id and email_verified = true

	router.GET("/lists/:id/leads/count/email_verified", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		count, err := CountEmailVerified(client, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(count)
	})

	// get leads count by list id and email_is_valid = yes

	router.GET("/lists/:id/leads/count/valid_emails", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		count, err := CountValidEmails(client, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(count)
	})

	// get leads count by list id and email_is_valid = no

	router.GET("/lists/:id/leads/count/invalid_emails", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		count, err := CountInvalidEmails(client, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(count)
	})

	// get leads count by list id and email_is_valid = unknown

	router.GET("/lists/:id/leads/count/unknown_emails", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		count, err := CountUnknownEmails(client, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(count)
	})

	// count for all emails no matter the list

	router.GET("/count_all", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		emailIsValid := r.URL.Query().Get("email_is_valid")
		count, err := CountAllEmails(client, emailIsValid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(count)
	})

	router.POST("/lists/:id/queue", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = AddListToQueue(client, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	// is list in queue

	router.GET("/lists/:id/queue", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		inQueue, err := IsListInQueue(client, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JsonResponse := struct {
			InQueue bool `json:"in_queue"`
		}{
			InQueue: inQueue}
		json.NewEncoder(w).Encode(JsonResponse)
	})

	// remove list from queue

	router.DELETE("/lists/:id/queue", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = RemoveListFromQueue(client, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	router.GET("/processQueue", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		go processQueue(client)
		w.WriteHeader(http.StatusNoContent)
	})

	// get data from csv file and add to list

	router.POST("/lists/:id/leads/csv", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = AddLeadsFromCSV(client, id, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	// download csv file of leads

	router.GET("/lists/:id/leads/csv", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = DownloadLeadsAsCSV(client, id, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	// enable cors and content type json from headers for all routes
	wrappedRouter := enableCORSAndJSONContentType(router)
	log.Fatal(http.ListenAndServe(":30001", wrappedRouter))
}
