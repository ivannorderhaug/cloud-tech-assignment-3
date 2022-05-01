package webhook

import (
	"bytes"
	"corona-information-service/internal/model"
	"corona-information-service/pkg/api"
	"corona-information-service/pkg/customjson"
	"corona-information-service/pkg/db"
	"corona-information-service/tools/hash"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// COLLECTION can be changed
const COLLECTION = "notifications"

var webhooks []model.Webhook

// InitializeWebhooks retrieves all the webhooks from firestore
func InitializeWebhooks() {
	_, err := GetAllWebhooks()
	if err != nil {
		return
	}
}

// GetWebhook gets single webhook from local slice of webhooks
func GetWebhook(webhookId string) (model.Webhook, bool) {
	for _, wh := range webhooks {
		if wh.ID == webhookId {
			return wh, true
		}
	}
	return model.Webhook{}, false
}

// GetAllWebhooks retrieves all the webhooks from firestore
func GetAllWebhooks() ([]model.Webhook, error) {
	documentsFromFirestore, err := db.GetAllDocumentsFromFirestore(hash.Hash(COLLECTION))
	if err != nil {
		return []model.Webhook{}, err
	}
	webhooks = make([]model.Webhook, 0)
	//Converts each document snapshot into a webhook interface and adds it to the global webhooks slice
	for _, documentSnapshot := range documentsFromFirestore {
		var webhook model.Webhook
		err = documentSnapshot.DataTo(&webhook)
		if err != nil {
			return []model.Webhook{}, err
		}
		webhooks = append(webhooks, webhook)
	}
	return webhooks, nil
}

// DeleteWebhook deletes webhook locally and then from firestore
func DeleteWebhook(webhookId string) (bool, error) {
	var deleted = false
	if len(webhooks) != 0 {
		for i, wh := range webhooks {
			if wh.ID == webhookId {
				copy(webhooks[i:], webhooks[i+1:])          // Shift webhooks[i+1:] left one index.
				webhooks[len(webhooks)-1] = model.Webhook{} // Erase last element.
				webhooks = webhooks[:len(webhooks)-1]
				deleted = true
			}
		}
	}

	if err := db.DeleteSingleDocumentFromFirestore(hash.Hash(COLLECTION), hash.Hash(webhookId)); err != nil {
		return false, err
	}

	return deleted, nil
}

// RegisterWebhook registers a new webhook, and adds it to firestore
func RegisterWebhook(r *http.Request) (map[string]string, int, error) {
	var webhook model.Webhook
	var response = make(map[string]string, 1)

	//decode post request into webhook body
	err := customjson.Decode(r, &webhook)
	if err != nil {
		return map[string]string{}, http.StatusInternalServerError, err
	}

	//Check for duplicates
	for _, wh := range webhooks {
		if wh.Url == webhook.Url && wh.Country == webhook.Country {
			response["message"] = "There already exists a webhook with url: " + webhook.Url + " and country: " + webhook.Country
			return response, http.StatusConflict, nil
		}
	}

	//checks if alpha3 code was used as param for country
	if len(webhook.Country) == 3 && strings.ToLower(webhook.Country) != "usa" {
		fmt.Println(webhook.Country)
		//gets country name from restcountries api
		country, err := api.GetCountryNameByAlphaCode(webhook.Country)
		if err != nil {
			return map[string]string{}, http.StatusInternalServerError, err
		}
		webhook.Country = country
	}

	//autogenerate random id
	id := autoId()
	webhook.ID = id

	//Adds webhook to database
	err = db.AddToFirestore(hash.Hash(COLLECTION), hash.Hash(id), webhook)
	if err != nil {
		return map[string]string{}, http.StatusInternalServerError, err
	}
	webhooks = append(webhooks, webhook)

	response["id"] = id

	return response, http.StatusCreated, nil
}

// RunWebhookRoutine runs webhook routine for all webhooks
func RunWebhookRoutine(country string) error {
	for i, webhook := range webhooks {
		if webhook.Country == country {

			webhook.ActualCalls = webhook.ActualCalls + 1

			//Updates webhook in memory
			webhooks[i].ActualCalls = webhook.ActualCalls

			//Updates webhook in db
			err := db.UpdateDocument(hash.Hash(COLLECTION), hash.Hash(webhook.ID), "actual_calls", webhook.ActualCalls)
			if err != nil {
				return err
			}

			if webhook.ActualCalls == webhook.Calls {
				webhook.ActualCalls = 0

				//Updates webhook in db
				err = db.UpdateDocument(hash.Hash(COLLECTION), hash.Hash(webhook.ID), "actual_calls", webhook.ActualCalls)
				if err != nil {
					return err
				}

				//Updates webhook in memory
				webhooks[i].ActualCalls = webhook.ActualCalls

				webhook.Invoked = time.Now().UTC().String()

				go callUrl(webhook.Url, webhook)
			}

		}
	}
	return nil
}

// autoId Randomly generated a 15 character long string
func autoId() string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	arr := make([]rune, 15)
	for i := range arr {
		arr[i] = letters[rand.Intn(len(letters))]
	}

	return string(arr)
}

// callUrl sends post request to specified url
func callUrl(url string, data interface{}) {
	payloadBuffer := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuffer).Encode(data)
	if err == nil {
		_, err = http.Post(url, "application/json", payloadBuffer)
		if err != nil {
			return
		}
	}
}
